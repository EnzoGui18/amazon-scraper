// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/cors"
)

// Product define a estrutura dos dados que queremos extrair.
// As tags `json:"..."` são usadas para definir como cada campo
// será chamado quando convertido para o formato JSON.
type Product struct {
	Title    string `json:"title"`
	Rating   string `json:"rating"`
	Reviews  string `json:"reviews"`
	ImageURL string `json:"imageUrl"`
}

func main() {
	mux := http.NewServeMux() // Cria um novo multiplexador de requisições
	mux.HandleFunc("/api/scrape", scrapeHandler)

	// 2. ADICIONE O MIDDLEWARE DO CORS AQUI
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Permite o seu frontend
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}).Handler(mux)

	port := ":8080"
	fmt.Printf("Servidor iniciando na porta %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil { // Use o 'handler' com CORS
		log.Fatal(err)
	}
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obter a palavra-chave da URL (?keyword=...)
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		// Se a palavra-chave não for fornecida, retorna um erro 400 (Bad Request)
		http.Error(w, "A palavra-chave (keyword) é obrigatória", http.StatusBadRequest)
		return
	}

	// 2. Montar a URL de busca da Amazon
	// Usamos a Amazon do Brasil. A URL pode precisar de ajustes no futuro.
	searchURL := fmt.Sprintf("https://www.amazon.com.br/s?k=%s", keyword)

	// 3. Fazer a requisição HTTP para a Amazon
	// É CRUCIAL definir um User-Agent para simular um navegador real.
	// Sem isso, a Amazon pode bloquear a requisição ou retornar uma página diferente.
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		http.Error(w, "Falha ao criar requisição", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Falha ao buscar a página da Amazon", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Erro de status: %d %s", res.StatusCode, res.Status)
		http.Error(w, fmt.Sprintf("Amazon retornou um status inesperado: %d", res.StatusCode), http.StatusInternalServerError)
		return
	}

	// 4. Carregar o HTML da resposta no goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		http.Error(w, "Falha ao analisar o HTML", http.StatusInternalServerError)
		return
	}

	// 5. Extrair os dados de cada produto
	var products []Product
	// O seletor abaixo procura por cada "card" de produto na página de resultados.
	doc.Find("div[data-component-type='s-search-result']").Each(func(i int, s *goquery.Selection) {
		// Seletor final e corrigido para o título
		title := s.Find("a h2").Text()

		rating := s.Find("span.a-icon-alt").First().Text()
		reviews := s.Find("span.a-size-base.s-underline-text").Text()
		imageUrl, _ := s.Find("img.s-image").Attr("src")

		// Adiciona o produto à nossa lista, limpando espaços em branco
		products = append(products, Product{
			Title:    strings.TrimSpace(title),
			Rating:   strings.TrimSpace(rating),
			Reviews:  strings.TrimSpace(reviews),
			ImageURL: strings.TrimSpace(imageUrl),
		})
	})

	// 6. Converter os dados para JSON e enviar como resposta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Falha ao codificar a resposta em JSON", http.StatusInternalServerError)
	}
}

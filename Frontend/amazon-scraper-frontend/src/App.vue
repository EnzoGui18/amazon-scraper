<script setup>
import { ref } from 'vue'

// Váriaveis reativas para guardar o estado da nossa aplicação
const keyword = ref('') // Guarda o que o usuário digita no input
const products = ref([]) // Guarda a lista de produtos retornada pela API
const isLoading = ref(false) // Controla a exibição da mensagem "Carregando..."
const error = ref(null) // Guarda qualquer mensagem de erro

// Função que será chamada quando o botão for clicado
async function performSearch() {
  if (keyword.value.trim() === '') {
    error.value = 'Por favor, digite um termo para a busca.'
    return
  }

  // Limpa o estado anterior e ativa o loading
  isLoading.value = true
  products.value = []
  error.value = null

  try {
    // Faz a chamada para a nossa API Go
    const response = await fetch(`http://localhost:8080/api/scrape?keyword=${keyword.value}`)
    
    if (!response.ok) {
      throw new Error('A resposta do servidor não foi bem-sucedida.')
    }
    
    // Converte a resposta para JSON e a armazena na nossa váriavel 'products'
    const data = await response.json()
    products.value = data || [] // Garante que seja um array, mesmo que a resposta seja nula

  } catch (err) {
    // Em caso de erro na chamada, guarda a mensagem de erro
    console.error('Erro ao buscar dados:', err)
    error.value = 'Falha ao buscar os produtos. Verifique se o backend está rodando.'
  } finally {
    // Garante que o loading seja desativado ao final, com sucesso ou erro
    isLoading.value = false
  }
}
</script>

<template>
  <div class="container">
    <header>
      <h1>Amazon Busca</h1>
      <p>Digite um produto para buscar na Amazon</p>
    </header>

    <div class="search-box">
      <input 
        v-model="keyword" 
        @keyup.enter="performSearch"
        type="text" 
        placeholder="Ex: notebook, smartphone, etc."
      />
      <button @click="performSearch" :disabled="isLoading">
        {{ isLoading ? 'Buscando...' : 'Buscar' }}
      </button>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="isLoading" class="loading-message">
      Carregando resultados...
    </div>

    <div v-if="products.length > 0" class="results-grid">
      <div v-for="product in products" :key="product.imageUrl" class="product-card">
        <img :src="product.imageUrl" :alt="product.title" />
        <div class="product-info">
          <h3>{{ product.title || 'Título não disponível' }}</h3>
          <p class="rating">{{ product.rating }}</p>
          <p class="reviews">{{ product.reviews }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
:root {
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  background-color: #242424;
  color: rgba(255, 255, 255, 0.87);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

.search-box {
  margin: 2rem 0;
  display: flex;
  justify-content: center;
  gap: 10px;
}

input {
  width: 300px;
  padding: 10px;
  font-size: 1rem;
  border-radius: 6px;
  border: 1px solid #555;
  background-color: #333;
  color: #fff;
}

button {
  padding: 10px 20px;
  font-size: 1rem;
  cursor: pointer;
  border-radius: 6px;
  border: 1px solid transparent;
  background-color: #1a1a1a;
  transition: border-color 0.25s;
}
button:hover {
  border-color: #646cff;
}
button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.loading-message, .error-message {
  margin-top: 2rem;
  font-size: 1.2rem;
}
.error-message {
  color: #ff6b6b;
}

.results-grid {
  margin-top: 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
  text-align: left;
}

.product-card {
  background-color: #2f2f2f;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #444;
  transition: transform 0.2s;
}
.product-card:hover {
  transform: translateY(-5px);
}

.product-card img {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.product-card .product-info {
  padding: 1rem;
}

.product-card h3 {
  font-size: 1rem;
  margin: 0 0 0.5rem 0;
  height: 40px; /* Evita que os cards mudem de tamanho */
  overflow: hidden;
}

.product-card .rating {
  font-size: 0.9rem;
  color: #ffc107;
}

.product-card .reviews {
  font-size: 0.8rem;
  color: #aaa;
}
</style>
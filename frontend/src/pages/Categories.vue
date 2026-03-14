<template>
  <div class="page">
    <h1 class="page-title">分类</h1>
    <div class="grid">
      <router-link v-for="cat in categories" :key="cat.name" :to="`/categories/${cat.name}`" class="card">
        <div class="card-name">{{ cat.name }}</div>
        <div class="card-count">{{ cat.count }} 篇</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories } from '../api/index.js'
const categories = ref([])
onMounted(async () => {
  const r = await getCategories()
  categories.value = r.data.categories || []
})
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 24px; }
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 20px; text-align: center; transition: box-shadow 0.2s; }
.card:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.card-name { font-size: 16px; font-weight: 600; color: #18181b; margin-bottom: 4px; }
.card-count { font-size: 13px; color: #a1a1aa; }
@media (max-width: 640px) { .grid { grid-template-columns: repeat(2, 1fr); } }
</style>

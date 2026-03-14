<template>
  <div class="page">
    <h1 class="page-title">标签</h1>
    <div class="tags">
      <router-link v-for="tag in tags" :key="tag.name" :to="`/tags/${tag.name}`" class="tag">
        #{{ tag.name }} <span class="tag-count">{{ tag.count }}</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getTags } from '../api/index.js'
const tags = ref([])
onMounted(async () => {
  const r = await getTags()
  tags.value = r.data.tags || []
})
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 24px; }
.tags { display: flex; flex-wrap: wrap; gap: 12px; }
.tag { background: #fff; border: 1px solid #e5e5e5; border-radius: 20px; padding: 6px 14px; font-size: 14px; color: #52525b; transition: all 0.2s; }
.tag:hover { background: #eff6ff; border-color: #93c5fd; color: #2563eb; }
.tag-count { color: #a1a1aa; font-size: 12px; margin-left: 4px; }
</style>

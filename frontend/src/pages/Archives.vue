<template>
  <div class="page">
    <h1 class="page-title">归档</h1>
    <div v-for="arc in archives" :key="arc.month" class="archive-group">
      <div class="archive-month">
        <span>{{ arc.month }}</span>
        <span class="archive-count">{{ arc.count }} 篇</span>
      </div>
      <div class="archive-posts">
        <div v-for="post in arc.posts" :key="post.id" class="archive-post">
          <span class="archive-date">{{ formatDay(post.created_at) }}</span>
          <router-link :to="`/posts/${post.id}`" class="archive-title">{{ post.title }}</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getArchives, getPosts } from '../api/index.js'

const archives = ref([])

function formatDay(d) {
  const date = new Date(d)
  return `${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`
}

onMounted(async () => {
  const [arcRes, postRes] = await Promise.all([getArchives(), getPosts()])
  const arcMap = {}
  for (const arc of (arcRes.data.archives || [])) {
    arcMap[arc.month] = { ...arc, posts: [] }
  }
  for (const post of (postRes.data.posts || [])) {
    const month = new Date(post.created_at).toISOString().slice(0, 7)
    if (arcMap[month]) arcMap[month].posts.push(post)
  }
  archives.value = Object.values(arcMap)
})
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 24px; }
.archive-group { margin-bottom: 24px; }
.archive-month { font-size: 15px; font-weight: 600; color: #52525b; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
.archive-count { font-size: 12px; font-weight: 400; color: #a1a1aa; }
.archive-posts { border-left: 2px solid #e5e5e5; padding-left: 16px; }
.archive-post { display: flex; gap: 12px; margin-bottom: 8px; }
.archive-date { font-size: 12px; color: #a1a1aa; width: 50px; flex-shrink: 0; }
.archive-title { font-size: 14px; color: #3f3f46; }
.archive-title:hover { color: #2563eb; }
</style>

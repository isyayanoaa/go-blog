<template>
  <div class="page">
    <!-- 最新文章 -->
    <section class="section">
      <h2 class="section-title">最新文章</h2>
      <div v-if="posts.length === 0" class="empty">暂无文章</div>
      <div v-for="post in posts" :key="post.id" class="post-card">
        <router-link :to="`/posts/${post.id}`">
          <h3 class="post-title">{{ post.title }}</h3>
        </router-link>
        <p v-if="post.summary" class="post-summary">{{ post.summary }}</p>
        <div class="post-meta">
          <span>{{ formatDate(post.created_at) }}</span>
          <span v-if="post.category" class="tag">{{ post.category }}</span>
          <span>👍 {{ post.like_count }}</span>
        </div>
      </div>
      <router-link to="/posts" class="link">查看全部文章 →</router-link>
    </section>

    <!-- 最新动态 -->
    <section class="section">
      <h2 class="section-title">最新动态</h2>
      <div v-if="moments.length === 0" class="empty">暂无动态</div>
      <div v-for="m in moments" :key="m.id" class="moment-card">
        <p class="moment-content">{{ m.content }}</p>
        <div v-if="m.images && m.images.length" class="moment-images">
          <img v-for="img in m.images" :key="img" :src="img" class="moment-img" />
        </div>
        <div class="moment-meta">{{ formatDate(m.created_at) }} · 👍 {{ m.like_count }}</div>
      </div>
      <router-link to="/moments" class="link">查看全部动态 →</router-link>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPosts, getMoments } from '../api/index.js'

const posts = ref([])
const moments = ref([])

function formatDate(d) {
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(async () => {
  try {
    const r = await getPosts({ limit: 5 })
    posts.value = r.data.posts || []
  } catch {}
  try {
    const r = await getMoments({ limit: 3 })
    moments.value = r.data.moments || []
  } catch {}
})
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.section { margin-bottom: 48px; }
.section-title { font-size: 20px; font-weight: 600; margin-bottom: 20px; color: #18181b; }
.empty { color: #a1a1aa; text-align: center; padding: 32px 0; }
.post-card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 20px; margin-bottom: 16px; transition: box-shadow 0.2s; }
.post-card:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.post-title { font-size: 18px; font-weight: 600; color: #18181b; margin-bottom: 8px; }
.post-title:hover { color: #2563eb; }
.post-summary { color: #71717a; font-size: 14px; margin-bottom: 12px; }
.post-meta { display: flex; gap: 12px; font-size: 12px; color: #a1a1aa; }
.tag { background: #f4f4f5; padding: 2px 8px; border-radius: 4px; }
.moment-card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 20px; margin-bottom: 16px; }
.moment-content { color: #3f3f46; font-size: 14px; white-space: pre-wrap; margin-bottom: 12px; }
.moment-images { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 12px; }
.moment-img { width: 80px; height: 80px; object-fit: cover; border-radius: 8px; }
.moment-meta { font-size: 12px; color: #a1a1aa; }
.link { font-size: 14px; color: #2563eb; }
.link:hover { color: #1d4ed8; }
</style>

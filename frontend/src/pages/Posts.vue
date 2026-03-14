<template>
  <div class="page">
    <h1 class="page-title">文章</h1>
    <div v-if="posts.length === 0" class="empty">暂无文章</div>
    <div v-for="post in posts" :key="post.id" class="post-card">
      <router-link :to="`/posts/${post.id}`">
        <h2 class="post-title">{{ post.title }}</h2>
      </router-link>
      <p v-if="post.summary" class="post-summary">{{ post.summary }}</p>
      <div class="post-meta">
        <span>{{ formatDate(post.created_at) }}</span>
        <span v-if="post.category" class="tag">{{ post.category }}</span>
        <span v-for="tag in post.tags" :key="tag" class="tag tag-blue">#{{ tag }}</span>
        <span>👍 {{ post.like_count }}</span>
        <span>👁 {{ post.view_count }}</span>
      </div>
      <div v-if="userStore.isAdmin" class="admin-actions">
        <button class="btn-small" @click="editPost(post)">编辑</button>
        <button class="btn-small btn-danger" @click="deletePost(post.id)">删除</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPosts, deletePost as apiDeletePost } from '../api/index.js'
import { useUserStore } from '../stores/user.js'

const posts = ref([])
const userStore = useUserStore()

function formatDate(d) {
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

async function load() {
  const r = await getPosts()
  posts.value = r.data.posts || []
}

function editPost(post) {
  // TODO: 编辑弹窗
  alert('编辑功能开发中')
}

async function deletePost(id) {
  if (!confirm('确定删除？')) return
  await apiDeletePost(id)
  load()
}

onMounted(load)
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 24px; color: #18181b; }
.empty { color: #a1a1aa; text-align: center; padding: 48px 0; }
.post-card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 20px; margin-bottom: 16px; }
.post-title { font-size: 18px; font-weight: 600; color: #18181b; margin-bottom: 8px; }
.post-title:hover { color: #2563eb; }
.post-summary { color: #71717a; font-size: 14px; margin-bottom: 12px; }
.post-meta { display: flex; gap: 8px; font-size: 12px; color: #a1a1aa; flex-wrap: wrap; }
.tag { background: #f4f4f5; padding: 2px 8px; border-radius: 4px; }
.tag-blue { background: #eff6ff; color: #2563eb; }
.admin-actions { margin-top: 12px; display: flex; gap: 8px; }
.btn-small { font-size: 12px; padding: 4px 10px; border: 1px solid #e5e5e5; border-radius: 4px; background: #fff; cursor: pointer; }
.btn-small:hover { background: #f4f4f5; }
.btn-danger { color: #ef4444; border-color: #fca5a5; }
.btn-danger:hover { background: #fef2f2; }
</style>

<template>
  <div class="comments">
    <h4 class="comments-title">评论 ({{ comments.length }})</h4>
    <div v-for="c in comments" :key="c.id" class="comment">
      <img :src="c.avatar" class="comment-avatar" />
      <div class="comment-body">
        <div class="comment-header">
          <span class="comment-author">{{ c.username }}</span>
          <span class="comment-time">{{ formatDate(c.created_at) }}</span>
          <button v-if="userStore.isAdmin" class="comment-delete" @click="handleDelete(c.id)">删除</button>
        </div>
        <p class="comment-content">
          <span v-if="c.parent_id" class="reply-tag">↩ 回复</span>
          {{ c.content }}
        </p>
        <button v-if="userStore.isLoggedIn" class="reply-btn" @click="replyTo = c.id">回复</button>
        <div v-if="replyTo === c.id" class="reply-form">
          <input v-model="replyContent" class="reply-input" placeholder="回复..." />
          <button class="reply-submit" @click="submitComment(c.id)">发送</button>
          <button class="reply-cancel" @click="replyTo = null">取消</button>
        </div>
      </div>
    </div>

    <div v-if="userStore.isLoggedIn && !replyTo" class="comment-form">
      <textarea v-model="newComment" class="comment-input" placeholder="写下你的评论..."></textarea>
      <button class="comment-submit" @click="submitComment(null)">发送</button>
    </div>
    <div v-if="!userStore.isLoggedIn" class="login-hint">
      <a href="#" @click.prevent="$emit('login')" class="login-link">登录</a> 后发表评论
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getComments, createComment, deleteComment } from '../api/index.js'
import { useUserStore } from '../stores/user.js'

const props = defineProps({ targetType: String, targetId: Number })
const userStore = useUserStore()
const comments = ref([])
const newComment = ref('')
const replyContent = ref('')
const replyTo = ref(null)

function formatDate(d) {
  return new Date(d).toLocaleString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

async function load() {
  const r = await getComments({ target_type: props.targetType, target_id: props.targetId })
  comments.value = r.data.comments || []
}

async function submitComment(parentId) {
  const content = parentId ? replyContent.value : newComment.value
  if (!content.trim()) return
  await createComment({ target_type: props.targetType, target_id: props.targetId, parent_id: parentId, content })
  newComment.value = ''
  replyContent.value = ''
  replyTo.value = null
  load()
}

async function handleDelete(id) {
  if (!confirm('确定删除？')) return
  await deleteComment(id)
  load()
}

onMounted(load)
</script>

<style>
.comments { margin-top: 24px; }
.comments-title { font-size: 14px; font-weight: 600; color: #52525b; margin-bottom: 16px; }
.comment { display: flex; gap: 12px; margin-bottom: 16px; }
.comment-avatar { width: 28px; height: 28px; border-radius: 50%; flex-shrink: 0; }
.comment-body { flex: 1; }
.comment-header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.comment-author { font-size: 13px; font-weight: 500; color: #3f3f46; }
.comment-time { font-size: 12px; color: #a1a1aa; }
.comment-delete { font-size: 11px; color: #ef4444; background: none; border: none; cursor: pointer; margin-left: auto; }
.comment-content { font-size: 13px; color: #52525b; }
.reply-tag { color: #2563eb; font-size: 11px; margin-right: 4px; }
.reply-btn { font-size: 11px; color: #a1a1aa; background: none; border: none; cursor: pointer; margin-top: 4px; }
.reply-form { margin-top: 8px; display: flex; gap: 8px; }
.reply-input { flex: 1; padding: 6px 10px; border: 1px solid #e5e5e5; border-radius: 4px; font-size: 13px; }
.reply-submit, .reply-cancel { font-size: 12px; padding: 4px 10px; border: 1px solid #e5e5e5; border-radius: 4px; background: #fff; cursor: pointer; }
.reply-submit { background: #18181b; color: #fff; border-color: #18181b; }
.comment-form { margin-top: 16px; }
.comment-input { width: 100%; padding: 10px; border: 1px solid #e5e5e5; border-radius: 8px; font-size: 14px; resize: vertical; min-height: 80px; }
.comment-submit { margin-top: 8px; font-size: 14px; padding: 8px 16px; background: #18181b; color: #fff; border: none; border-radius: 6px; cursor: pointer; }
.login-hint { margin-top: 12px; font-size: 13px; color: #a1a1aa; }
.login-link { color: #2563eb; }
</style>

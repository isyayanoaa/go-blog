<template>
  <div class="page">
    <h1 class="page-title">动态</h1>
    <div v-if="moments.length === 0" class="empty">暂无动态</div>
    <div v-for="m in moments" :key="m.id" class="moment-card">
      <p v-if="m.content" class="moment-content">{{ m.content }}</p>
      <div v-if="m.images && m.images.length" class="moment-images" :class="getImgClass(m.images.length)">
        <div v-for="(img, i) in m.images" :key="img" class="img-wrap" @click="openLightbox(m.images, i)">
          <img :src="img" class="moment-img" />
        </div>
      </div>
      <div v-if="m.video" class="moment-video">
        <video :src="m.video" controls controlsList="nodownload" @contextmenu.prevent class="video-player"></video>
      </div>
      <div class="moment-footer">
        <span class="moment-time">{{ formatDate(m.created_at) }}</span>
        <div class="moment-actions">
          <button class="action-btn" :class="{ liked: likedMoments.has(m.id) }" @click="toggleLike(m)">
            👍 {{ m.like_count }}
          </button>
          <button v-if="userStore.isAdmin" class="action-btn btn-danger" @click="handleDelete(m.id)">删除</button>
        </div>
      </div>
      <CommentSection class="mt-4" target-type="moment" :target-id="m.id" />
    </div>

    <!-- 灯箱 -->
    <div v-if="lightbox.show" class="lightbox" @click="closeLightbox">
      <button class="lb-close" @click="closeLightbox">✕</button>
      <button v-if="lightbox.list.length > 1" class="lb-arrow lb-prev" @click.stop="lightboxStep(-1)">‹</button>
      <img
        :src="lightbox.list[lightbox.index]"
        class="lb-img"
        @click.stop
      />
      <button v-if="lightbox.list.length > 1" class="lb-arrow lb-next" @click.stop="lightboxStep(1)">›</button>
      <div v-if="lightbox.list.length > 1" class="lb-dots">
        <span v-for="(_, i) in lightbox.list" :key="i" class="lb-dot" :class="{ active: i === lightbox.index }" @click.stop="lightbox.index = i"></span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getMoments, deleteMoment, addLike, removeLike } from '../api/index.js'
import { useUserStore } from '../stores/user.js'
import CommentSection from '../components/CommentSection.vue'

const moments = ref([])
const userStore = useUserStore()
const likedMoments = ref(new Set())

const lightbox = ref({ show: false, list: [], index: 0 })

function openLightbox(images, index) {
  lightbox.value = { show: true, list: images, index }
  document.body.style.overflow = 'hidden'
}

function closeLightbox() {
  lightbox.value.show = false
  document.body.style.overflow = ''
}

function lightboxStep(dir) {
  const len = lightbox.value.list.length
  lightbox.value.index = (lightbox.value.index + dir + len) % len
}

function onKeydown(e) {
  if (!lightbox.value.show) return
  if (e.key === 'Escape') closeLightbox()
  if (e.key === 'ArrowLeft') lightboxStep(-1)
  if (e.key === 'ArrowRight') lightboxStep(1)
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

function formatDate(d) {
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function getImgClass(len) {
  if (len === 1) return 'grid-1'
  if (len <= 4) return 'grid-2'
  return 'grid-3'
}

async function toggleLike(m) {
  if (!userStore.isLoggedIn) return alert('请先登录')
  if (likedMoments.value.has(m.id)) {
    await removeLike({ target_type: 'moment', target_id: m.id })
    m.like_count--
    likedMoments.value.delete(m.id)
  } else {
    await addLike({ target_type: 'moment', target_id: m.id })
    m.like_count++
    likedMoments.value.add(m.id)
  }
}

async function handleDelete(id) {
  if (!confirm('确定删除？')) return
  await deleteMoment(id)
  moments.value = moments.value.filter(m => m.id !== id)
}

onMounted(async () => {
  const r = await getMoments()
  moments.value = r.data.moments || []
})
</script>

<style>
.page { max-width: 640px; margin: 0 auto; padding: 0 24px; }
.page-title { font-size: 24px; font-weight: 700; margin-bottom: 24px; color: #18181b; }
.empty { color: #a1a1aa; text-align: center; padding: 48px 0; }
.moment-card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 20px; margin-bottom: 16px; }
.moment-content { color: #3f3f46; font-size: 14px; white-space: pre-wrap; line-height: 1.6; margin-bottom: 12px; }
.moment-images { display: grid; gap: 4px; margin-bottom: 12px; }
.grid-1 { grid-template-columns: 1fr; max-width: 320px; }
.grid-2 { grid-template-columns: repeat(2, 1fr); }
.grid-3 { grid-template-columns: repeat(3, 1fr); }
.img-wrap { overflow: hidden; border-radius: 6px; cursor: zoom-in; aspect-ratio: 1; }
.moment-img { width: 100%; height: 100%; object-fit: cover; display: block; }
.img-wrap:hover .moment-img { opacity: 0.95; }
.moment-video { margin-bottom: 12px; }
.video-player { width: 100%; border-radius: 8px; max-height: 360px; background: #000; }
.moment-footer { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: #a1a1aa; }
.moment-actions { display: flex; gap: 8px; }
.action-btn { font-size: 12px; padding: 4px 10px; border: 1px solid #e5e5e5; border-radius: 4px; background: #fff; cursor: pointer; }
.action-btn:hover { background: #f4f4f5; }
.action-btn.liked { background: #2563eb; color: #fff; border-color: #2563eb; }
.btn-danger { color: #ef4444; border-color: #fca5a5; }
.mt-4 { margin-top: 16px; }

/* 灯箱 */
.lightbox { position: fixed; inset: 0; background: rgba(0,0,0,0.92); z-index: 1000; display: flex; align-items: center; justify-content: center; }
.lb-img { max-width: 90vw; max-height: 88vh; object-fit: contain; border-radius: 4px; }
.lb-close { position: fixed; top: 20px; right: 24px; background: none; border: none; color: #fff; font-size: 28px; cursor: pointer; z-index: 10; opacity: 0.8; }
.lb-close:hover { opacity: 1; }
.lb-arrow { position: fixed; top: 50%; transform: translateY(-50%); background: rgba(255,255,255,0.15); border: none; color: #fff; font-size: 48px; cursor: pointer; padding: 8px 16px; border-radius: 4px; line-height: 1; }
.lb-arrow:hover { background: rgba(255,255,255,0.3); }
.lb-prev { left: 16px; }
.lb-next { right: 16px; }
.lb-dots { position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%); display: flex; gap: 8px; }
.lb-dot { width: 8px; height: 8px; border-radius: 50%; background: rgba(255,255,255,0.4); cursor: pointer; }
.lb-dot.active { background: #fff; }
</style>

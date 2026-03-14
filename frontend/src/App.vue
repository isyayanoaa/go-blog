<template>
  <div id="app">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">go-blog</router-link>
        <nav class="nav">
          <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-link">{{ item.label }}</router-link>
        </nav>
        <div class="header-right">
          <template v-if="userStore.isLoggedIn">
            <img :src="userStore.user.avatar" class="avatar" />
            <span class="username">{{ userStore.user.username }}</span>
            <button v-if="userStore.isAdmin" class="btn btn-primary" @click="$router.push('/write')">写文章</button>
            <button v-if="userStore.isAdmin" class="btn" @click="showWriteMoment = true">发动态</button>
            <button class="btn" @click="handleLogout">退出</button>
          </template>
          <template v-else>
            <button class="btn btn-primary" @click="handleLogin">GitHub 登录</button>
          </template>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main :class="isWritePage ? 'main-full' : 'main'">
      <router-view />
    </main>

    <!-- 写文章弹窗 -->
    <div v-if="showWritePost" class="modal-mask" @click.self="showWritePost = false">
      <div class="modal">
        <div class="modal-header">
          <h3>写文章</h3>
          <button class="modal-close" @click="showWritePost = false">✕</button>
        </div>
        <div class="modal-body">
          <input v-model="postForm.title" class="input" placeholder="文章标题" />
          <input v-model="postForm.category" class="input" placeholder="分类（可选）" />
          <input v-model="postForm.tags" class="input" placeholder="标签，逗号分隔（可选）" />
          <input v-model="postForm.summary" class="input" placeholder="摘要（可选）" />
          <div style="position:relative">
            <textarea
              ref="postContentRef"
              v-model="postForm.content"
              class="textarea"
              placeholder="正文（支持 Markdown）。可直接粘贴图片自动上传。"
              @paste="onPostPaste"
            ></textarea>
            <span v-if="imageUploading" style="position:absolute;bottom:8px;right:12px;font-size:12px;color:#2563eb">图片上传中...</span>
          </div>
          <label class="checkbox-label">
            <input type="checkbox" v-model="postForm.published" /> 立即发布
          </label>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="showWritePost = false">取消</button>
          <button class="btn btn-primary" @click="submitPost" :disabled="posting">{{ posting ? '发布中...' : '发布' }}</button>
        </div>
      </div>
    </div>

    <!-- 发动态弹窗 -->
    <div v-if="showWriteMoment" class="modal-mask" @click.self="showWriteMoment = false">
      <div class="modal modal-small">
        <div class="modal-header">
          <h3>发动态</h3>
          <button class="modal-close" @click="showWriteMoment = false">✕</button>
        </div>
        <div class="modal-body">
          <textarea v-model="momentForm.content" class="textarea" placeholder="分享一下你的想法..." style="min-height:100px"></textarea>

          <!-- 媒体预览 -->
          <div v-if="momentForm.previews.length" class="media-previews">
            <div v-for="(p, i) in momentForm.previews" :key="i" class="media-preview-item">
              <img v-if="p.type === 'image'" :src="p.url" class="preview-img" />
              <video v-else :src="p.url" class="preview-video" controls></video>
              <button class="preview-remove" @click="removeMedia(i)">✕</button>
              <span v-if="p.uploading" class="preview-status">上传中...</span>
              <span v-else-if="p.error" class="preview-status error">失败</span>
            </div>
          </div>

          <!-- 选文件按钮 -->
          <div class="media-actions">
            <label class="btn media-btn">
              🖼 图片
              <input type="file" accept="image/*" multiple style="display:none" @change="onFileSelect($event, 'image')" />
            </label>
            <label class="btn media-btn">
              🎬 视频
              <input type="file" accept="video/*" style="display:none" @change="onFileSelect($event, 'video')" />
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="cancelMoment">取消</button>
          <button class="btn btn-primary" @click="submitMoment" :disabled="posting || momentForm.previews.some(p=>p.uploading)">
            {{ posting ? '发布中...' : momentForm.previews.some(p=>p.uploading) ? '上传中...' : '发布' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from './stores/user.js'
import { getMe, logout, createPost, createMoment, uploadFile } from './api/index.js'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const isWritePage = computed(() => route.path.startsWith('/write'))
const showWritePost = ref(false)
const showWriteMoment = ref(false)
const posting = ref(false)

const postForm = ref({ title: '', content: '', summary: '', category: '', tags: '', published: true })
const momentForm = ref({ content: '', previews: [] })
const postContentRef = ref(null)
const imageUploading = ref(false)
// previews: [{ type: 'image'|'video', url: localObjectURL, remoteUrl: '', uploading: true, error: false }]

const navItems = [
  { path: '/posts', label: '文章' },
  { path: '/moments', label: '动态' },
  { path: '/archives', label: '归档' },
  { path: '/categories', label: '分类' },
  { path: '/tags', label: '标签' },
  { path: '/about', label: '关于' },
]

onMounted(async () => {
  try {
    const res = await getMe()
    if (res.data.user) userStore.setUser(res.data.user)
  } catch {}
})

function handleLogin() {
  window.location.href = 'http://localhost:8080/api/v1/auth/github'
}

async function handleLogout() {
  await logout()
  userStore.logout()
  router.push('/')
}

async function submitPost() {
  if (!postForm.value.title || !postForm.value.content) return alert('标题和正文不能为空')
  posting.value = true
  try {
    const tags = postForm.value.tags ? postForm.value.tags.split(',').map(t => t.trim()).filter(Boolean) : []
    await createPost({
      title: postForm.value.title,
      content: postForm.value.content,
      summary: postForm.value.summary,
      category: postForm.value.category,
      tags,
      status: postForm.value.published ? 'published' : 'draft',
    })
    showWritePost.value = false
    postForm.value = { title: '', content: '', summary: '', category: '', tags: '', published: true }
    router.push('/posts')
  } catch (e) {
    alert('发布失败：' + (e.response?.data?.error || e.message))
  } finally {
    posting.value = false
  }
}

async function onPostPaste(e) {
  const items = Array.from(e.clipboardData?.items || [])
  const imageItem = items.find(item => item.type.startsWith('image/'))
  if (!imageItem) return // 普通文字粘贴，不拦截
  e.preventDefault()

  const file = imageItem.getAsFile()
  if (!file) return

  imageUploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file, `paste-${Date.now()}.png`)
    const res = await uploadFile(fd)
    const url = res.data.url
    // 在光标位置插入 Markdown 图片语法
    const ta = postContentRef.value
    const start = ta.selectionStart
    const end = ta.selectionEnd
    const md = `![image](${url})`
    const content = postForm.value.content
    postForm.value.content = content.slice(0, start) + md + content.slice(end)
    // 移动光标到插入内容后面
    await new Promise(r => setTimeout(r, 0))
    ta.selectionStart = ta.selectionEnd = start + md.length
  } catch {
    alert('图片上传失败')
  } finally {
    imageUploading.value = false
  }
}

async function onFileSelect(e, type) {
  const files = Array.from(e.target.files)
  e.target.value = '' // 允许重复选同一文件
  for (const file of files) {
    const localUrl = URL.createObjectURL(file)
    const item = { type, url: localUrl, remoteUrl: '', uploading: true, error: false }
    momentForm.value.previews.push(item)
    const fd = new FormData()
    fd.append('file', file)
    try {
      const res = await uploadFile(fd)
      item.remoteUrl = res.data.url
      item.uploading = false
    } catch {
      item.uploading = false
      item.error = true
    }
  }
}

function removeMedia(i) {
  const item = momentForm.value.previews[i]
  URL.revokeObjectURL(item.url)
  momentForm.value.previews.splice(i, 1)
}

function cancelMoment() {
  momentForm.value.previews.forEach(p => URL.revokeObjectURL(p.url))
  momentForm.value = { content: '', previews: [] }
  showWriteMoment.value = false
}

async function submitMoment() {
  if (!momentForm.value.content && momentForm.value.previews.length === 0) return alert('内容不能为空')
  if (momentForm.value.previews.some(p => p.uploading)) return alert('还有文件正在上传，请稍候')
  if (momentForm.value.previews.some(p => p.error)) return alert('部分文件上传失败，请移除后重试')
  posting.value = true
  try {
    const images = momentForm.value.previews.filter(p => p.type === 'image').map(p => p.remoteUrl)
    const video = momentForm.value.previews.find(p => p.type === 'video')?.remoteUrl || ''
    await createMoment({ content: momentForm.value.content, images, video })
    cancelMoment()
    router.push('/moments')
  } catch (e) {
    alert('发布失败：' + (e.response?.data?.error || e.message))
  } finally {
    posting.value = false
  }
}
</script>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f8f9fa; color: #1a1a1a; }
a { color: inherit; text-decoration: none; }
.header { background: #fff; border-bottom: 1px solid #e5e5e5; height: 56px; position: sticky; top: 0; z-index: 100; }
.header-inner { max-width: 960px; margin: 0 auto; padding: 0 24px; height: 100%; display: flex; align-items: center; justify-content: space-between; }
.logo { font-size: 20px; font-weight: 700; color: #18181b; }
.nav { display: flex; gap: 24px; }
.nav-link { font-size: 14px; color: #52525b; transition: color 0.2s; }
.nav-link:hover, .nav-link.router-link-active { color: #18181b; font-weight: 600; }
.header-right { display: flex; align-items: center; gap: 12px; }
.avatar { width: 32px; height: 32px; border-radius: 50%; }
.username { font-size: 14px; color: #52525b; }
.btn { font-size: 14px; padding: 6px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; cursor: pointer; transition: all 0.2s; }
.btn:hover { background: #f4f4f5; }
.btn-primary { background: #18181b; color: #fff; border-color: #18181b; }
.btn-primary:hover { background: #27272a; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.main { min-height: calc(100vh - 56px); padding: 32px 0; }
.main-full { height: calc(100vh - 56px); padding: 0; overflow: hidden; }

/* 弹窗 */
.modal-mask { position: fixed; inset: 0; background: rgba(0,0,0,0.4); z-index: 200; display: flex; align-items: center; justify-content: center; }
.modal { background: #fff; border-radius: 12px; width: 640px; max-width: 95vw; max-height: 90vh; display: flex; flex-direction: column; }
.modal-small { width: 480px; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 20px 24px; border-bottom: 1px solid #e5e5e5; }
.modal-header h3 { font-size: 16px; font-weight: 600; }
.modal-close { background: none; border: none; font-size: 18px; cursor: pointer; color: #a1a1aa; }
.modal-close:hover { color: #18181b; }
.modal-body { padding: 20px 24px; display: flex; flex-direction: column; gap: 12px; overflow-y: auto; }
.modal-footer { padding: 16px 24px; border-top: 1px solid #e5e5e5; display: flex; justify-content: flex-end; gap: 8px; }
.input { width: 100%; padding: 8px 12px; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 14px; outline: none; }
.input:focus { border-color: #a1a1aa; }
.textarea { width: 100%; padding: 10px 12px; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 14px; resize: vertical; min-height: 240px; outline: none; font-family: 'Menlo', monospace; }
.checkbox-label { display: flex; align-items: center; gap: 8px; font-size: 14px; color: #52525b; cursor: pointer; }
.media-previews { display: flex; flex-wrap: wrap; gap: 8px; }
.media-preview-item { position: relative; }
.preview-img { width: 80px; height: 80px; object-fit: cover; border-radius: 8px; display: block; }
.preview-video { width: 160px; height: 90px; border-radius: 8px; display: block; }
.preview-remove { position: absolute; top: 2px; right: 2px; background: rgba(0,0,0,0.5); color: #fff; border: none; border-radius: 50%; width: 18px; height: 18px; font-size: 10px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.preview-status { position: absolute; bottom: 4px; left: 0; right: 0; text-align: center; font-size: 10px; color: #fff; background: rgba(0,0,0,0.4); border-radius: 0 0 8px 8px; padding: 2px; }
.preview-status.error { background: rgba(220,38,38,0.7); }
.media-actions { display: flex; gap: 8px; }
.media-btn { font-size: 13px; padding: 6px 12px; cursor: pointer; }
</style>

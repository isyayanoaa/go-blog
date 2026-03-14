<template>
  <div class="write-page">
    <div class="toolbar-top">
      <div class="toolbar-left">
        <button class="back-btn" @click="$router.back()">← 返回</button>
        <input v-model="form.title" class="title-input" placeholder="文章标题..." />
      </div>
      <div class="toolbar-right">
        <input v-model="form.category" class="meta-input" placeholder="分类" />
        <input v-model="form.tags" class="meta-input" placeholder="标签（逗号分隔）" />
        <label class="publish-toggle">
          <input type="checkbox" v-model="form.published" />
          <span>{{ form.published ? '发布' : '草稿' }}</span>
        </label>
        <button class="btn-publish" @click="submit" :disabled="submitting">
          {{ submitting ? '保存中...' : (isEdit ? '保存修改' : '发布文章') }}
        </button>
      </div>
    </div>
    <div class="editor-wrap">
      <TiptapEditor v-if="ready" v-model="form.content" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import TiptapEditor from '../components/TiptapEditor.vue'
import { createPost, updatePost, getPost } from '../api/index.js'

const router = useRouter()
const route = useRoute()
const submitting = ref(false)
const ready = ref(false)

const isEdit = !!route.params.id

const form = ref({
  title: '',
  content: '',
  category: '',
  tags: '',
  published: true,
})

onMounted(async () => {
  if (isEdit) {
    try {
      const r = await getPost(route.params.id)
      const p = r.data.post
      form.value.title = p.title
      form.value.content = p.content
      form.value.category = p.category || ''
      form.value.tags = (p.tags || []).join(', ')
      form.value.published = p.status === 'published'
    } catch (e) {
      alert('加载文章失败')
      router.back()
    }
  }
  ready.value = true
})

async function submit() {
  if (!form.value.title.trim()) return alert('请填写标题')
  if (!form.value.content || form.value.content === '<p></p>') return alert('请填写正文')
  submitting.value = true
  try {
    const tags = form.value.tags ? form.value.tags.split(',').map(t => t.trim()).filter(Boolean) : []
    const tmp = document.createElement('div')
    tmp.innerHTML = form.value.content
    const summary = tmp.textContent.slice(0, 100)
    const payload = {
      title: form.value.title,
      content: form.value.content,
      summary,
      category: form.value.category,
      tags,
      status: form.value.published ? 'published' : 'draft',
    }
    if (isEdit) {
      await updatePost(route.params.id, payload)
      router.push('/posts/' + route.params.id)
    } else {
      await createPost(payload)
      router.push('/posts')
    }
  } catch (e) {
    alert('保存失败：' + (e.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.write-page { display: flex; flex-direction: column; height: 100vh; background: #fff; overflow: hidden; }
.toolbar-top { display: flex; align-items: center; justify-content: space-between; padding: 0 24px; height: 56px; border-bottom: 1px solid #e5e5e5; flex-shrink: 0; gap: 16px; }
.toolbar-left { display: flex; align-items: center; gap: 16px; flex: 1; min-width: 0; }
.toolbar-right { display: flex; align-items: center; gap: 12px; flex-shrink: 0; }
.back-btn { font-size: 14px; color: #71717a; background: none; border: none; cursor: pointer; white-space: nowrap; padding: 0; }
.back-btn:hover { color: #18181b; }
.title-input { font-size: 18px; font-weight: 600; border: none; outline: none; flex: 1; min-width: 0; color: #18181b; }
.title-input::placeholder { color: #d4d4d8; }
.meta-input { font-size: 13px; padding: 5px 10px; border: 1px solid #e5e5e5; border-radius: 6px; outline: none; width: 100px; }
.meta-input:focus { border-color: #a1a1aa; }
.publish-toggle { display: flex; align-items: center; gap: 6px; font-size: 13px; color: #52525b; cursor: pointer; white-space: nowrap; }
.btn-publish { font-size: 14px; padding: 7px 16px; background: #18181b; color: #fff; border: none; border-radius: 6px; cursor: pointer; white-space: nowrap; }
.btn-publish:hover { background: #27272a; }
.btn-publish:disabled { opacity: 0.5; cursor: not-allowed; }
.editor-wrap { flex: 1; overflow: hidden; display: flex; flex-direction: column; }
</style>

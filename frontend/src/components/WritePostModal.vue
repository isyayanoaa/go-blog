<template>
  <n-modal v-model:show="show" preset="card" title="写文章" style="width: 800px; max-width: 95vw;">
    <div class="flex flex-col gap-3">
      <n-input v-model:value="form.title" placeholder="文章标题" size="large" />
      <div class="flex gap-2">
        <n-input v-model:value="form.category" placeholder="分类" style="flex:1" />
        <n-input v-model:value="tagsInput" placeholder="标签（逗号分隔）" style="flex:2" />
        <n-select v-model:value="form.status" :options="statusOptions" style="width:100px" />
      </div>
      <n-input v-model:value="form.summary" placeholder="摘要（可选）" type="textarea" :rows="2" />
      <n-input v-model:value="form.content" placeholder="Markdown 内容..." type="textarea" :rows="16" style="font-family: monospace; font-size: 13px;" />
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="show = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="save">保存</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { createPost, updatePost } from '../api/index.js'

const props = defineProps({ post: Object })
const emit = defineEmits(['update:show', 'saved'])
const show = defineModel('show')
const message = useMessage()
const saving = ref(false)
const tagsInput = ref('')

const form = ref({ title: '', content: '', summary: '', cover: '', category: '', tags: [], status: 'published' })
const statusOptions = [
  { label: '发布', value: 'published' },
  { label: '草稿', value: 'draft' },
]

watch(() => props.post, (p) => {
  if (p) {
    form.value = { ...p }
    tagsInput.value = (p.tags || []).join(', ')
  } else {
    form.value = { title: '', content: '', summary: '', cover: '', category: '', tags: [], status: 'published' }
    tagsInput.value = ''
  }
}, { immediate: true })

async function save() {
  form.value.tags = tagsInput.value.split(',').map(t => t.trim()).filter(Boolean)
  saving.value = true
  try {
    if (props.post?.id) {
      await updatePost(props.post.id, form.value)
    } else {
      await createPost(form.value)
    }
    message.success('保存成功')
    emit('saved')
  } finally { saving.value = false }
}
</script>

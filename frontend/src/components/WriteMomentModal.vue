<template>
  <n-modal v-model:show="show" preset="card" title="发动态" style="width: 480px; max-width: 95vw;">
    <div class="flex flex-col gap-3">
      <n-input v-model:value="content" type="textarea" :rows="5" placeholder="说点什么..." />
      <n-upload multiple accept="image/*" :custom-request="uploadImage" list-type="image-card">
        <n-button>上传图片</n-button>
      </n-upload>
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="show = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="save">发布</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup>
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { createMoment, uploadFile } from '../api/index.js'

const show = defineModel('show')
const emit = defineEmits(['saved'])
const message = useMessage()
const content = ref('')
const images = ref([])
const saving = ref(false)

async function uploadImage({ file, onFinish, onError }) {
  const fd = new FormData()
  fd.append('file', file.file)
  try {
    const r = await uploadFile(fd)
    images.value.push('http://localhost:8080' + r.data.url)
    onFinish()
  } catch { onError() }
}

async function save() {
  saving.value = true
  try {
    await createMoment({ content: content.value, images: images.value })
    message.success('发布成功')
    content.value = ''
    images.value = []
    emit('saved')
  } finally { saving.value = false }
}
</script>

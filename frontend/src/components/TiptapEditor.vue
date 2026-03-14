<template>
  <div class="editor-container">
    <!-- 工具栏 -->
    <div class="toolbar" v-if="editor">
      <!-- 标题 -->
      <button @click="editor.chain().focus().toggleHeading({level:1}).run()" :class="{active: editor.isActive('heading',{level:1})}" title="一级标题 (#+空格)">H1</button>
      <button @click="editor.chain().focus().toggleHeading({level:2}).run()" :class="{active: editor.isActive('heading',{level:2})}" title="二级标题 (##+空格)">H2</button>
      <button @click="editor.chain().focus().toggleHeading({level:3}).run()" :class="{active: editor.isActive('heading',{level:3})}" title="三级标题 (###+空格)">H3</button>
      <div class="divider"></div>
      
      <!-- 文本格式 -->
      <button @click="editor.chain().focus().toggleBold().run()" :class="{active: editor.isActive('bold')}" title="加粗 (Cmd+B)"><b>B</b></button>
      <button @click="editor.chain().focus().toggleItalic().run()" :class="{active: editor.isActive('italic')}" title="斜体 (Cmd+I)"><i>I</i></button>
      <button @click="editor.chain().focus().toggleUnderline().run()" :class="{active: editor.isActive('underline')}" title="下划线 (Cmd+U)"><u>U</u></button>
      <button @click="editor.chain().focus().toggleStrike().run()" :class="{active: editor.isActive('strike')}" title="删除线"><s>S</s></button>
      <button @click="editor.chain().focus().toggleCode().run()" :class="{active: editor.isActive('code')}" title="行内代码">`</button>
      <button @click="setHighlight" :class="{active: editor.isActive('highlight')}" title="高亮背景">🖍</button>
      <button @click="setTextColor" title="文字颜色">🎨</button>
      <div class="divider"></div>
      
      <!-- 列表 -->
      <button @click="editor.chain().focus().toggleBulletList().run()" :class="{active: editor.isActive('bulletList')}" title="无序列表 (-+空格)">• 列表</button>
      <button @click="editor.chain().focus().toggleOrderedList().run()" :class="{active: editor.isActive('orderedList')}" title="有序列表 (1.+空格)">1. 列表</button>
      <button @click="editor.chain().focus().toggleTaskList().run()" :class="{active: editor.isActive('taskList')}" title="任务清单">☑ 任务</button>
      <div class="divider"></div>
      
      <!-- 块级 -->
      <button @click="editor.chain().focus().toggleBlockquote().run()" :class="{active: editor.isActive('blockquote')}" title="引用块 (>+空格)">" 引用</button>
      <button @click="editor.chain().focus().toggleCodeBlock().run()" :class="{active: editor.isActive('codeBlock')}" title="代码块 (```+空格)">{ }</button>
      <select v-if="editor.isActive('codeBlock')" class="lang-select" @change="setCodeLang($event.target.value)" :value="currentLang">
        <option value="">自动</option>
        <option v-for="lang in languages" :key="lang.value" :value="lang.value">{{ lang.label }}</option>
      </select>
      <button @click="insertHr" title="分割线 (---)">—</button>
      <button @click="insertTable" title="表格">⊞ 表格</button>
      <div class="divider"></div>
      
      <!-- 媒体 -->
      <label class="toolbar-btn" title="图片 (支持粘贴)">
        🖼 图片
        <input type="file" accept="image/*" multiple style="display:none" @change="insertImages" />
      </label>
      <button @click="setLink" :class="{active: editor.isActive('link')}" title="链接 (Cmd+K)">🔗 链接</button>
      <div class="divider"></div>
      
      <!-- 对齐 -->
      <button @click="editor.chain().focus().setTextAlign('left').run()" :class="{active: editor.isActive({textAlign:'left'})}" title="左对齐">⬅</button>
      <button @click="editor.chain().focus().setTextAlign('center').run()" :class="{active: editor.isActive({textAlign:'center'})}" title="居中">⬌</button>
      <button @click="editor.chain().focus().setTextAlign('right').run()" :class="{active: editor.isActive({textAlign:'right'})}" title="右对齐">➡</button>
    </div>

    <!-- 编辑区 -->
    <editor-content :editor="editor" class="editor-body" />

    <!-- 上传提示 -->
    <div v-if="uploading" class="upload-toast">📤 图片上传中...</div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Underline from '@tiptap/extension-underline'
import { Color } from '@tiptap/extension-color'
import { TextStyle } from '@tiptap/extension-text-style'
import { Table, TableRow, TableCell, TableHeader } from '@tiptap/extension-table'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import Typography from '@tiptap/extension-typography'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { createLowlight, common } from 'lowlight'
import { uploadFile } from '../api/index.js'

const props = defineProps({ modelValue: String })
const emit = defineEmits(['update:modelValue'])
const uploading = ref(false)

const languages = [
  { label: 'JavaScript', value: 'javascript' },
  { label: 'TypeScript', value: 'typescript' },
  { label: 'Go', value: 'go' },
  { label: 'Python', value: 'python' },
  { label: 'Rust', value: 'rust' },
  { label: 'Java', value: 'java' },
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'C#', value: 'csharp' },
  { label: 'PHP', value: 'php' },
  { label: 'Ruby', value: 'ruby' },
  { label: 'Swift', value: 'swift' },
  { label: 'Kotlin', value: 'kotlin' },
  { label: 'Dart', value: 'dart' },
  { label: 'Scala', value: 'scala' },
  { label: 'Shell/Bash', value: 'bash' },
  { label: 'SQL', value: 'sql' },
  { label: 'HTML', value: 'html' },
  { label: 'CSS', value: 'css' },
  { label: 'Vue', value: 'xml' },
  { label: 'JSON', value: 'json' },
  { label: 'YAML', value: 'yaml' },
  { label: 'Markdown', value: 'markdown' },
  { label: 'Dockerfile', value: 'dockerfile' },
  { label: 'Nginx', value: 'nginx' },
  { label: 'XML', value: 'xml' },
  { label: 'GraphQL', value: 'graphql' },
  { label: 'Protobuf', value: 'protobuf' },
]

const currentLang = computed(() => {
  if (!editor.value) return ''
  return editor.value.getAttributes('codeBlock').language || ''
})

const lowlight = createLowlight(common)

const editor = useEditor({
  content: props.modelValue || '',
  extensions: [
    StarterKit.configure({ codeBlock: false }),
    Placeholder.configure({ placeholder: '开始写作... 输入 / 触发命令，# 空格 创建标题' }),
    Image.configure({ inline: false, allowBase64: false }),
    CodeBlockLowlight.configure({ lowlight }),
    Link.configure({ openOnClick: false, HTMLAttributes: { target: '_blank', rel: 'noopener noreferrer' } }),
    Underline,
    Table.configure({ resizable: true }),
    TableRow,
    TableCell,
    TableHeader,
    TaskList,
    TaskItem.configure({ nested: true }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Highlight,
    Typography,
    TextStyle,
    Color,
  ],
  onUpdate({ editor }) {
    emit('update:modelValue', editor.getHTML())
  },
  editorProps: {
    handlePaste(view, event) {
      const items = Array.from(event.clipboardData?.items || [])
      const imageItem = items.find(i => i.type.startsWith('image/'))
      if (!imageItem) return false
      event.preventDefault()
      const file = imageItem.getAsFile()
      if (file) uploadImageFile(file)
      return true
    },
  },
})

watch(() => props.modelValue, (val) => {
  if (editor.value && val !== editor.value.getHTML()) {
    editor.value.commands.setContent(val || '', false)
  }
})

async function uploadImageFile(file) {
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file, file.name || `paste-${Date.now()}.png`)
    const res = await uploadFile(fd)
    editor.value.chain().focus().setImage({ src: res.data.url }).run()
  } catch { alert('图片上传失败') }
  finally { uploading.value = false }
}

async function insertImages(e) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  e.target.value = ''
  uploading.value = true
  try {
    for (const file of files) {
      const fd = new FormData()
      fd.append('file', file)
      const res = await uploadFile(fd)
      editor.value.chain().focus().setImage({ src: res.data.url }).run()
    }
  } catch { alert('图片上传失败') }
  finally { uploading.value = false }
}

function setHighlight() {
  editor.value.chain().focus().toggleHighlight().run()
}

function setTextColor() {
  const color = prompt('请输入颜色（如 #ff0000 或 red）:', '#ff0000')
  if (!color) return
  editor.value.chain().focus().setColor(color).run()
}

function setLink() {
  const url = prompt('请输入链接地址:', 'https://')
  if (!url) return
  if (editor.value.isActive('link')) {
    editor.value.chain().focus().unsetLink().run()
  } else {
    editor.value.chain().focus().setLink({ href: url }).run()
  }
}

function setCodeLang(lang) {
  editor.value.chain().focus().updateAttributes('codeBlock', { language: lang || null }).run()
}

function insertHr() {
  editor.value.chain().focus().setHorizontalRule().run()
}

function insertTable() {
  editor.value.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
}

onBeforeUnmount(() => editor.value?.destroy())
</script>

<style>
.editor-container { position: relative; display: flex; flex-direction: column; height: 100%; }
.toolbar { display: flex; align-items: center; gap: 2px; padding: 6px 12px; border-bottom: 1px solid #e5e5e5; flex-wrap: wrap; flex-shrink: 0; background: #fafafa; }
.toolbar button, .toolbar .toolbar-btn { font-size: 13px; padding: 4px 8px; border: none; border-radius: 4px; background: none; cursor: pointer; color: #52525b; white-space: nowrap; }
.toolbar button:hover, .toolbar .toolbar-btn:hover { background: #f0f0f0; }
.toolbar button.active { background: #e0e7ff; color: #2563eb; }
.toolbar .divider { width: 1px; height: 18px; background: #e5e5e5; margin: 0 4px; }
.lang-select { font-size: 12px; padding: 3px 6px; border: 1px solid #e5e5e5; border-radius: 4px; background: #fff; color: #52525b; cursor: pointer; }
.editor-body { flex: 1; overflow-y: auto; padding: 32px 48px; }
.editor-body .tiptap { outline: none; min-height: 100%; font-size: 15px; line-height: 1.8; color: #3f3f46; }
.editor-body .tiptap p.is-editor-empty:first-child::before { content: attr(data-placeholder); color: #d4d4d8; pointer-events: none; float: left; height: 0; }
.tiptap h1 { font-size: 2em; font-weight: 700; margin: 1.2em 0 0.4em; color: #18181b; }
.tiptap h2 { font-size: 1.5em; font-weight: 700; margin: 1.2em 0 0.4em; color: #18181b; border-bottom: 1px solid #e5e5e5; padding-bottom: 0.3em; }
.tiptap h3 { font-size: 1.2em; font-weight: 600; margin: 1em 0 0.3em; color: #18181b; }
.tiptap p { margin: 0.6em 0; }
.tiptap strong { font-weight: 700; }
.tiptap em { font-style: italic; }
.tiptap u { text-decoration: underline; }
.tiptap s { text-decoration: line-through; color: #a1a1aa; }
.tiptap code { background: #f4f4f5; padding: 2px 6px; border-radius: 4px; font-size: 0.85em; font-family: 'Menlo', monospace; color: #e11d48; }
.tiptap mark { background: #fef08a; padding: 1px 2px; border-radius: 2px; }
.tiptap ul { list-style-type: disc; padding-left: 1.8em; margin: 0.5em 0; }
.tiptap ol { list-style-type: decimal; padding-left: 1.8em; margin: 0.5em 0; }
.tiptap ul ul { list-style-type: circle; }
.tiptap li { margin: 0.2em 0; display: list-item; }
.tiptap ul[data-type="taskList"] { list-style: none; padding-left: 0; }
.tiptap ul[data-type="taskList"] li { display: flex; align-items: flex-start; gap: 8px; }
.tiptap ul[data-type="taskList"] li > label { flex-shrink: 0; margin-top: 4px; }
.tiptap ul[data-type="taskList"] li > div { flex: 1; }
.tiptap blockquote { border-left: 3px solid #2563eb; margin: 1em 0; padding: 4px 16px; color: #71717a; background: #f8faff; border-radius: 0 6px 6px 0; }
.tiptap pre { background: #1e1e2e; color: #cdd6f4; padding: 16px 20px; border-radius: 10px; overflow-x: auto; margin: 1em 0; }
.tiptap pre code { background: none; color: inherit; padding: 0; font-size: 0.9em; }
.tiptap hr { border: none; border-top: 1px solid #e5e5e5; margin: 1.5em 0; }
.tiptap a { color: #2563eb; text-decoration: underline; }
.tiptap img { max-width: 100%; border-radius: 8px; margin: 0.5em 0; display: block; }
.tiptap img.ProseMirror-selectednode { outline: 2px solid #2563eb; }
.tiptap table { border-collapse: collapse; width: 100%; margin: 1em 0; }
.tiptap th, .tiptap td { border: 1px solid #e5e5e5; padding: 8px 12px; text-align: left; }
.tiptap th { background: #f4f4f5; font-weight: 600; }
.tiptap .selectedCell { background: #e0e7ff; }
.upload-toast { position: absolute; bottom: 16px; right: 16px; background: #18181b; color: #fff; font-size: 13px; padding: 8px 14px; border-radius: 6px; z-index: 10; }
</style>

<template>
  <div class="page">
    <template v-if="post">
      <h1 class="article-title">{{ post.title }}</h1>
      <div class="article-meta">
        <span>{{ formatDate(post.created_at) }}</span>
        <span v-if="post.category" class="tag">{{ post.category }}</span>
        <span v-for="tag in post.tags" :key="tag" class="tag tag-blue">#{{ tag }}</span>
        <span>👁 {{ post.view_count }}</span>
      </div>
      <div v-if="userStore.isAdmin" class="admin-actions">
        <button class="btn-small" @click="router.push('/write/' + post.id)">编辑</button>
        <button class="btn-small btn-danger" @click="handleDelete">删除</button>
      </div>

      <article ref="articleRef" class="article-content" v-html="renderedContent"></article>

      <div class="like-bar">
        <button class="like-btn" :class="{ liked }" @click="toggleLike">
          👍 {{ post.like_count }}
        </button>
      </div>

      <CommentSection target-type="post" :target-id="post.id" />
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from '../utils/markdown.js'
import { getPost, addLike, removeLike, deletePost } from '../api/index.js'
import { useUserStore } from '../stores/user.js'
import CommentSection from '../components/CommentSection.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const post = ref(null)
const liked = ref(false)
const showEdit = ref(false)
const articleRef = ref(null)

// 运行结果：key = 代码块索引，value = { loading, output, error }
const runResults = ref({})

const renderedContent = computed(() => {
  if (!post.value?.content) return ''
  const c = post.value.content.trim()
  if (c.startsWith('<')) return c
  return marked(c)
})

function formatDate(d) {
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// 给每个代码块注入运行按钮 + 语言标签 + 复制按钮
function injectCodeActions() {
  if (!articleRef.value) return
  const blocks = articleRef.value.querySelectorAll('pre')
  blocks.forEach((pre, idx) => {
    if (pre.querySelector('.code-toolbar')) return // 已注入过
    const code = pre.querySelector('code')
    if (!code) return

    // 检测语言
    const langClass = Array.from(code.classList).find(c => c.startsWith('language-'))
    const lang = langClass ? langClass.replace('language-', '') : ''

    // 工具栏容器
    const toolbar = document.createElement('div')
    toolbar.className = 'code-toolbar'

    // 语言标签
    if (lang) {
      const langBadge = document.createElement('span')
      langBadge.className = 'code-lang'
      langBadge.textContent = lang
      toolbar.appendChild(langBadge)
    }

    // 复制按钮
    const copyBtn = document.createElement('button')
    copyBtn.className = 'code-btn'
    copyBtn.textContent = '复制'
    copyBtn.onclick = () => {
      navigator.clipboard.writeText(code.textContent || '')
      copyBtn.textContent = '已复制!'
      setTimeout(() => { copyBtn.textContent = '复制' }, 1500)
    }
    toolbar.appendChild(copyBtn)

    // 运行按钮（仅 Go 语言）
    if (lang === 'go' || lang === 'golang') {
      const runBtn = document.createElement('button')
      runBtn.className = 'code-btn code-run-btn'
      runBtn.textContent = '▶ 运行'
      runBtn.onclick = () => runCode(idx, code.textContent || '', runBtn, pre)
      toolbar.appendChild(runBtn)
    }

    // 把工具栏插在 pre 里顶部
    pre.style.position = 'relative'
    pre.style.paddingTop = '40px'
    pre.insertBefore(toolbar, pre.firstChild)

    // 创建输出区（隐藏状态）
    const outputEl = document.createElement('div')
    outputEl.className = 'code-output'
    outputEl.id = `code-output-${idx}`
    outputEl.style.display = 'none'
    pre.appendChild(outputEl)
  })
}

// Go 标准库包名 → import 路径映射（常用包）
const GO_STD_PKGS = {
  fmt: 'fmt', os: 'os', io: 'io', log: 'log', math: 'math',
  rand: 'math/rand', big: 'math/big', bits: 'math/bits',
  strings: 'strings', strconv: 'strconv', unicode: 'unicode', utf8: 'unicode/utf8',
  bytes: 'bytes', bufio: 'bufio', errors: 'errors',
  sort: 'sort', slices: 'slices', maps: 'maps',
  time: 'time', context: 'context',
  sync: 'sync', atomic: 'sync/atomic',
  http: 'net/http', url: 'net/url', net: 'net',
  json: 'encoding/json', xml: 'encoding/xml', csv: 'encoding/csv',
  base64: 'encoding/base64', hex: 'encoding/hex', binary: 'encoding/binary',
  md5: 'crypto/md5', sha256: 'crypto/sha256', sha512: 'crypto/sha512',
  rand2: 'crypto/rand',
  filepath: 'path/filepath', path: 'path',
  regexp: 'regexp', template: 'text/template', html: 'html/template',
  reflect: 'reflect', runtime: 'runtime',
  flag: 'flag', testing: 'testing',
  heap: 'container/heap', list: 'container/list', ring: 'container/ring',
  ioutil: 'io/ioutil', fs: 'io/fs',
  exec: 'os/exec', signal: 'os/signal',
  image: 'image', color: 'image/color', png: 'image/png',
  zip: 'archive/zip', tar: 'archive/tar',
  gzip: 'compress/gzip', zlib: 'compress/zlib',
}

function wrapGoCode(raw) {
  const code = raw.trim()
  // 已经是完整程序，直接用
  if (code.startsWith('package ')) return code

  // 分离 import 行和正文行
  const lines = code.split('\n')
  const manualImports = new Set()
  const bodyLines = []
  let inImportBlock = false

  for (const line of lines) {
    const trimmed = line.trim()
    if (/^import\s*\($/.test(trimmed)) { inImportBlock = true; continue }
    if (inImportBlock) {
      if (trimmed === ')') { inImportBlock = false; continue }
      const m = trimmed.match(/"([^"]+)"/)
      if (m) manualImports.add(m[1])
      continue
    }
    const single = trimmed.match(/^import\s+"([^"]+)"/)
    if (single) { manualImports.add(single[1]); continue }
    bodyLines.push(line)
  }

  // 自动检测正文里用到了哪些标准库包
  const body = bodyLines.join('\n')
  const autoImports = new Set(manualImports)

  // 扫描 pkgName.Xxx 调用
  const usedPkgs = new Set()
  const callRe = /\b([a-z][a-zA-Z0-9]*)\.([A-Za-z])/g
  let m
  while ((m = callRe.exec(body)) !== null) {
    usedPkgs.add(m[1])
  }

  // 对照标准库表自动补 import
  for (const pkg of usedPkgs) {
    if (GO_STD_PKGS[pkg] && !Array.from(autoImports).some(p => p.endsWith('/' + pkg) || p === pkg)) {
      autoImports.add(GO_STD_PKGS[pkg])
    }
  }

  // 构建 import 块
  let importSection = ''
  if (autoImports.size > 0) {
    importSection = 'import (\n' + Array.from(autoImports).map(p => `\t"${p}"`).join('\n') + '\n)\n\n'
  }

  return `package main\n\n${importSection}func main() {\n${bodyLines.map(l => '\t' + l).join('\n')}\n}\n`
}

async function runCode(idx, code, btn, pre) {
  btn.textContent = '⏳ 运行中...'
  btn.disabled = true
  const outputEl = pre.querySelector(`#code-output-${idx}`)
  if (outputEl) { outputEl.style.display = 'block'; outputEl.textContent = '运行中...' }

  const wrappedCode = wrapGoCode(code)

  try {
    const res = await fetch('/api/v1/playground/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code: wrappedCode }),
    })
    const data = await res.json()
    // play.golang.org 返回格式：{ Errors, Events: [{Message}] }
    // go.dev 返回格式也类似
    const hasError = data.Errors || data.errors
    const output = data.Events?.map(e => e.Message).join('')
      || data.events?.map(e => e.message).join('')
      || data.output
      || '(无输出)'
    if (outputEl) {
      outputEl.style.display = 'block'
      outputEl.className = hasError ? 'code-output code-output-error' : 'code-output'
      outputEl.textContent = hasError ? '编译错误：\n' + (data.Errors || data.errors) : output
    }
  } catch (e) {
    if (outputEl) {
      outputEl.style.display = 'block'
      outputEl.className = 'code-output code-output-error'
      outputEl.textContent = '运行失败：' + e.message
    }
  } finally {
    btn.textContent = '▶ 运行'
    btn.disabled = false
  }
}

async function load() {
  const r = await getPost(route.params.id)
  post.value = r.data.post
  await nextTick()
  injectCodeActions()
}

async function toggleLike() {
  if (!userStore.isLoggedIn) return alert('请先登录')
  if (liked.value) {
    await removeLike({ target_type: 'post', target_id: post.value.id })
    post.value.like_count--
  } else {
    await addLike({ target_type: 'post', target_id: post.value.id })
    post.value.like_count++
  }
  liked.value = !liked.value
}

async function handleDelete() {
  if (!confirm('确定删除？')) return
  await deletePost(route.params.id)
  router.push('/posts')
}

onMounted(load)
</script>

<style>
.page { max-width: 720px; margin: 0 auto; padding: 0 24px; }
.article-title { font-size: 32px; font-weight: 700; color: #18181b; margin-bottom: 16px; }
.article-meta { display: flex; gap: 12px; font-size: 14px; color: #a1a1aa; margin-bottom: 24px; flex-wrap: wrap; }
.tag { background: #f4f4f5; padding: 2px 8px; border-radius: 4px; }
.tag-blue { background: #eff6ff; color: #2563eb; }
.admin-actions { margin-bottom: 24px; display: flex; gap: 8px; }
.btn-small { font-size: 12px; padding: 4px 10px; border: 1px solid #e5e5e5; border-radius: 4px; background: #fff; cursor: pointer; }
.btn-danger { color: #ef4444; }
.article-content { margin-bottom: 32px; }
.article-content h1,.article-content h2,.article-content h3 { font-weight: 700; margin: 1.5em 0 0.5em; }
.article-content h1 { font-size: 1.8em; }
.article-content h2 { font-size: 1.4em; }
.article-content p { line-height: 1.8; margin: 0.8em 0; color: #3f3f46; }
.article-content code { background: #f4f4f5; padding: 2px 6px; border-radius: 4px; font-size: 0.85em; }
.article-content pre { background: #1e1e2e; color: #cdd6f4; padding: 16px 20px; border-radius: 10px; overflow-x: auto; margin: 1em 0; }
.article-content pre code { background: none; color: inherit; padding: 0; font-size: 0.9em; font-family: 'Menlo', 'Monaco', monospace; }
.article-content blockquote { border-left: 3px solid #2563eb; margin: 1em 0; padding: 4px 16px; color: #71717a; background: #f8faff; border-radius: 0 6px 6px 0; }
.article-content img { max-width: 100%; border-radius: 8px; margin: 0.5em 0; }
.article-content table { border-collapse: collapse; width: 100%; margin: 1em 0; }
.article-content th, .article-content td { border: 1px solid #e5e5e5; padding: 8px 12px; }
.article-content th { background: #f4f4f5; font-weight: 600; }
.article-content hr { border: none; border-top: 1px solid #e5e5e5; margin: 1.5em 0; }
.article-content a { color: #2563eb; }
.article-content ul { list-style-type: disc; padding-left: 1.8em; margin: 0.5em 0; }
.article-content ol { list-style-type: decimal; padding-left: 1.8em; margin: 0.5em 0; }
.article-content ul ul { list-style-type: circle; }
.article-content li { margin: 0.3em 0; display: list-item; }
/* 代码块工具栏 */
.article-content .code-toolbar { position: absolute; top: 0; left: 0; right: 0; display: flex; align-items: center; gap: 8px; padding: 6px 12px; background: rgba(255,255,255,0.06); border-bottom: 1px solid rgba(255,255,255,0.08); }
.article-content .code-lang { font-size: 12px; color: #71717a; flex: 1; font-family: 'Menlo', monospace; }
.article-content .code-btn { font-size: 12px; padding: 2px 10px; border: 1px solid rgba(255,255,255,0.15); border-radius: 4px; background: rgba(255,255,255,0.08); color: #a1a1aa; cursor: pointer; transition: all 0.15s; }
.article-content .code-btn:hover { background: rgba(255,255,255,0.15); color: #fff; }
.article-content .code-run-btn { color: #4ade80; border-color: rgba(74,222,128,0.3); }
.article-content .code-run-btn:hover { background: rgba(74,222,128,0.15); color: #4ade80; }
.article-content .code-run-btn:disabled { opacity: 0.5; cursor: not-allowed; }
/* 输出区 */
.article-content .code-output { margin-top: 12px; padding: 12px 16px; background: #0f0f17; border-top: 1px solid rgba(255,255,255,0.08); border-radius: 0 0 8px 8px; font-family: 'Menlo', monospace; font-size: 13px; color: #a3e635; white-space: pre-wrap; word-break: break-all; }
.article-content .code-output-error { color: #f87171; }
.like-bar { padding: 16px 0; border-top: 1px solid #e5e5e5; border-bottom: 1px solid #e5e5e5; margin-bottom: 24px; }
.like-btn { font-size: 14px; padding: 8px 16px; border: 1px solid #e5e5e5; border-radius: 20px; background: #fff; cursor: pointer; }
.like-btn:hover { background: #f4f4f5; }
.like-btn.liked { background: #2563eb; color: #fff; border-color: #2563eb; }
</style>

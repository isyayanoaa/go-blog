import { marked } from 'marked'
import hljs from 'highlight.js'

// 配置 marked 使用 highlight.js 做代码高亮
marked.setOptions({
  highlight(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try { return hljs.highlight(code, { language: lang }).value } catch {}
    }
    return hljs.highlightAuto(code).value
  },
  breaks: true,   // 换行符转 <br>
  gfm: true,      // GitHub Flavored Markdown
})

export { marked }

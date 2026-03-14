<template>
  <div class="page">
    <div class="about-card">
      <div class="about-avatar">👨‍💻</div>
      <h1 class="about-name">Hi, I'm a Go developer</h1>
      <p class="about-desc">这是我的个人博客，记录 Go 学习笔记、技术思考和日常动态。</p>
      <div class="about-links">
        <a href="https://github.com/isyayanoaa" target="_blank" class="about-link">GitHub</a>
      </div>
    </div>
  </div>
</template>

<style>
.page { max-width: 560px; margin: 0 auto; padding: 0 24px; }
.about-card { background: #fff; border: 1px solid #e5e5e5; border-radius: 16px; padding: 40px; text-align: center; }
.about-avatar { width: 80px; height: 80px; border-radius: 50%; background: #f4f4f5; margin: 0 auto 24px; display: flex; align-items: center; justify-content: center; font-size: 40px; }
.about-name { font-size: 24px; font-weight: 700; color: #18181b; margin-bottom: 12px; }
.about-desc { color: #71717a; line-height: 1.6; margin-bottom: 24px; }
.about-links { display: flex; justify-content: center; gap: 24px; }
.about-link { color: #52525b; font-size: 14px; }
.about-link:hover { color: #18181b; }
</style>

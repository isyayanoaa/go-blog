import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/',                    component: () => import('../pages/Home.vue') },
  { path: '/posts',               component: () => import('../pages/Posts.vue') },
  { path: '/posts/:id',           component: () => import('../pages/PostDetail.vue') },
  { path: '/moments',             component: () => import('../pages/Moments.vue') },
  { path: '/categories',          component: () => import('../pages/Categories.vue') },
  { path: '/categories/:name',    component: () => import('../pages/Categories.vue') },
  { path: '/tags',                component: () => import('../pages/Tags.vue') },
  { path: '/tags/:name',          component: () => import('../pages/Tags.vue') },
  { path: '/archives',            component: () => import('../pages/Archives.vue') },
  { path: '/about',               component: () => import('../pages/About.vue') },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})

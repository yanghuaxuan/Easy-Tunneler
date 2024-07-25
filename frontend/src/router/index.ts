import { createRouter, createWebHistory } from 'vue-router'

import 'vue-router'

export {}

declare module 'vue-router' {
  interface RouteMeta {
    transition?: string
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'tunnels',
      component: () => import('@/views/TunnelView.vue'),
      meta: { transition: "slide-right" }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/SettingsView.vue'),
      meta: { transition: 'slide-left' }
    }
  ]
})

export default router

import { createRouter, createWebHistory } from 'vue-router'

import 'vue-router'

// To ensure it is treated as a module, add at least one `export` statement
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
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/SettingsView.vue'),
      meta: { transition: 'slide-left' }
    }
  ]
})

export default router

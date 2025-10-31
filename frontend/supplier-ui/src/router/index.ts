import { nextTick } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import auth from '@/auth'

// https://router.vuejs.org/guide/advanced/meta.html
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth: boolean
  }
}

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/Login.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    redirect: '/home',
    component: () => import('@/layouts/Main.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '/home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: '/products',
        children: [
          {
            path: '',
            name: 'Products',
            component: () => import('@/views/CoreProducts.vue'),
          },
        ],
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  // nav guard: redirect to /login if the page requires auth and the user is not authed
  if (to.matched.some(r => r.meta.requiresAuth) && !auth.user.authenticated) {
    var destPath = '/login'
    if (to.path.length > 1) {
      destPath += '?to=' + to.path
    }
    router.push(destPath)
    return
  }

  next()
})

router.afterEach((to, from) => {
  nextTick(() => {
    document.title = to.name?.toString() + ' | Northwind'
  })
})

export default router

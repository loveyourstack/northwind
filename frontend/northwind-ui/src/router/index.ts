import { nextTick } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

// https://router.vuejs.org/guide/advanced/meta.html
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth: boolean
  }
}

const routes = [
  {
    path: '/',
    redirect: '/home',
    component: () => import('@/layouts/Main.vue'),
    children: [
      {
        path: '/home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: '/categories',
        children: [
          {
            path: '',
            name: 'Categories',
            component: () => import('@/views/Categories.vue'),
          },
        ],
      },
      {
        path: '/customers',
        children: [
          {
            path: '',
            name: 'Customers',
            component: () => import('@/views/Customers.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New customer',
            component: () => import('@/views/CustomerDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Customer detail',
            component: () => import('@/views/CustomerDetail.vue'),
          },
        ],
      },
      {
        path: '/employees',
        children: [
          {
            path: '',
            name: 'Employees',
            component: () => import('@/views/Employees.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New employee',
            component: () => import('@/views/EmployeeDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Employee detail',
            component: () => import('@/views/EmployeeDetail.vue'),
          },
        ]
      },
      {
        path: '/orders',
        children: [
          {
            path: '',
            name: 'Orders',
            component: () => import('@/views/Orders.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New order',
            component: () => import('@/views/OrderDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Order detail',
            component: () => import('@/views/OrderDetail.vue'),
          },
        ]
      },
      {
        path: '/products',
        children: [
          {
            path: '',
            name: 'Products',
            component: () => import('@/views/Products.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New product',
            component: () => import('@/views/ProductDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Product detail',
            component: () => import('@/views/ProductDetail.vue'),
          },
        ],
      },
      {
        path: '/suppliers',
        children: [
          {
            path: '',
            name: 'Suppliers',
            component: () => import('@/views/Suppliers.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New supplier',
            component: () => import('@/views/SupplierDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Supplier detail',
            component: () => import('@/views/SupplierDetail.vue'),
          },
        ],
      },
      {
        path: '/territories',
        children: [
          {
            path: '',
            name: 'Territories',
            component: () => import('@/views/Territories.vue'),
          },
        ]
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

router.afterEach((to, from) => {
  nextTick(() => {
    document.title = to.name?.toString() + ' | Northwind'
  })
})

export default router

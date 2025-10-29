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
            component: () => import('@/views/CoreCategories.vue'),
          },
        ],
      },
      {
        path: '/countries',
        children: [
          {
            path: '',
            name: 'Countries',
            component: () => import('@/views/CoreCountries.vue'),
          },
        ]
      },
      {
        path: '/customers',
        children: [
          {
            path: '',
            name: 'Customers',
            component: () => import('@/views/SalesCustomers.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New customer',
            component: () => import('@/views/SalesCustomerDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Customer detail',
            component: () => import('@/views/SalesCustomerDetail.vue'),
          },
        ],
      },
      {
        path: '/database',
        children: [
          {
            path: '',
            name: 'Database',
            component: () => import('@/views/TechDatabase.vue'),
          },
        ]
      },
      {
        path: '/employees',
        children: [
          {
            path: '',
            name: 'Employees',
            component: () => import('@/views/HrEmployees.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New employee',
            component: () => import('@/views/HrEmployeeDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Employee detail',
            component: () => import('@/views/HrEmployeeDetail.vue'),
          },
        ]
      },
      {
        path: '/meeting-schedule',
        children: [
          {
            path: '',
            name: 'MeetingSchedule',
            component: () => import('@/views/HrMeetingSchedule.vue'),
          },
        ]
      },
      {
        path: '/orders',
        children: [
          {
            path: '',
            name: 'Orders',
            component: () => import('@/views/SalesOrders.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New order',
            component: () => import('@/views/SalesOrderDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Order detail',
            component: () => import('@/views/SalesOrderDetail.vue'),
          },
        ]
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
      {
        path: '/shippers',
        children: [
          {
            path: '',
            name: 'Shippers',
            component: () => import('@/views/SalesShippers.vue'),
          },
        ]
      },
      {
        path: '/suppliers',
        children: [
          {
            path: '',
            name: 'Suppliers',
            component: () => import('@/views/CoreSuppliers.vue'),
          },
          {
            path: 'new',
            props: () => { return { id: 0 } },
            name: 'New supplier',
            component: () => import('@/views/CoreSupplierDetail.vue'),
          },
          {
            path: ':id',
            props: (route: any) => {
              return {
                id: parseInt(route.params.id),
              }
            },
            name: 'Supplier detail',
            component: () => import('@/views/CoreSupplierDetail.vue'),
          },
        ],
      },
      {
        path: '/territories',
        children: [
          {
            path: '',
            name: 'Territories',
            component: () => import('@/views/SalesTerritories.vue'),
          },
        ]
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.afterEach((to, from) => {
  nextTick(() => {
    document.title = to.name?.toString() + ' | Northwind'
  })
})

export default router

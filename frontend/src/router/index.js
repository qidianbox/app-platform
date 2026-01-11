import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue')
      },
      {
        path: 'apps',
        name: 'Apps',
        component: () => import('@/views/app/index.vue')
      },
      {
        path: 'modules',
        name: 'Modules',
        component: () => import('@/views/module/index.vue')
      },
      {
        path: 'system/audit',
        name: 'AuditLog',
        component: () => import('@/views/system/AuditLog.vue'),
        meta: { title: '操作审计日志' }
      }
    ]
  },
  // APP详情页面独立于主布局，拥有自己的顶部导航
  {
    path: '/apps/:id/config',
    name: 'AppConfig',
    component: () => import('@/views/app/config/index.vue')
  },
  // 功能版本管理页面
  {
    path: '/apps/:appId/collections/:collectionId/versions',
    name: 'FeatureVersions',
    component: () => import('@/views/app/config/FeatureVersions.vue'),
    meta: { title: '版本管理' }
  },
  // 模块版本管理页面
  {
    path: '/apps/:appId/modules/:moduleCode/versions',
    name: 'ModuleVersions',
    component: () => import('@/views/app/config/ModuleVersions.vue'),
    meta: { title: '模块版本管理' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router

import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Dashboard from '@/views/Dashboard.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录 | 农产品溯源系统' }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { title: '注册 | 农产品溯源系统' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      requiresAuth: true,
      title: '控制面板 | 农产品溯源系统'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // 设置页面标题
  document.title = to.meta.title || '农产品溯源系统'

  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')

    if (!token) {
      next('/login')
      return
    }

    // 验证令牌是否有效
    try {
      const isValid = await authStore.checkAuth()
      if (!isValid) {
        next('/login')
        return
      }
    } catch (error) {
      console.error('令牌验证失败:', error)
      next('/login')
      return
    }
  }

  // 如果已登录但访问登录页，重定向到仪表板
  if ((to.name === 'Login' || to.name === 'Register') && localStorage.getItem('token')) {
    next('/dashboard')
    return
  }

  next()
})

export default router

import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Dashboard from '@/views/Dashboard.vue'
import Explorer from '@/views/Explorer.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Explorer,
    meta: { title: '区块浏览器 | 农产品溯源系统' }
  },
  {
    path: '/explorer',
    redirect: '/'
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || '农产品溯源系统'

  // 公开页面：根路径、explorer、登录、注册，直接放行
  const publicPaths = ['/', '/explorer', '/login', '/register']
  if (publicPaths.includes(to.path)) {
    next()
    return
  }

  // 只有需要认证的页面才检查
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')

    if (!token) {
      next('/login')
      return
    }

    // 验证令牌是否有效
    const authStore = useAuthStore()
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

  next()
})

export default router

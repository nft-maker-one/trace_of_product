import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()

  // 从 localStorage 恢复用户信息
  const getStoredUser = () => {
    const stored = localStorage.getItem('user')
    if (!stored) return null

    const userData = JSON.parse(stored)
    // 如果旧数据中没有nickname，使用username作为nicknameaine
    if (!userData.nickname && userData.username) {
      userData.nickname = userData.username
    }
    return userData
  }

  const token = ref(localStorage.getItem('token') || '')
  const user = ref(getStoredUser())
  const loading = ref(false)
  const error = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(userData) {
    user.value = userData
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    error.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('rememberUser')
    localStorage.removeItem('username')

    // 重定向到登录页
    router.push('/login')
  }

  async function login(username, password) {
    loading.value = true
    error.value = null

    try {
      const response = await api.login(username, password)
      console.log("登录响应:", response)

      if (response.data.token) {
        setToken(response.data.token)

        // 设置用户信息
        setUser({
          id: response.data.user.id,
          username,
          nickname: response.data.user.nickname,
          email: response.data.user.email,
          avatarUrl: response.data.user.avatar_url,
          loginTime: new Date().toISOString(),
          lastLoginAt: response.data.user.last_login_at,
          createdAt: response.data.user.created_at,
          updatedAt: response.data.user.updated_at
        })
        console.log("用户信息:", user)

        return {
          success: true,
          data: response
        }
      } else {
        error.value = response.message || '登录失败'
        return {
          success: false,
          message: response.message || '登录失败'
        }
      }
    } catch (error) {
      error.value = error.message || '网络错误，请检查连接'
      return {
        success: false,
        message: error.message || '网络错误，请检查连接'
      }
    } finally {
      loading.value = false
    }
  }

  async function register(userData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.register({
        username: userData.username,
        email: userData.email,
        password: userData.password,
        invite_code: userData.inviteCode || null,
        agree_terms: true
      })

      if (response.success || response.code === 0) {
        return {
          success: true,
          message: response.message || '注册成功'
        }
      } else {
        error.value = response.message || '注册失败'
        return {
          success: false,
          message: response.message || '注册失败'
        }
      }
    } catch (error) {
      error.value = error.message || '网络错误，请检查连接'
      return {
        success: false,
        message: error.message || '网络错误，请检查连接'
      }
    } finally {
      loading.value = false
    }
  }

  async function checkAuth() {
    if (!token.value) {
      return false
    }

    try {
      const response = await api.checkAuth()
      return response.status === 'ok'
    } catch (error) {
      // 认证失败，清除令牌
      console.error('令牌验证失败:', error)
      clearAuth()
      return false
    }
  }

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    setToken,
    setUser,
    clearAuth,
    login,
    register,
    checkAuth
  }
})

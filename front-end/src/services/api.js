import axios from 'axios'

const API_BASE_URL = 'http://localhost:8081'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30秒超时
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response) {
      // 服务器返回错误状态码
      switch (error.response.status) {
        case 401:
          // 未授权，清除令牌并跳转到登录页
          localStorage.removeItem('token')
          if (window.location.pathname !== '/login') {
            window.location.href = '/login'
          }
          break
        case 403:
          console.error('访问被拒绝')
          break
        case 404:
          console.error('请求的资源不存在')
          break
        case 500:
          console.error('服务器内部错误')
          break
        default:
          console.error('请求失败:', error.response.status)
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('网络错误，请检查连接')
    } else {
      // 请求配置出错
      console.error('请求配置错误:', error.message)
    }

    return Promise.reject(error)
  }
)

export default {
  // 认证相关
  login(user_name, password) {
    return api.post('/login', {
      user_name: user_name,
      password: password
    })
  },

  checkAuth() {
    return api.get('/menu')
  },

  // 节点相关
  getNodes() {
    return api.get('/nodes')
  },

  // 数据上传
  uploadData(data) {
    return api.post('/upload', data)
  },

  // 数据查询
  queryData(id, node) {
    return api.get(`/message?id=${id}&node=${node}`)
  },

  // 注册
  register(userData) {
    return api.post('/register', {
      username: userData.username,
      email: userData.email,
      password: userData.password,
      invite_code: userData.invite_code,
      agree_terms: userData.agree_terms
    })
  },

// 更新用户资料
  updateProfile(profileData) {
    return api.post('/user/update/profile', profileData)
  },

// 更新用户头像
  updateAvatar(avatarData) {
    return api.post('/user/update/avatar', avatarData)
  },

  // 测试节点连接
  pingNode(addr) {
    return api.get('/node/ping', { params: { addr } })
  },
}

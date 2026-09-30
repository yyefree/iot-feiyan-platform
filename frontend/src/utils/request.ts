import request from '@/utils/request'

const BASE_URL = '/api/v1'

const http = request.create({
  baseURL: BASE_URL,
  timeout: 10000
})

// 请求拦截器
http.interceptors.request.use(
  (config) => {
    // 可以在这里添加token等
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
http.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code === 0) {
      return res
    } else {
      return Promise.reject(new Error(res.message || '请求失败'))
    }
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

export default http

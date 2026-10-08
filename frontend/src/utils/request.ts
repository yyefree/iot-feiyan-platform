import axios from 'axios'
import type { AxiosRequestConfig, AxiosResponse } from 'axios'

// 创建 axios 实例
const service = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 可以在这里添加token等
    return config
  },
  (error: any) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code === 0) {
      return res
    } else {
      return Promise.reject(new Error(res.message || '请求失败'))
    }
  },
  (error: any) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

export default service

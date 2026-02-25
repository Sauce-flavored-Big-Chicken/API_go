import axios from 'axios'
import { useAuthStore } from '../store/auth'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const absoluteUrlPattern = /^https?:\/\//i

export const apiClient = axios.create({
  baseURL,
  timeout: 20000,
})

export function resolveAssetUrl(path: string): string {
  if (!path) return ''
  if (absoluteUrlPattern.test(path)) return path

  const origin = typeof window !== 'undefined' ? window.location.origin : 'http://localhost'
  const apiBase = apiClient.defaults.baseURL
    ? new URL(String(apiClient.defaults.baseURL), origin).toString()
    : origin

  return new URL(path, apiBase).toString()
}

apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export type ApiResponse<T = unknown> = {
  code: number
  msg: string
  data?: T
  total?: number
  token?: string
  [key: string]: unknown
}

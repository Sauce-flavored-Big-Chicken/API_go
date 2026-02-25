import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { api } from './api/service'
import { useAuthStore } from './store/auth'
import './Login.css'

type LoginType = 'login' | 'register' | 'phone'

export default function LoginPage() {
  const navigate = useNavigate()
  const setToken = useAuthStore((s) => s.setToken)
  const [type, setType] = useState<LoginType>('login')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [phone, setPhone] = useState('')
  const [countdown, setCountdown] = useState(0)
  const [smsCode, setSmsCode] = useState('')

  const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    const form = new FormData(e.currentTarget)
    try {
      const { data } = await api.login({
        userName: String(form.get('userName')),
        password: String(form.get('password')),
      })
      if (data?.token) {
        setToken(data.token)
        navigate('/dashboard')
      } else {
        setError(data?.msg || '登录失败')
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.msg || err.message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    const form = new FormData(e.currentTarget)
    try {
      const { data } = await api.register({
        userName: String(form.get('userName')),
        nickName: String(form.get('nickName')),
        password: String(form.get('password')),
        phoneNumber: String(form.get('phoneNumber')),
        sex: String(form.get('sex') || '0'),
      })
      if (data?.code === 200) {
        alert('注册成功，请登录')
        setType('login')
      } else {
        setError(data?.msg || '注册失败')
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.msg || err.message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handlePhoneLogin = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    const form = new FormData(e.currentTarget)
    try {
      const { data } = await api.phoneLogin({
        phone: String(form.get('phone')),
        smsCode: String(form.get('smsCode')),
      })
      if (data?.token) {
        setToken(data.token)
        navigate('/dashboard')
      } else {
        setError(data?.msg || '登录失败')
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.msg || err.message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleGetCode = async () => {
    if (!phone || countdown > 0) return
    try {
      const { data } = await api.smsCode(phone)
      if (data?.code === 200) {
        setSmsCode(data.data || '')
        setCountdown(60)
        const timer = setInterval(() => {
          setCountdown((c) => {
            if (c <= 1) {
              clearInterval(timer)
              return 0
            }
            return c - 1
          })
        }, 1000)
      } else {
        setError(data?.msg || '获取验证码失败')
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.msg || err.message)
      }
    }
  }

  return (
    <div className="login-page">
      <div className="login-card">
        <div className="login-header">
          <h1>Digital Community</h1>
          <p>管理后台</p>
        </div>

        <div className="login-tabs">
          <button
            className={type === 'login' ? 'active' : ''}
            onClick={() => setType('login')}
          >
            账号登录
          </button>
          <button
            className={type === 'phone' ? 'active' : ''}
            onClick={() => setType('phone')}
          >
            手机登录
          </button>
          <button
            className={type === 'register' ? 'active' : ''}
            onClick={() => setType('register')}
          >
            注册
          </button>
        </div>

        {error && <div className="login-error">{error}</div>}

        {type === 'login' && (
          <form className="login-form" onSubmit={handleLogin}>
            <input name="userName" placeholder="用户名" defaultValue="test01" required />
            <input name="password" type="password" placeholder="密码" defaultValue="123456" required />
            <button type="submit" disabled={loading}>
              {loading ? '登录中...' : '登录'}
            </button>
          </form>
        )}

        {type === 'phone' && (
          <form className="login-form" onSubmit={handlePhoneLogin}>
            <input
              name="phone"
              placeholder="手机号"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              required
            />
            <div className="sms-code-row">
              <input
                name="smsCode"
                placeholder="短信验证码"
                value={smsCode}
                onChange={(e) => setSmsCode(e.target.value)}
                required
              />
              <button
                type="button"
                onClick={handleGetCode}
                disabled={!phone || countdown > 0}
                className="sms-code-btn"
              >
                {countdown > 0 ? `${countdown}s` : '获取验证码'}
              </button>
            </div>
            <button type="submit" disabled={loading}>
              {loading ? '登录中...' : '手机登录'}
            </button>
          </form>
        )}

        {type === 'register' && (
          <form className="login-form" onSubmit={handleRegister}>
            <input name="userName" placeholder="用户名" required />
            <input name="nickName" placeholder="昵称" />
            <input name="password" type="password" placeholder="密码" required />
            <input name="phoneNumber" placeholder="手机号" required />
            <select name="sex">
              <option value="0">男</option>
              <option value="1">女</option>
            </select>
            <button type="submit" disabled={loading}>
              {loading ? '注册中...' : '注册'}
            </button>
          </form>
        )}
      </div>
    </div>
  )
}

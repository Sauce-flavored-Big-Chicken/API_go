import { useEffect, useMemo, useState } from 'react'
import { Navigate, NavLink, Route, Routes, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { api, endpointCatalog } from './api/service'
import { apiClient, resolveAssetUrl } from './api/client'
import { useAuthStore } from './store/auth'
import LoginPage from './Login'
import './App.css'

type RunnerState = {
  loading: boolean
  label: string
  result: unknown
  error: string
}

const initialRunnerState: RunnerState = {
  loading: false,
  label: '',
  result: null,
  error: '',
}

function sortById<T>(arr: T[]): T[] {
  return [...arr].sort((a, b) => {
    const idA = (a as Record<string, unknown>).id
    const idB = (b as Record<string, unknown>).id
    if (idA === undefined || idA === null) return 1
    if (idB === undefined || idB === null) return -1
    return Number(idA) - Number(idB)
  })
}

function useRunner() {
  const [state, setState] = useState<RunnerState>(initialRunnerState)

  const run = async (
    label: string,
    request: () => Promise<unknown>,
    onSuccess?: (data: unknown) => void,
  ) => {
    setState({ loading: true, label, result: null, error: '' })
    try {
      const data = await request()
      onSuccess?.(data)
      setState({ loading: false, label, result: data, error: '' })
    } catch (error) {
      let message = '请求失败'
      if (axios.isAxiosError(error)) {
        const msg = (error.response?.data as { msg?: string } | undefined)?.msg
        message = msg || error.message || message
      } else if (error instanceof Error) {
        message = error.message
      }
      setState({ loading: false, label, result: null, error: message })
    }
  }

  return { state, run }
}

function numberFromForm(data: FormData, key: string, fallback: number) {
  const raw = String(data.get(key) || '')
  const value = Number(raw)
  return Number.isFinite(value) && value > 0 ? value : fallback
}

function textFromForm(data: FormData, key: string) {
  return String(data.get(key) || '').trim()
}

function ResponsePane({ state }: { state: RunnerState }) {
  const getStatusBadge = () => {
    if (state.loading) {
      return <span className="status-badge status-loading">加载中</span>
    }
    if (state.error) {
      return <span className="status-badge status-error">失败</span>
    }
    if (state.result) {
      return <span className="status-badge status-success">成功</span>
    }
    return <span className="status-badge status-idle">等待调用</span>
  }

  return (
    <section className="response-panel">
      <div className="response-header">
        <span className="response-title">{state.label || '响应结果'}</span>
        <div className="response-status">{getStatusBadge()}</div>
      </div>
      <div className="response-body">
        {state.loading && <div className="response-empty">请求中...</div>}
        {state.error && <div className="response-error-text">{state.error}</div>}
        {state.result ? (
          <pre><code>{JSON.stringify(state.result, null, 2)}</code></pre>
        ) : null}
        {!state.loading && !state.error && !state.result && (
          <div className="response-empty">调用接口后结果将显示在这里</div>
        )}
      </div>
    </section>
  )
}

function ImageUrlPickerField({
  label,
  name,
  defaultValue = '',
  value: controlledValue,
  onValueChange,
  placeholder,
  required = false,
}: {
  label: string
  name?: string
  defaultValue?: string
  value?: string
  onValueChange?: (value: string) => void
  placeholder: string
  required?: boolean
}) {
  const [innerValue, setInnerValue] = useState(defaultValue)
  const [showPicker, setShowPicker] = useState(false)
  const [images, setImages] = useState<Array<{ name: string; url: string; thumbUrl?: string }>>([])
  const [pageNum, setPageNum] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const pageSize = 24

  useEffect(() => {
    if (controlledValue === undefined) {
      setInnerValue(defaultValue)
    }
  }, [defaultValue, controlledValue])

  const value = controlledValue !== undefined ? controlledValue : innerValue
  const setValue = (next: string) => {
    if (onValueChange) onValueChange(next)
    if (controlledValue === undefined) setInnerValue(next)
  }

  const loadImages = async (nextPage = 1, append = false) => {
    setLoading(true)
    try {
      const res = await api.imageList({ pageNum: nextPage, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        const list = data.data as typeof images
        setImages((prev) => (append ? [...prev, ...list] : list))
        setPageNum(nextPage)
        setTotal(Number(data.total || 0))
      }
    } finally {
      setLoading(false)
    }
  }

  const hasMore = images.length < total

  return (
    <div className="form-group">
      <label>{label}</label>
      <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
        <input name={name} value={value} onChange={(e) => setValue(e.target.value)} placeholder={placeholder} required={required} style={{ flex: 1 }} />
        <button
          type="button"
          className="btn btn-secondary"
          onClick={() => {
            const nextVisible = !showPicker
            setShowPicker(nextVisible)
            if (nextVisible && images.length === 0) {
              void loadImages(1, false)
            }
          }}
        >
          选择图片
        </button>
      </div>
      {value && (
        <div style={{ marginTop: 8 }}>
          <img src={resolveAssetUrl(value)} alt="预览" style={{ maxWidth: 200, maxHeight: 120, borderRadius: 6, objectFit: 'cover' }} />
        </div>
      )}
      {showPicker && (
        <div style={{ marginTop: 8, maxHeight: 220, overflow: 'auto', border: '1px solid #e5e7eb', borderRadius: 8, padding: 8 }}>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 8 }}>
            {images.map((img) => (
              <div
                key={`${name}-${img.url}`}
                onClick={() => {
                  setValue(img.url)
                  setShowPicker(false)
                }}
                style={{ cursor: 'pointer', border: value === img.url ? '2px solid #2563eb' : '1px solid #e5e7eb', borderRadius: 4, overflow: 'hidden' }}
              >
                <img src={resolveAssetUrl(img.thumbUrl || img.url)} alt={img.name} loading="lazy" decoding="async" style={{ width: '100%', height: 50, objectFit: 'cover' }} />
              </div>
            ))}
          </div>
          <div style={{ marginTop: 8, display: 'flex', justifyContent: 'center' }}>
            {hasMore ? (
              <button type="button" className="btn btn-secondary btn-sm" disabled={loading} onClick={() => void loadImages(pageNum + 1, true)}>
                {loading ? '加载中...' : '加载更多'}
              </button>
            ) : (
              <span style={{ fontSize: 12, color: '#9ca3af' }}>已加载全部图片</span>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function DashboardPage() {
  const navigate = useNavigate()
  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>仪表盘</h2>
        <p>欢迎使用数字社区管理后台</p>
      </div>

      <div className="dashboard-stats">
        <div className="stat-card">
          <div className="stat-icon blue">📰</div>
          <div className="stat-content">
            <h4>新闻资讯</h4>
            <p>管理</p>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon green">📢</div>
          <div className="stat-content">
            <h4>公告通知</h4>
            <p>管理</p>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon purple">🎉</div>
          <div className="stat-content">
            <h4>社区活动</h4>
            <p>管理</p>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon orange">👥</div>
          <div className="stat-content">
            <h4>友邻帖子</h4>
            <p>管理</p>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>快捷操作</h3>
        </div>
        <div className="card-body">
          <div className="quick-actions">
            <button className="quick-action-btn" onClick={() => navigate('/user')}>
              <span className="icon">👤</span>
              <span>用户信息</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/news')}>
              <span className="icon">📰</span>
              <span>新闻管理</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/notice')}>
              <span className="icon">📢</span>
              <span>公告管理</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/activity')}>
              <span className="icon">🎉</span>
              <span>活动管理</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/neighbor')}>
              <span className="icon">💬</span>
              <span>友邻帖子</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/upload')}>
              <span className="icon">📁</span>
              <span>文件上传</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/green')}>
              <span className="icon">🌿</span>
              <span>绿动未来</span>
            </button>
            <button className="quick-action-btn" onClick={() => navigate('/playground')}>
              <span className="icon">🧪</span>
              <span>API 测试台</span>
            </button>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>使用说明</h3>
        </div>
        <div className="card-body">
          <p style={{ color: '#6b7280', lineHeight: 1.8 }}>
            该管理端覆盖数字社区文档全部接口，支持以下操作：
          </p>
          <ul style={{ color: '#6b7280', lineHeight: 2, paddingLeft: 20 }}>
            <li>用户信息查看与修改</li>
            <li>新闻、公告的浏览与互动（点赞、评论）</li>
            <li>社区活动的查看、报名、签到、评论</li>
            <li>友邻帖子的浏览与评论</li>
            <li>文件上传功能</li>
            <li>API 测试台可测试所有接口</li>
          </ul>
        </div>
      </div>
    </div>
  )
}

function UserPage() {
  const { state, run } = useRunner()
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [users, setUsers] = useState<unknown[]>([])
  const [showModal, setShowModal] = useState(false)
  const [editingUser, setEditingUser] = useState<unknown | null>(null)
  const [modalType, setModalType] = useState<'add' | 'edit'>('add')

  const loadUsers = () => {
    run('加载用户列表', async () => {
      const res = await api.userList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setUsers(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }

  useEffect(() => {
    run('加载用户列表', async () => {
      const res = await api.userList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setUsers(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }, [pageNum, pageSize])

  const handleAdd = () => {
    setModalType('add')
    setEditingUser(null)
    setShowModal(true)
  }

  const handleEdit = (user: unknown) => {
    setModalType('edit')
    setEditingUser(user)
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除该用户吗？')) return
    run('删除用户', async () => {
      const res = await api.userDelete(id)
      loadUsers()
      return res.data
    })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      userName: textFromForm(form, 'userName'),
      password: textFromForm(form, 'password'),
      nickName: textFromForm(form, 'nickName'),
      phone: textFromForm(form, 'phone'),
      sex: textFromForm(form, 'sex'),
      email: textFromForm(form, 'email'),
      address: textFromForm(form, 'address'),
      introduction: textFromForm(form, 'introduction'),
    }

    if (modalType === 'add') {
      run('创建用户', async () => {
        const res = await api.userCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadUsers()
        }
        return res.data
      })
    } else {
      run('更新用户', async () => {
        const res = await api.userUpdate((editingUser as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadUsers()
        }
        return res.data
      })
    }
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>用户管理</h2>
        <p>管理系统用户，支持增删改查</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>用户列表</h3>
          <div className="card-actions">
            <button className="btn btn-primary" onClick={loadUsers}>刷新</button>
            <button className="btn btn-primary" onClick={handleAdd}>新增用户</button>
          </div>
        </div>
        <div className="card-body">
          <table className="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>用户名</th>
                <th>昵称</th>
                <th>手机号</th>
                <th>邮箱</th>
                <th>性别</th>
                <th>积分</th>
                <th>余额</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {(users as Array<{
                id: number
                userName: string
                nickName: string
                phone: string
                email: string
                sex: string
                score: number
                balance: number
              }>).map((user) => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>{user.userName}</td>
                  <td>{user.nickName}</td>
                  <td>{user.phone}</td>
                  <td>{user.email || '-'}</td>
                  <td>{user.sex === '0' ? '男' : '女'}</td>
                  <td>{user.score}</td>
                  <td>¥{user.balance}</td>
                  <td>
                    <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(user)}>编辑</button>
                    <button className="btn btn-danger btn-sm" onClick={() => handleDelete(user.id)}>删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {total > 0 && (
            <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>当前用户信息</h3>
          <div className="card-actions">
            <button className="btn btn-primary" onClick={() => run('获取用户信息', async () => (await api.getUserInfo()).data)}>
              刷新
            </button>
          </div>
        </div>
        <div className="card-body">
          <form onSubmit={(e) => {
            e.preventDefault()
            const form = new FormData(e.currentTarget)
            run('更新用户信息', async () => (await api.updateUserInfo({
              nickName: textFromForm(form, 'nickName'),
              avatar: textFromForm(form, 'avatar'),
              email: textFromForm(form, 'email'),
              phonenumber: textFromForm(form, 'phonenumber'),
              sex: textFromForm(form, 'sex'),
              address: textFromForm(form, 'address'),
              introduction: textFromForm(form, 'introduction'),
            })).data)
          }}>
            <div className="form-row">
              <div className="form-group">
                <label>昵称</label>
                <input name="nickName" placeholder="请输入昵称" required />
              </div>
              <div className="form-group">
                <label>头像URL</label>
                <input name="avatar" placeholder="头像链接" />
              </div>
            </div>
            <div className="form-row">
              <div className="form-group">
                <label>邮箱</label>
                <input name="email" type="email" placeholder="邮箱" />
              </div>
              <div className="form-group">
                <label>手机号</label>
                <input name="phonenumber" placeholder="手机号" />
              </div>
            </div>
            <div className="form-row">
              <div className="form-group">
                <label>性别</label>
                <select name="sex">
                  <option value="0">男</option>
                  <option value="1">女</option>
                </select>
              </div>
              <div className="form-group">
                <label>地址</label>
                <input name="address" placeholder="地址" />
              </div>
            </div>
            <div className="form-group">
              <label>个人简介</label>
              <textarea name="introduction" placeholder="个人简介" rows={2} />
            </div>
            <button type="submit" className="btn btn-primary">保存修改</button>
          </form>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>修改密码</h3>
        </div>
        <div className="card-body">
          <form onSubmit={(e) => {
            e.preventDefault()
            const form = new FormData(e.currentTarget)
            run('修改密码', async () => (await api.resetPwd({
              oldPassword: textFromForm(form, 'oldPassword'),
              newPassword: textFromForm(form, 'newPassword'),
            })).data)
          }}>
            <div className="form-row">
              <div className="form-group">
                <label>原密码</label>
                <input name="oldPassword" type="password" placeholder="原密码" required />
              </div>
              <div className="form-group">
                <label>新密码</label>
                <input name="newPassword" type="password" placeholder="新密码" required />
              </div>
            </div>
            <button type="submit" className="btn btn-primary">修改密码</button>
          </form>
        </div>
      </div>

      <ResponsePane state={state} />

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{modalType === 'add' ? '新增用户' : '编辑用户'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>×</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="modal-body">
                <div className="form-group">
                  <label>用户名</label>
                  <input
                    name="userName"
                    placeholder="用户名"
                    required
                    defaultValue={(editingUser as { userName?: string })?.userName || ''}
                    disabled={modalType === 'edit'}
                  />
                </div>
                {modalType === 'add' && (
                  <div className="form-group">
                    <label>密码</label>
                    <input name="password" type="password" placeholder="密码" required />
                  </div>
                )}
                <div className="form-group">
                  <label>昵称</label>
                  <input
                    name="nickName"
                    placeholder="昵称"
                    defaultValue={(editingUser as { nickName?: string })?.nickName || ''}
                  />
                </div>
                <div className="form-group">
                  <label>手机号</label>
                  <input
                    name="phone"
                    placeholder="手机号"
                    required
                    defaultValue={(editingUser as { phone?: string })?.phone || ''}
                  />
                </div>
                <div className="form-group">
                  <label>性别</label>
                  <select name="sex" defaultValue={(editingUser as { sex?: string })?.sex || '0'}>
                    <option value="0">男</option>
                    <option value="1">女</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>邮箱</label>
                  <input
                    name="email"
                    type="email"
                    placeholder="邮箱"
                    defaultValue={(editingUser as { email?: string })?.email || ''}
                  />
                </div>
                <div className="form-group">
                  <label>地址</label>
                  <input
                    name="address"
                    placeholder="地址"
                    defaultValue={(editingUser as { address?: string })?.address || ''}
                  />
                </div>
                <div className="form-group">
                  <label>个人简介</label>
                  <textarea
                    name="introduction"
                    placeholder="个人简介"
                    rows={2}
                    defaultValue={(editingUser as { introduction?: string })?.introduction || ''}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

function NewsPage() {
  const { state, run } = useRunner()
  const [activeTab, setActiveTab] = useState<'news' | 'category'>('news')
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [list, setList] = useState<unknown[]>([])
  const [categories, setCategories] = useState<unknown[]>([])
  const [showModal, setShowModal] = useState(false)
  const [editingItem, setEditingItem] = useState<unknown | null>(null)
  const [modalType, setModalType] = useState<'add' | 'edit'>('add')

  const loadNews = () => {
    run('加载新闻列表', async () => {
      const res = await api.pressNewsList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setList(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }

  const loadCategories = () => {
    run('加载分类列表', async () => {
      const res = await api.pressCategoryList()
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setCategories(data.data as unknown[])
      }
      return data
    })
  }

  useEffect(() => {
    if (activeTab === 'news') {
      loadNews()
    } else {
      loadCategories()
    }
  }, [activeTab, pageNum, pageSize])

  const handleAdd = () => {
    setModalType('add')
    setEditingItem(null)
    setShowModal(true)
  }

  const handleEdit = (item: unknown) => {
    setModalType('edit')
    setEditingItem(item)
    setShowModal(true)
  }

  const handleDelete = async (id: number, type: 'news' | 'category') => {
    if (!confirm(`确定要删除这条${type === 'news' ? '新闻' : '分类'}吗？`)) return
    run('删除', async () => {
      const res = type === 'news' 
        ? await api.pressNewsDelete(id)
        : await api.pressCategoryDelete(id)
      if (type === 'news') loadNews()
      else loadCategories()
      return res.data
    })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      title: textFromForm(form, 'title'),
      subTitle: textFromForm(form, 'subTitle'),
      content: textFromForm(form, 'content'),
      categoryId: numberFromForm(form, 'categoryId', 0),
      type: textFromForm(form, 'type'),
      imageUrls: textFromForm(form, 'imageUrls'),
    }

    if (modalType === 'add') {
      run('创建新闻', async () => {
        const res = await api.pressNewsCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadNews()
        }
        return res.data
      })
    } else {
      run('更新新闻', async () => {
        const res = await api.pressNewsUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadNews()
        }
        return res.data
      })
    }
  }

  const handleCategorySubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      name: textFromForm(form, 'name'),
      sort: numberFromForm(form, 'sort', 0),
      status: textFromForm(form, 'status'),
    }

    if (modalType === 'add') {
      run('创建分类', async () => {
        const res = await api.pressCategoryCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadCategories()
        }
        return res.data
      })
    } else {
      run('更新分类', async () => {
        const res = await api.pressCategoryUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadCategories()
        }
        return res.data
      })
    }
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>新闻资讯</h2>
        <p>管理新闻内容和分类</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>管理</h3>
          <div className="card-actions">
            <button className={`btn ${activeTab === 'news' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setActiveTab('news')}>新闻列表</button>
            <button className={`btn ${activeTab === 'category' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setActiveTab('category')}>分类管理</button>
          </div>
        </div>
      </div>

      {activeTab === 'news' && (
        <div className="card">
          <div className="card-header">
            <h3>新闻列表</h3>
            <div className="card-actions">
              <button className="btn btn-secondary" onClick={loadNews}>刷新</button>
              <button className="btn btn-primary" onClick={handleAdd}>新增新闻</button>
            </div>
          </div>
          <div className="card-body">
            <table className="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>标题</th>
                  <th>副标题</th>
                  <th>分类</th>
                  <th>点赞数</th>
                  <th>阅读数</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {(list as Array<{
                  id: number
                  title: string
                  subTitle: string
                  type: string
                  likeNum: number
                  readNum: number
                }>).map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.title}</td>
                    <td>{item.subTitle || '-'}</td>
                    <td>{item.type || '-'}</td>
                    <td>{item.likeNum}</td>
                    <td>{item.readNum}</td>
                    <td>
                      <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(item.id, 'news')}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {total > 0 && (
              <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
                <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
                <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
                <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'category' && (
        <div className="card">
          <div className="card-header">
            <h3>分类列表</h3>
            <div className="card-actions">
              <button className="btn btn-secondary" onClick={loadCategories}>刷新</button>
              <button className="btn btn-primary" onClick={handleAdd}>新增分类</button>
            </div>
          </div>
          <div className="card-body">
            <table className="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>分类名称</th>
                  <th>排序</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {(categories as Array<{
                  id: number
                  name: string
                  sort: number
                }>).map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.name}</td>
                    <td>{item.sort}</td>
                    <td>
                      <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(item.id, 'category')}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <ResponsePane state={state} />

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{modalType === 'add' ? `新增${activeTab === 'news' ? '新闻' : '分类'}` : `编辑${activeTab === 'news' ? '新闻' : '分类'}`}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>×</button>
            </div>
            {activeTab === 'news' ? (
              <form onSubmit={handleSubmit}>
                <div className="modal-body">
                  <div className="form-group">
                    <label>标题</label>
                    <input name="title" placeholder="新闻标题" required defaultValue={(editingItem as { title?: string })?.title || ''} />
                  </div>
                  <div className="form-group">
                    <label>副标题</label>
                    <input name="subTitle" placeholder="副标题" defaultValue={(editingItem as { subTitle?: string })?.subTitle || ''} />
                  </div>
                  <div className="form-group">
                    <label>分类ID</label>
                    <input name="categoryId" type="number" placeholder="分类ID" defaultValue={(editingItem as { type?: string })?.type || ''} />
                  </div>
                  <ImageUrlPickerField
                    label="图片URL"
                    name="imageUrls"
                    placeholder="图片URL"
                    defaultValue={(editingItem as { cover?: string })?.cover || ''}
                  />
                  <div className="form-group">
                    <label>内容</label>
                    <textarea name="content" placeholder="新闻内容" rows={4} required defaultValue={(editingItem as { content?: string })?.content || ''} />
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                  <button type="submit" className="btn btn-primary">提交</button>
                </div>
              </form>
            ) : (
              <form onSubmit={handleCategorySubmit}>
                <div className="modal-body">
                  <div className="form-group">
                    <label>分类名称</label>
                    <input name="name" placeholder="分类名称" required defaultValue={(editingItem as { name?: string })?.name || ''} />
                  </div>
                  <div className="form-group">
                    <label>排序</label>
                    <input name="sort" type="number" placeholder="排序" defaultValue={(editingItem as { sort?: number })?.sort || 0} />
                  </div>
                  <div className="form-group">
                    <label>状态</label>
                    <select name="status" defaultValue={(editingItem as { status?: string })?.status || '0'}>
                      <option value="0">正常</option>
                      <option value="1">禁用</option>
                    </select>
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                  <button type="submit" className="btn btn-primary">提交</button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function NoticePage() {
  const { state, run } = useRunner()
  const [activeTab, setActiveTab] = useState<'notice' | 'rotation'>('notice')
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [list, setList] = useState<unknown[]>([])
  const [rotationList, setRotationList] = useState<unknown[]>([])
  const [rotationType, setRotationType] = useState('1')
  const [showModal, setShowModal] = useState(false)
  const [editingItem, setEditingItem] = useState<unknown | null>(null)
  const [modalType, setModalType] = useState<'add' | 'edit'>('add')

  const loadNotices = () => {
    run('加载公告列表', async () => {
      const res = await api.noticeList({ pageNum, pageSize, noticeStatus: '' })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setList(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }

  const loadRotations = () => {
    run('加载轮播图列表', async () => {
      const res = await api.rotationList({ pageNum: 1, pageSize: 100, type: rotationType })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setRotationList(data.data as unknown[])
      }
      return data
    })
  }

  useEffect(() => {
    if (activeTab === 'notice') {
      loadNotices()
    } else {
      loadRotations()
    }
  }, [activeTab, pageNum, pageSize, rotationType])

  const handleAdd = () => {
    setModalType('add')
    setEditingItem(null)
    setShowModal(true)
  }

  const handleEdit = (item: unknown) => {
    setModalType('edit')
    setEditingItem(item)
    setShowModal(true)
  }

  const handleDeleteNotice = async (id: number) => {
    if (!confirm('确定要删除这条公告吗？')) return
    run('删除公告', async () => {
      const res = await api.noticeDelete(id)
      loadNotices()
      return res.data
    })
  }

  const handleDeleteRotation = async (id: number) => {
    if (!confirm('确定要删除这个轮播图吗？')) return
    run('删除轮播图', async () => {
      const res = await api.rotationDelete(id)
      loadRotations()
      return res.data
    })
  }

  const handleNoticeSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      title: textFromForm(form, 'title'),
      noticeContent: textFromForm(form, 'noticeContent'),
      noticeStatus: textFromForm(form, 'noticeStatus'),
      createBy: textFromForm(form, 'createBy'),
    }

    if (modalType === 'add') {
      run('创建公告', async () => {
        const res = await api.noticeCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadNotices()
        }
        return res.data
      })
    } else {
      run('更新公告', async () => {
        const res = await api.noticeUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadNotices()
        }
        return res.data
      })
    }
  }

  const handleRotationSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      title: textFromForm(form, 'title'),
      picPath: textFromForm(form, 'picPath'),
      link: textFromForm(form, 'link'),
      type: parseInt(textFromForm(form, 'type') || '1'),
    }

    if (modalType === 'add') {
      run('创建轮播图', async () => {
        const res = await api.rotationCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadRotations()
        }
        return res.data
      })
    } else {
      run('更新轮播图', async () => {
        const res = await api.rotationUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadRotations()
        }
        return res.data
      })
    }
  }

  const handleTabChange = (tab: 'notice' | 'rotation') => {
    setActiveTab(tab)
    setPageNum(1)
  }

  const handleTypeChange = (type: string) => {
    setRotationType(type)
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>公告管理</h2>
        <p>管理公告通知和轮播图</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>管理</h3>
          <div className="card-actions">
            <button className={`btn ${activeTab === 'notice' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => handleTabChange('notice')}>公告列表</button>
            <button className={`btn ${activeTab === 'rotation' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => handleTabChange('rotation')}>轮播图</button>
          </div>
        </div>
      </div>

      {activeTab === 'notice' && (
        <div className="card">
          <div className="card-header">
            <h3>公告列表</h3>
            <div className="card-actions">
              <button className="btn btn-secondary" onClick={loadNotices}>刷新</button>
              <button className="btn btn-primary" onClick={handleAdd}>新增公告</button>
            </div>
          </div>
          <div className="card-body">
            <table className="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>标题</th>
                  <th>发布单位</th>
                  <th>发布时间</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {(list as Array<{
                  id: number
                  noticeTitle: string
                  releaseUnit: string
                  createTime: string
                  noticeStatus: string
                }>).map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.noticeTitle}</td>
                    <td>{item.releaseUnit || '-'}</td>
                    <td>{item.createTime || '-'}</td>
                    <td>{item.noticeStatus === '1' ? '已发布' : '未发布'}</td>
                    <td>
                      <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDeleteNotice(item.id)}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {total > 0 && (
              <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
                <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
                <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
                <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'rotation' && (
        <div className="card">
          <div className="card-header">
            <h3>轮播图列表</h3>
            <div className="card-actions">
              <select value={rotationType} onChange={(e) => handleTypeChange(e.target.value)} style={{ padding: '6px 10px', borderRadius: 6, border: '1px solid #ddd' }}>
                <option value="1">类型1</option>
                <option value="2">类型2</option>
              </select>
              <button className="btn btn-secondary" onClick={loadRotations}>刷新</button>
              <button className="btn btn-primary" onClick={handleAdd}>新增轮播图</button>
            </div>
          </div>
          <div className="card-body">
            <table className="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>标题</th>
                  <th>图片</th>
                  <th>类型</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {(rotationList as Array<{
                  id: number
                  advTitle: string
                  advImg: string
                  type: string
                }>).map((item) => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td>{item.advTitle}</td>
                    <td>{item.advImg ? <img src={resolveAssetUrl(item.advImg)} alt="" style={{ width: 60, height: 30, objectFit: 'cover' }} /> : '-'}</td>
                    <td>{item.type === '1' ? '类型1' : '类型2'}</td>
                    <td>
                      <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDeleteRotation(item.id)}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <ResponsePane state={state} />

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>
                {activeTab === 'notice' 
                  ? (modalType === 'add' ? '新增公告' : '编辑公告')
                  : (modalType === 'add' ? '新增轮播图' : '编辑轮播图')
                }
              </h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>×</button>
            </div>
            {activeTab === 'notice' ? (
              <form onSubmit={handleNoticeSubmit}>
                <div className="modal-body">
                  <div className="form-group">
                    <label>标题</label>
                    <input name="title" placeholder="公告标题" required defaultValue={(editingItem as { noticeTitle?: string })?.noticeTitle || ''} />
                  </div>
                  <div className="form-group">
                    <label>发布单位</label>
                    <input name="createBy" placeholder="发布单位" defaultValue={(editingItem as { releaseUnit?: string })?.releaseUnit || ''} />
                  </div>
                  <div className="form-group">
                    <label>状态</label>
                    <select name="noticeStatus" defaultValue={(editingItem as { noticeStatus?: string })?.noticeStatus || '0'}>
                      <option value="0">未发布</option>
                      <option value="1">已发布</option>
                    </select>
                  </div>
                  <div className="form-group">
                    <label>内容</label>
                    <textarea name="noticeContent" placeholder="公告内容" rows={4} required defaultValue={(editingItem as { contentNotice?: string })?.contentNotice || ''} />
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                  <button type="submit" className="btn btn-primary">提交</button>
                </div>
              </form>
            ) : (
              <form onSubmit={handleRotationSubmit}>
                <div className="modal-body">
                  <div className="form-group">
                    <label>标题</label>
                    <input name="title" placeholder="轮播图标题" required defaultValue={(editingItem as { advTitle?: string })?.advTitle || ''} />
                  </div>
                  <ImageUrlPickerField
                    label="图片URL"
                    name="picPath"
                    placeholder="图片URL"
                    defaultValue={(editingItem as { advImg?: string })?.advImg || ''}
                  />
                  <div className="form-group">
                    <label>链接</label>
                    <input name="link" placeholder="跳转链接" defaultValue={(editingItem as { link?: string })?.link || ''} />
                  </div>
                  <div className="form-group">
                    <label>类型</label>
                    <select name="type" defaultValue={(editingItem as { type?: string })?.type || rotationType}>
                      <option value="1">类型1</option>
                      <option value="2">类型2</option>
                    </select>
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                  <button type="submit" className="btn btn-primary">提交</button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function NeighborPage() {
  const { state, run } = useRunner()
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [list, setList] = useState<unknown[]>([])
  const [showModal, setShowModal] = useState(false)
  const [editingItem, setEditingItem] = useState<unknown | null>(null)
  const [modalType, setModalType] = useState<'add' | 'edit'>('add')
  const [showCommentModal, setShowCommentModal] = useState(false)
  const [commentTargetId, setCommentTargetId] = useState<number | null>(null)
  const [commentContent, setCommentContent] = useState('')
  const [images, setImages] = useState<Array<{ name: string; url: string; thumbUrl?: string }>>([])
  const [imagePageNum, setImagePageNum] = useState(1)
  const [imageTotal, setImageTotal] = useState(0)
  const [imageLoading, setImageLoading] = useState(false)
  const [imageLoadingAll, setImageLoadingAll] = useState(false)
  const [showImagePicker, setShowImagePicker] = useState(false)
  const [pickerTarget, setPickerTarget] = useState<'post' | 'avatar'>('post')
  const [tempImgUrl, setTempImgUrl] = useState('')
  const [tempUserImgUrl, setTempUserImgUrl] = useState('')
  const [tempNickName, setTempNickName] = useState('')

  const loadList = () => {
    run('加载帖子列表', async () => {
      const res = await api.neighborList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setList(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }

  const imagePageSize = 24

  const mergeImages = (
    prev: Array<{ name: string; url: string; thumbUrl?: string }>,
    next: Array<{ name: string; url: string; thumbUrl?: string }>,
  ) => {
    const map = new Map(prev.map((item) => [item.url, item]))
    for (const item of next) {
      map.set(item.url, item)
    }
    return Array.from(map.values())
  }

  const loadImages = async (nextPage = 1, append = false) => {
    setImageLoading(true)
    try {
      const res = await api.imageList({ pageNum: nextPage, pageSize: imagePageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        const nextList = data.data as typeof images
        setImages((prev) => (append ? mergeImages(prev, nextList) : nextList))
        setImagePageNum(nextPage)
        setImageTotal(Number(data.total || 0))
        return { loadedCount: nextList.length, total: Number(data.total || 0) }
      }
      return { loadedCount: 0, total: 0 }
    } finally {
      setImageLoading(false)
    }
  }

  useEffect(() => {
    loadList()
  }, [pageNum, pageSize])

  useEffect(() => {
    if (showModal) {
      if (images.length === 0) {
        void loadImages(1, false)
      }
      setShowImagePicker(false)
      if (editingItem) {
        setTempNickName((editingItem as { publishName?: string })?.publishName || '')
        setTempImgUrl((editingItem as { imgUrl?: string })?.imgUrl || '')
        setTempUserImgUrl((editingItem as { userImgUrl?: string })?.userImgUrl || '')
      } else {
        setTempNickName('')
        setTempImgUrl('')
        setTempUserImgUrl('')
      }
    }
  }, [showModal, editingItem, images.length])

  const handleAdd = () => {
    setModalType('add')
    setEditingItem(null)
    setShowModal(true)
  }

  const handleEdit = (item: unknown) => {
    setModalType('edit')
    setEditingItem(item)
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条帖子吗？')) return
    run('删除帖子', async () => {
      const res = await api.neighborDelete(id)
      loadList()
      return res.data
    })
  }

  const handleAddComment = (id: number) => {
    setCommentTargetId(id)
    setCommentContent('')
    setShowCommentModal(true)
  }

  const submitComment = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!commentTargetId || !commentContent.trim()) return
    run('发布评论', async () => {
      const res = await api.neighborAddComment({ neighborhoodId: commentTargetId, content: commentContent.trim() })
      if (res.data?.code === 200) {
        setShowCommentModal(false)
        setCommentTargetId(null)
        setCommentContent('')
        loadList()
      }
      return res.data
    })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      content: textFromForm(form, 'content'),
      nickName: tempNickName,
      imgUrl: tempImgUrl,
      userImgUrl: tempUserImgUrl,
    }

    if (modalType === 'add') {
      run('创建帖子', async () => {
        const res = await api.neighborCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadList()
        }
        return res.data
      })
    } else {
      run('更新帖子', async () => {
        const res = await api.neighborUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadList()
        }
        return res.data
      })
    }
  }

  const selectImage = (url: string) => {
    if (pickerTarget === 'avatar') {
      setTempUserImgUrl(url)
    } else {
      setTempImgUrl(url)
    }
    setShowImagePicker(false)
  }

  const loadAllImages = async () => {
    if (imageLoadingAll || imageLoading || !hasMoreImages) return
    setImageLoadingAll(true)
    try {
      let nextPage = imagePageNum
      for (let i = 0; i < 20; i++) {
        const result = await loadImages(nextPage + 1, true)
        nextPage += 1
        if (!result || result.loadedCount < imagePageSize || nextPage * imagePageSize >= result.total) {
          break
        }
      }
    } finally {
      setImageLoadingAll(false)
    }
  }

  const hasMoreImages = images.length < imageTotal

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>友邻帖子</h2>
        <p>管理友邻社区帖子</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>帖子列表</h3>
          <div className="card-actions">
            <button className="btn btn-secondary" onClick={loadList}>刷新</button>
            <button className="btn btn-primary" onClick={handleAdd}>新增帖子</button>
          </div>
        </div>
        <div className="card-body">
          <table className="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>用户</th>
                <th>内容</th>
                <th>图片</th>
                <th>点赞数</th>
                <th>评论数</th>
                <th>发布时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {(list as Array<{
                id: number
                publishName: string
                userImgUrl: string
                publishContent: string
                imgUrl: string
                likeNum: number
                commentNum: number
                publishTime: string
              }>).map((item) => (
                <tr key={item.id}>
                  <td>{item.id}</td>
                  <td>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      {item.userImgUrl ? (
                        <img src={resolveAssetUrl(item.userImgUrl)} alt="头像" style={{ width: 28, height: 28, borderRadius: '50%', objectFit: 'cover' }} />
                      ) : (
                        <div style={{ width: 28, height: 28, borderRadius: '50%', background: '#e5e7eb' }} />
                      )}
                      <span>{item.publishName || '-'}</span>
                    </div>
                  </td>
                  <td style={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{item.publishContent}</td>
                  <td>{item.imgUrl ? <img src={resolveAssetUrl(item.imgUrl)} alt="" style={{ width: 40, height: 40, objectFit: 'cover' }} /> : '-'}</td>
                  <td>{item.likeNum}</td>
                  <td>{item.commentNum}</td>
                  <td>{item.publishTime || '-'}</td>
                  <td>
                    <button className="btn btn-secondary btn-sm" onClick={() => handleAddComment(item.id)}>评论</button>
                    <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                    <button className="btn btn-danger btn-sm" onClick={() => handleDelete(item.id)}>删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {total > 0 && (
            <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <ResponsePane state={state} />

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{modalType === 'add' ? '新增帖子' : '编辑帖子'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>×</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="modal-body">
                <div className="form-group">
                  <label>用户昵称</label>
                  <input
                    value={tempNickName}
                    onChange={(e) => setTempNickName(e.target.value)}
                    placeholder="用户昵称"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>内容</label>
                  <textarea name="content" placeholder="帖子内容" rows={4} required defaultValue={(editingItem as { publishContent?: string })?.publishContent || ''} />
                </div>
                <div className="form-group">
                  <label>图片</label>
                  <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                    <input 
                      value={tempImgUrl} 
                      onChange={(e) => setTempImgUrl(e.target.value)} 
                      placeholder="图片URL" 
                      style={{ flex: 1 }}
                    />
                    <button
                      type="button"
                      className="btn btn-secondary"
                      onClick={() => {
                        setPickerTarget('post')
                        const nextVisible = !showImagePicker
                        setShowImagePicker(nextVisible)
                        if (nextVisible && images.length === 0) {
                          void loadImages(1, false)
                        }
                      }}
                    >
                      选择图片
                    </button>
                  </div>
                  {tempImgUrl && (
                    <div style={{ marginTop: 8 }}>
                      <img src={resolveAssetUrl(tempImgUrl)} alt="预览" style={{ maxWidth: 200, maxHeight: 150, borderRadius: 8 }} />
                    </div>
                  )}
                  {showImagePicker && (
                    <div style={{ marginTop: 8, maxHeight: 200, overflow: 'auto', border: '1px solid #e5e7eb', borderRadius: 8, padding: 8 }}>
                      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 8 }}>
                        {images.map((img) => (
                          <div 
                            key={img.url} 
                            onClick={() => selectImage(img.url)}
                            style={{ cursor: 'pointer', border: (pickerTarget === 'post' ? tempImgUrl : tempUserImgUrl) === img.url ? '2px solid #2563eb' : '1px solid #e5e7eb', borderRadius: 4, overflow: 'hidden' }}
                          >
                            <img src={resolveAssetUrl(img.thumbUrl || img.url)} alt={img.name} loading="lazy" decoding="async" style={{ width: '100%', height: 50, objectFit: 'cover' }} />
                          </div>
                        ))}
                      </div>
                      <div style={{ marginTop: 8, display: 'flex', justifyContent: 'center', gap: 8 }}>
                        {hasMoreImages ? (
                          <>
                            <button
                              type="button"
                              className="btn btn-secondary btn-sm"
                              disabled={imageLoading || imageLoadingAll}
                              onClick={() => void loadImages(imagePageNum + 1, true)}
                            >
                              {imageLoading ? '加载中...' : '加载更多'}
                            </button>
                            <button
                              type="button"
                              className="btn btn-secondary btn-sm"
                              disabled={imageLoading || imageLoadingAll}
                              onClick={() => void loadAllImages()}
                            >
                              {imageLoadingAll ? '加载中...' : '加载全部'}
                            </button>
                          </>
                        ) : (
                          <span style={{ fontSize: 12, color: '#9ca3af' }}>已加载全部图片</span>
                        )}
                      </div>
                    </div>
                  )}
                </div>
                <div className="form-group">
                  <label>用户头像</label>
                  <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                    <input
                      value={tempUserImgUrl}
                      onChange={(e) => setTempUserImgUrl(e.target.value)}
                      placeholder="用户头像URL"
                      style={{ flex: 1 }}
                    />
                    <button
                      type="button"
                      className="btn btn-secondary"
                      onClick={() => {
                        setPickerTarget('avatar')
                        const nextVisible = !showImagePicker
                        setShowImagePicker(nextVisible)
                        if (nextVisible && images.length === 0) {
                          void loadImages(1, false)
                        }
                      }}
                    >
                      选择头像
                    </button>
                  </div>
                  {tempUserImgUrl && (
                    <div style={{ marginTop: 8 }}>
                      <img src={resolveAssetUrl(tempUserImgUrl)} alt="头像预览" style={{ width: 60, height: 60, borderRadius: '50%', objectFit: 'cover' }} />
                    </div>
                  )}
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showCommentModal && (
        <div className="modal-overlay" onClick={() => setShowCommentModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>发布评论</h3>
              <button className="modal-close" onClick={() => setShowCommentModal(false)}>×</button>
            </div>
            <form onSubmit={submitComment}>
              <div className="modal-body">
                <div className="form-group">
                  <label>评论内容</label>
                  <textarea
                    value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    placeholder="输入评论内容"
                    rows={4}
                    required
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowCommentModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">发布</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

function ActivityPage() {
  const { state, run } = useRunner()
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [list, setList] = useState<unknown[]>([])
  const [showModal, setShowModal] = useState(false)
  const [editingItem, setEditingItem] = useState<unknown | null>(null)
  const [modalType, setModalType] = useState<'add' | 'edit'>('add')

  const loadList = () => {
    run('加载活动列表', async () => {
      const res = await api.activityList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setList(sortById(data.data as unknown[]))
        setTotal(data.total as number)
      }
      return data
    })
  }

  useEffect(() => {
    loadList()
  }, [pageNum, pageSize])

  const handleAdd = () => {
    setModalType('add')
    setEditingItem(null)
    setShowModal(true)
  }

  const handleEdit = (item: unknown) => {
    setModalType('edit')
    setEditingItem(item)
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个活动吗？')) return
    run('删除活动', async () => {
      const res = await api.activityDelete(id)
      loadList()
      return res.data
    })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const payload = {
      title: textFromForm(form, 'title'),
      content: textFromForm(form, 'content'),
      picPath: textFromForm(form, 'picPath'),
      categoryId: numberFromForm(form, 'categoryId', 0),
      address: textFromForm(form, 'address'),
      totalCount: numberFromForm(form, 'totalCount', 0),
      isTop: textFromForm(form, 'isTop'),
      createBy: textFromForm(form, 'createBy'),
    }

    if (modalType === 'add') {
      run('创建活动', async () => {
        const res = await api.activityCreate(payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadList()
        }
        return res.data
      })
    } else {
      run('更新活动', async () => {
        const res = await api.activityUpdate((editingItem as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowModal(false)
          loadList()
        }
        return res.data
      })
    }
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>社区活动</h2>
        <p>管理社区活动</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>活动列表</h3>
          <div className="card-actions">
            <button className="btn btn-secondary" onClick={loadList}>刷新</button>
            <button className="btn btn-primary" onClick={handleAdd}>新增活动</button>
          </div>
        </div>
        <div className="card-body">
          <table className="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>标题</th>
                <th>分类</th>
                <th>地址</th>
                <th>报名人数</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {(list as Array<{
                id: number
                title: string
                category: string
                address: string
                signUpNum: number
                maxNum: number
                isTop: string
              }>).map((item) => (
                <tr key={item.id}>
                  <td>{item.id}</td>
                  <td>{item.isTop === '1' && <span style={{color: '#eab308', marginRight: 4}}>置顶</span>}{item.title}</td>
                  <td>{item.category || '-'}</td>
                  <td>{item.address || '-'}</td>
                  <td>{item.signUpNum}/{item.maxNum || '不限'}</td>
                  <td>{item.isTop === '1' ? '置顶' : '正常'}</td>
                  <td>
                    <button className="btn btn-secondary btn-sm" onClick={() => handleEdit(item)}>编辑</button>
                    <button className="btn btn-danger btn-sm" onClick={() => handleDelete(item.id)}>删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {total > 0 && (
            <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <ResponsePane state={state} />

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{modalType === 'add' ? '新增活动' : '编辑活动'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>×</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="modal-body">
                <div className="form-group">
                  <label>标题</label>
                  <input name="title" placeholder="活动标题" required defaultValue={(editingItem as { title?: string })?.title || ''} />
                </div>
                <div className="form-group">
                  <label>分类ID</label>
                  <input name="categoryId" type="number" placeholder="分类ID" defaultValue={(editingItem as { category?: string })?.category || ''} />
                </div>
                <div className="form-group">
                  <label>地址</label>
                  <input name="address" placeholder="活动地址" defaultValue={(editingItem as { address?: string })?.address || ''} />
                </div>
                <div className="form-group">
                  <label>人数上限</label>
                  <input name="totalCount" type="number" placeholder="0表示不限" defaultValue={(editingItem as { maxNum?: number })?.maxNum || ''} />
                </div>
                <ImageUrlPickerField
                  label="图片URL"
                  name="picPath"
                  placeholder="图片URL"
                  defaultValue={(editingItem as { picPath?: string })?.picPath || ''}
                />
                <div className="form-group">
                  <label>主办方</label>
                  <input name="createBy" placeholder="主办方" defaultValue={(editingItem as { sponsor?: string })?.sponsor || ''} />
                </div>
                <div className="form-group">
                  <label>是否置顶</label>
                  <select name="isTop" defaultValue={(editingItem as { isTop?: string })?.isTop || '0'}>
                    <option value="0">否</option>
                    <option value="1">是</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>内容</label>
                  <textarea name="content" placeholder="活动内容" rows={4} required defaultValue={(editingItem as { content?: string })?.content || ''} />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

function UploadPage() {
  const { state, run } = useRunner()
  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>文件上传</h2>
        <p>上传图片、文档等文件</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>上传文件</h3>
        </div>
        <div className="card-body">
          <form onSubmit={(e) => {
            e.preventDefault()
            const form = new FormData(e.currentTarget)
            const file = form.get('file') as File | null
            if (!file || !file.name) return
            run('上传文件', async () => (await api.upload(file)).data)
          }}>
            <div className="form-group">
              <label>选择文件</label>
              <input name="file" type="file" required />
            </div>
            <button type="submit" className="btn btn-primary">上传文件</button>
          </form>
        </div>
      </div>

      <ResponsePane state={state} />
    </div>
  )
}

function ImagePage() {
  const { state, run } = useRunner()
  const [images, setImages] = useState<Array<{ name: string; url: string; thumbUrl?: string; size: number; created: string }>>([])
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(24)
  const [total, setTotal] = useState(0)

  const loadImages = () => {
    run('加载图片列表', async () => {
      const res = await api.imageList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setImages(sortById(data.data as typeof images))
        setTotal(Number(data.total || 0))
      }
      return data
    })
  }

  useEffect(() => {
    loadImages()
  }, [pageNum, pageSize])

  const handleDelete = async (url: string) => {
    if (!confirm('确定要删除这张图片吗？')) return
    run('删除图片', async () => {
      const res = await api.imageDelete(url)
      loadImages()
      return res.data
    })
  }

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>图片管理</h2>
        <p>管理所有上传的图片</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>图片列表</h3>
          <div className="card-actions">
            <button className="btn btn-secondary" onClick={loadImages}>刷新</button>
          </div>
        </div>
        <div className="card-body">
          {images.length === 0 ? (
            <div className="empty-state">
              <div className="icon">🖼️</div>
              <p>暂无图片</p>
            </div>
          ) : (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))', gap: 16 }}>
              {images.map((img) => (
                <div key={img.url} style={{ border: '1px solid #e5e7eb', borderRadius: 8, overflow: 'hidden' }}>
                  <div style={{ height: 140, overflow: 'hidden', background: '#f9fafb', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    <img src={resolveAssetUrl(img.thumbUrl || img.url)} alt={img.name} loading="lazy" decoding="async" style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }} />
                  </div>
                  <div style={{ padding: 8, borderTop: '1px solid #e5e7eb' }}>
                    <p style={{ margin: 0, fontSize: 12, color: '#6b7280', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{img.name}</p>
                    <p style={{ margin: '4px 0 0', fontSize: 11, color: '#9ca3af' }}>{formatSize(img.size)}</p>
                    <button className="btn btn-danger btn-sm" style={{ marginTop: 8, width: '100%' }} onClick={() => handleDelete(img.url)}>删除</button>
                  </div>
                </div>
              ))}
            </div>
          )}
          {total > 0 && (
            <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {pageNum} / {Math.max(1, Math.ceil(total / pageSize))} 页，共 {total} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={pageNum >= Math.ceil(total / pageSize)} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <ResponsePane state={state} />
    </div>
  )
}

function FilePage() {
  const { state, run } = useRunner()
  const [files, setFiles] = useState<Array<{ name: string; url: string; size: number; created: string }>>([])
  const [pageNum, setPageNum] = useState(1)
  const [pageSize] = useState(20)
  const [total, setTotal] = useState(0)

  const loadFiles = () => {
    run('加载文件列表', async () => {
      const res = await api.fileList({ pageNum, pageSize })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setFiles(sortById(data.data as typeof files))
        setTotal(Number(data.total || 0))
      }
      return data
    })
  }

  useEffect(() => {
    loadFiles()
  }, [pageNum, pageSize])

  const handleDelete = async (url: string) => {
    if (!confirm('确定要删除这个文件吗？')) return
    run('删除文件', async () => {
      const res = await api.fileDelete(url)
      loadFiles()
      return res.data
    })
  }

  const handleDownload = (url: string) => {
    window.open(resolveAssetUrl(url), '_blank')
  }

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  const getFileIcon = (name: string) => {
    const ext = name.split('.').pop()?.toLowerCase() || ''
    if (['doc', 'docx'].includes(ext)) return '📄'
    if (['xls', 'xlsx'].includes(ext)) return '📊'
    if (['ppt', 'pptx'].includes(ext)) return '📽️'
    if (['pdf'].includes(ext)) return '📕'
    if (['txt'].includes(ext)) return '📝'
    if (['zip', 'rar', '7z'].includes(ext)) return '📦'
    if (['mp3', 'wav', 'aac'].includes(ext)) return '🎵'
    if (['mp4', 'avi', 'mov'].includes(ext)) return '🎬'
    return '📁'
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>文件管理</h2>
        <p>管理所有上传的非图片文件</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>文件列表</h3>
          <div className="card-actions">
            <button className="btn btn-secondary" onClick={loadFiles}>刷新</button>
          </div>
        </div>
        <div className="card-body">
          {files.length === 0 ? (
            <div className="empty-state">
              <div className="icon">📁</div>
              <p>暂无文件</p>
            </div>
          ) : (
            <table className="data-table">
              <thead>
                <tr>
                  <th>类型</th>
                  <th>文件名</th>
                  <th>大小</th>
                  <th>上传时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {files.map((file) => (
                  <tr key={file.url}>
                    <td style={{ fontSize: 20 }}>{getFileIcon(file.name)}</td>
                    <td style={{ maxWidth: 300, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{file.name}</td>
                    <td>{formatSize(file.size)}</td>
                    <td>{file.created}</td>
                    <td>
                      <button className="btn btn-secondary btn-sm" onClick={() => handleDownload(file.url)}>下载</button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(file.url)} style={{ marginLeft: 8 }}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {total > 0 && (
            <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={pageNum === 1} onClick={() => setPageNum(pageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {pageNum} / {totalPages} 页，共 {total} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={pageNum >= totalPages} onClick={() => setPageNum(pageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <ResponsePane state={state} />
    </div>
  )
}

function GreenFuturePage() {
  const { state, run } = useRunner()
  const [activeTab, setActiveTab] = useState<'card' | 'question' | 'analysis'>('card')
  const [cards, setCards] = useState<unknown[]>([])
  const [questionType, setQuestionType] = useState('1')
  const [questionLevel, setQuestionLevel] = useState('1')
  const [questions, setQuestions] = useState<unknown[]>([])
  const [questionList, setQuestionList] = useState<unknown[]>([])
  const [questionPage, setQuestionPage] = useState(1)
  const [questionTotal, setQuestionTotal] = useState(0)
  const [paperScore, setPaperScore] = useState('40')
  const [paperAnswers, setPaperAnswers] = useState('[\n  {"qid": 101, "answer": "B"}\n]')
  const [dataSeries, setDataSeries] = useState<unknown[]>([])

  const [showCardModal, setShowCardModal] = useState(false)
  const [cardModalType, setCardModalType] = useState<'add' | 'edit'>('add')
  const [editingCard, setEditingCard] = useState<unknown | null>(null)
  const [cardForm, setCardForm] = useState({ icon: '', title: '', num: '', unit: '', trend: '', sort: 0 })

  const [showQuestionModal, setShowQuestionModal] = useState(false)
  const [questionModalType, setQuestionModalType] = useState<'add' | 'edit'>('add')
  const [editingQuestion, setEditingQuestion] = useState<unknown | null>(null)
  const [questionForm, setQuestionForm] = useState({ questionType: '1', level: '1', question: '', optionA: '', optionB: '', optionC: '', optionD: '', optionE: '', optionF: '', answer: '', score: 2, status: '0' })

  const [showSeriesModal, setShowSeriesModal] = useState(false)
  const [seriesModalType, setSeriesModalType] = useState<'add' | 'edit'>('add')
  const [editingSeries, setEditingSeries] = useState<unknown | null>(null)
  const [seriesItems, setSeriesItems] = useState<Array<{ name: string; data: string }>>([{ name: '', data: '' }])

  const loadDataCard = () => {
    run('加载数据卡片', async () => {
      const res = await api.dataCard()
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setCards(sortById(data.data as unknown[]))
      }
      return data
    })
  }

  const loadQuestionAdminList = () => {
    run('加载题库列表', async () => {
      const res = await api.questionAdminList({ pageNum: questionPage, pageSize: 10 })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setQuestionList(sortById(data.data as unknown[]))
        setQuestionTotal(Number(data.total || 0))
      }
      return data
    })
  }

  const loadQuestions = () => {
    run('随机抽题', async () => {
      const res = await api.questionList(questionType, questionLevel)
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setQuestions(data.data as unknown[])
      }
      return data
    })
  }

  const submitPaper = () => {
    run('提交答案', async () => {
      let parsedAnswers: Array<{ qid: number; answer: string }> = []
      try {
        parsedAnswers = JSON.parse(paperAnswers) as Array<{ qid: number; answer: string }>
      } catch {
        throw new Error('答案 JSON 格式错误')
      }
      const score = Number(paperScore)
      const res = await api.savePaper({
        score: Number.isFinite(score) ? score : paperScore,
        answer: parsedAnswers,
      })
      return res.data
    })
  }

  const loadDataSeries = () => {
    run('加载数据系列', async () => {
      const res = await api.dataSeriesList()
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setDataSeries(data.data as unknown[])
      }
      return data
    })
  }

  const handleCardSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const payload = { ...cardForm, sort: Number(cardForm.sort) || 0 }
    if (cardModalType === 'add') {
      run('创建卡片', async () => {
        const res = await api.dataCardCreate(payload)
        if (res.data?.code === 200) {
          setShowCardModal(false)
          loadDataCard()
        }
        return res.data
      })
    } else {
      run('更新卡片', async () => {
        const res = await api.dataCardUpdate((editingCard as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowCardModal(false)
          loadDataCard()
        }
        return res.data
      })
    }
  }

  const handleCardDelete = (id: number) => {
    if (!confirm('确定删除?')) return
    run('删除卡片', async () => {
      const res = await api.dataCardDelete(id)
      loadDataCard()
      return res.data
    })
  }

  const handleQuestionSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const payload = { ...questionForm, score: Number(questionForm.score) || 0 }
    if (questionModalType === 'add') {
      run('创建题目', async () => {
        const res = await api.questionCreate(payload)
        if (res.data?.code === 200) {
          setShowQuestionModal(false)
          loadQuestionAdminList()
        }
        return res.data
      })
    } else {
      run('更新题目', async () => {
        const res = await api.questionUpdate((editingQuestion as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowQuestionModal(false)
          loadQuestionAdminList()
        }
        return res.data
      })
    }
  }

  const handleQuestionDelete = (id: number) => {
    if (!confirm('确定删除?')) return
    run('删除题目', async () => {
      const res = await api.questionDelete(id)
      loadQuestionAdminList()
      return res.data
    })
  }

  const handleSeriesSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const validItems = seriesItems.filter(item => item.name.trim())
    if (validItems.length === 0) {
      alert('请至少添加一条数据')
      return
    }
    const dataArray = validItems.map(item => {
      const nums = item.data.split(',').map(s => parseFloat(s.trim())).filter(n => !isNaN(n))
      return { name: item.name.trim(), data: nums }
    })
    const payload = { data: JSON.stringify(dataArray) }
    if (seriesModalType === 'add') {
      run('创建数据', async () => {
        const res = await api.dataSeriesCreate(payload)
        if (res.data?.code === 200) {
          setShowSeriesModal(false)
          loadDataSeries()
        }
        return res.data
      })
    } else {
      run('更新数据', async () => {
        const res = await api.dataSeriesUpdate((editingSeries as { id: number }).id, payload)
        if (res.data?.code === 200) {
          setShowSeriesModal(false)
          loadDataSeries()
        }
        return res.data
      })
    }
  }

  const handleSeriesDelete = (id: number) => {
    if (!confirm('确定删除?')) return
    run('删除数据', async () => {
      const res = await api.dataSeriesDelete(id)
      loadDataSeries()
      return res.data
    })
  }

  useEffect(() => {
    if (activeTab === 'card') {
      loadDataCard()
      return
    }
    if (activeTab === 'question') {
      loadQuestionAdminList()
      return
    }
    if (activeTab === 'analysis') {
      loadDataSeries()
    }
  }, [activeTab])

  useEffect(() => {
    if (activeTab === 'question') {
      loadQuestionAdminList()
    }
  }, [questionPage])

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>绿动未来</h2>
        <p>数据卡片、随机抽题与数据分析管理</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>管理面板</h3>
          <div className="card-actions">
            <button className={`btn ${activeTab === 'card' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setActiveTab('card')}>数据卡片</button>
            <button className={`btn ${activeTab === 'question' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setActiveTab('question')}>环保答题</button>
            <button className={`btn ${activeTab === 'analysis' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setActiveTab('analysis')}>数据分析</button>
          </div>
        </div>
      </div>

      {activeTab === 'card' && (
        <div className="card">
          <div className="card-header">
            <h3>数据卡片管理</h3>
            <div className="card-actions">
              <button className="btn btn-secondary" onClick={loadDataCard}>刷新</button>
              <button className="btn btn-primary" onClick={() => { setCardModalType('add'); setEditingCard(null); setCardForm({ icon: '', title: '', num: '', unit: '', trend: '', sort: 0 }); setShowCardModal(true); }}>新增</button>
            </div>
          </div>
          <div className="card-body">
            <div className="green-card-grid">
              {(cards as Array<{ id: number; title: string; icon: string; num: string; unit: string; trend: string; sort: number }>).map((item) => (
                <div key={item.id} className="green-card-item">
                  <p className="green-card-title">{item.title}</p>
                  <h4 className="green-card-number">{item.num}</h4>
                  <p className="green-card-unit">{item.unit}</p>
                  <p className="green-card-trend">趋势: {item.trend || '-'}</p>
                  <p className="green-card-icon">图标: {item.icon || '-'}</p>
                  <p className="green-card-sort">排序: {item.sort}</p>
                  <div style={{ marginTop: 8, display: 'flex', gap: 4 }}>
                    <button className="btn btn-secondary btn-sm" onClick={() => { setCardModalType('edit'); setEditingCard(item); setCardForm({ icon: String(item.icon || ''), title: String(item.title || ''), num: String(item.num || ''), unit: String(item.unit || ''), trend: String(item.trend || ''), sort: Number(item.sort) || 0 }); setShowCardModal(true); }}>编辑</button>
                    <button className="btn btn-danger btn-sm" onClick={() => handleCardDelete(item.id)}>删除</button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'question' && (
        <>
          <div className="card">
            <div className="card-header">
              <h3>题库管理</h3>
              <div className="card-actions">
                <button className="btn btn-secondary" onClick={loadQuestionAdminList}>刷新</button>
                <button className="btn btn-primary" onClick={() => { setQuestionModalType('add'); setEditingQuestion(null); setQuestionForm({ questionType: '1', level: '1', question: '', optionA: '', optionB: '', optionC: '', optionD: '', optionE: '', optionF: '', answer: '', score: 2, status: '0' }); setShowQuestionModal(true); }}>新增题目</button>
              </div>
            </div>
            <div className="card-body">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>题型</th>
                    <th>难度</th>
                    <th>题目</th>
                    <th>答案</th>
                    <th>分值</th>
                    <th>状态</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {(questionList as Array<{
                    id: number
                    questionType: string
                    level: string
                    question: string
                    optionA: string
                    optionB: string
                    optionC: string
                    optionD: string
                    answer: string
                    score: number
                    status: string
                  }>).map((item) => (
                    <tr key={item.id}>
                      <td>{item.id}</td>
                      <td>{item.questionType === '1' ? '选择题' : item.questionType === '4' ? '判断题' : item.questionType}</td>
                      <td>{item.level === '1' ? '简单' : item.level === '2' ? '中等' : item.level === '3' ? '困难' : item.level}</td>
                      <td style={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{item.question}</td>
                      <td>{item.answer}</td>
                      <td>{item.score}</td>
                      <td>{item.status === '0' ? '启用' : '禁用'}</td>
                      <td>
                        <button className="btn btn-secondary btn-sm" onClick={() => { setQuestionModalType('edit'); setEditingQuestion(item); setQuestionForm({ questionType: String(item.questionType || '1'), level: String(item.level || '1'), question: String(item.question || ''), optionA: String(item.optionA || ''), optionB: String(item.optionB || ''), optionC: String(item.optionC || ''), optionD: String(item.optionD || ''), optionE: '', optionF: '', answer: String(item.answer || ''), score: Number(item.score) || 2, status: String(item.status || '0') }); setShowQuestionModal(true); }}>编辑</button>
                        <button className="btn btn-danger btn-sm" onClick={() => handleQuestionDelete(item.id)}>删除</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {questionTotal > 0 && (
                <div className="pagination" style={{ marginTop: 16, display: 'flex', gap: 8, alignItems: 'center' }}>
                  <button className="btn btn-secondary btn-sm" disabled={questionPage === 1} onClick={() => setQuestionPage(questionPage - 1)}>上一页</button>
                  <span style={{ color: '#6b7280' }}>第 {questionPage} / {Math.ceil(questionTotal / 10)} 页，共 {questionTotal} 条</span>
                  <button className="btn btn-secondary btn-sm" disabled={questionPage >= Math.ceil(questionTotal / 10)} onClick={() => setQuestionPage(questionPage + 1)}>下一页</button>
                </div>
              )}
            </div>
          </div>

          <div className="card">
            <div className="card-header">
              <h3>随机抽题测试</h3>
              <div className="card-actions">
                <button className="btn btn-primary" onClick={loadQuestions}>抽题</button>
              </div>
            </div>
            <div className="card-body">
              <div className="form-row">
                <div className="form-group">
                  <label>题型</label>
                  <select value={questionType} onChange={(e) => setQuestionType(e.target.value)}>
                    <option value="1">选择题</option>
                    <option value="4">判断题</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>难度</label>
                  <select value={questionLevel} onChange={(e) => setQuestionLevel(e.target.value)}>
                    <option value="1">简单</option>
                    <option value="2">中等</option>
                    <option value="3">困难</option>
                  </select>
                </div>
              </div>
              <table className="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>题目</th>
                    <th>选项</th>
                    <th>答案</th>
                    <th>分值</th>
                  </tr>
                </thead>
                <tbody>
                  {(questions as Array<{
                    id: number
                    question: string
                    optionA: string
                    optionB: string
                    optionC: string
                    optionD: string
                    optionE: string
                    optionF: string
                    answer: string
                    score: number
                  }>).map((item) => (
                    <tr key={item.id}>
                      <td>{item.id}</td>
                      <td>{item.question}</td>
                      <td>
                        {[item.optionA, item.optionB, item.optionC, item.optionD, item.optionE, item.optionF]
                          .filter((option) => option)
                          .join(' / ') || '-'}
                      </td>
                      <td>{item.answer}</td>
                      <td>{item.score}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div className="card">
            <div className="card-header">
              <h3>提交答案</h3>
              <div className="card-actions">
                <button className="btn btn-primary" onClick={submitPaper}>提交</button>
              </div>
            </div>
            <div className="card-body">
              <div className="form-group">
                <label>分数</label>
                <input value={paperScore} onChange={(e) => setPaperScore(e.target.value)} placeholder="分数" />
              </div>
              <div className="form-group">
                <label>答案 JSON</label>
                <textarea value={paperAnswers} onChange={(e) => setPaperAnswers(e.target.value)} rows={5} placeholder='[{"qid": 101, "answer": "B"}]' />
              </div>
            </div>
          </div>
        </>
      )}

      {activeTab === 'analysis' && (
        <>
          <div className="card">
            <div className="card-header">
              <h3>数据分析管理</h3>
              <div className="card-actions">
                <button className="btn btn-secondary" onClick={loadDataSeries}>刷新</button>
                <button className="btn btn-primary" onClick={() => { setSeriesModalType('add'); setEditingSeries(null); setSeriesItems([{ name: '', data: '' }]); setShowSeriesModal(true); }}>新增数据</button>
              </div>
            </div>
            <div className="card-body">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>列表Key</th>
                    <th>名称</th>
                    <th>数据</th>
                    <th>排序</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {(dataSeries as Array<{ id: number; listKey: string; name: string; data: string; sort: number }>).map((item) => (
                    <tr key={item.id}>
                      <td>{item.id}</td>
                      <td>{item.listKey}</td>
                      <td>{item.name}</td>
                      <td style={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{item.data}</td>
                      <td>{item.sort}</td>
                      <td>
                        <button className="btn btn-secondary btn-sm" onClick={() => { 
                          setSeriesModalType('edit'); 
                          setEditingSeries(item); 
                          try {
                            const parsed = JSON.parse(String(item.data || '[]'))
                            const items = parsed.map((p: { name: string; data: number[] }) => ({ 
                              name: p.name, 
                              data: Array.isArray(p.data) ? p.data.join(', ') : '' 
                            }))
                            setSeriesItems(items.length ? items : [{ name: '', data: '' }])
                          } catch {
                            setSeriesItems([{ name: '', data: '' }])
                          }
                          setShowSeriesModal(true); 
                        }}>编辑</button>
                        <button className="btn btn-danger btn-sm" onClick={() => handleSeriesDelete(item.id)}>删除</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          </>
      )}

      <ResponsePane state={state} />

      {showCardModal && (
        <div className="modal-overlay" onClick={() => setShowCardModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{cardModalType === 'add' ? '新增卡片' : '编辑卡片'}</h3>
              <button className="modal-close" onClick={() => setShowCardModal(false)}>×</button>
            </div>
            <form onSubmit={handleCardSubmit}>
              <div className="modal-body">
                <ImageUrlPickerField
                  label="图标URL"
                  value={cardForm.icon}
                  onValueChange={(v) => setCardForm({ ...cardForm, icon: v })}
                  placeholder="/static/image/icon/news_hot.png"
                />
                <div className="form-group">
                  <label>标题</label>
                  <input value={cardForm.title} onChange={(e) => setCardForm({ ...cardForm, title: e.target.value })} placeholder="AQI 指数" required />
                </div>
                <div className="form-group">
                  <label>数值</label>
                  <input value={cardForm.num} onChange={(e) => setCardForm({ ...cardForm, num: e.target.value })} placeholder="45" required />
                </div>
                <div className="form-group">
                  <label>单位</label>
                  <input value={cardForm.unit} onChange={(e) => setCardForm({ ...cardForm, unit: e.target.value })} placeholder="优" required />
                </div>
                <ImageUrlPickerField
                  label="趋势图标"
                  value={cardForm.trend}
                  onValueChange={(v) => setCardForm({ ...cardForm, trend: v })}
                  placeholder="/static/image/icon/down_arrow.png"
                />
                <div className="form-group">
                  <label>排序</label>
                  <input type="number" value={cardForm.sort} onChange={(e) => setCardForm({ ...cardForm, sort: Number(e.target.value) || 0 })} placeholder="0" />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowCardModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showQuestionModal && (
        <div className="modal-overlay" onClick={() => setShowQuestionModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{questionModalType === 'add' ? '新增题目' : '编辑题目'}</h3>
              <button className="modal-close" onClick={() => setShowQuestionModal(false)}>×</button>
            </div>
            <form onSubmit={handleQuestionSubmit}>
              <div className="modal-body">
                <div className="form-group">
                  <label>题型</label>
                  <select value={questionForm.questionType} onChange={(e) => setQuestionForm({ ...questionForm, questionType: e.target.value })}>
                    <option value="1">选择题</option>
                    <option value="4">判断题</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>难度</label>
                  <select value={questionForm.level} onChange={(e) => setQuestionForm({ ...questionForm, level: e.target.value })}>
                    <option value="1">简单</option>
                    <option value="2">中等</option>
                    <option value="3">困难</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>题目</label>
                  <textarea value={questionForm.question} onChange={(e) => setQuestionForm({ ...questionForm, question: e.target.value })} rows={3} required />
                </div>
                {questionForm.questionType === '1' && (
                  <>
                    <div className="form-group">
                      <label>选项A</label>
                      <input value={questionForm.optionA} onChange={(e) => setQuestionForm({ ...questionForm, optionA: e.target.value })} placeholder="A选项内容" />
                    </div>
                    <div className="form-group">
                      <label>选项B</label>
                      <input value={questionForm.optionB} onChange={(e) => setQuestionForm({ ...questionForm, optionB: e.target.value })} placeholder="B选项内容" />
                    </div>
                    <div className="form-group">
                      <label>选项C</label>
                      <input value={questionForm.optionC} onChange={(e) => setQuestionForm({ ...questionForm, optionC: e.target.value })} placeholder="C选项内容" />
                    </div>
                    <div className="form-group">
                      <label>选项D</label>
                      <input value={questionForm.optionD} onChange={(e) => setQuestionForm({ ...questionForm, optionD: e.target.value })} placeholder="D选项内容" />
                    </div>
                    <div className="form-group">
                      <label>选项E</label>
                      <input value={questionForm.optionE} onChange={(e) => setQuestionForm({ ...questionForm, optionE: e.target.value })} placeholder="E选项内容" />
                    </div>
                    <div className="form-group">
                      <label>选项F</label>
                      <input value={questionForm.optionF} onChange={(e) => setQuestionForm({ ...questionForm, optionF: e.target.value })} placeholder="F选项内容" />
                    </div>
                  </>
                )}
                <div className="form-group">
                  <label>答案</label>
                  {questionForm.questionType === '4' ? (
                    <select value={questionForm.answer} onChange={(e) => setQuestionForm({ ...questionForm, answer: e.target.value })} required>
                      <option value="">请选择</option>
                      <option value="1">正确</option>
                      <option value="0">错误</option>
                    </select>
                  ) : (
                    <input value={questionForm.answer} onChange={(e) => setQuestionForm({ ...questionForm, answer: e.target.value })} placeholder="A/B/C/D/E/F" required />
                  )}
                </div>
                <div className="form-group">
                  <label>分值</label>
                  <input type="number" value={questionForm.score} onChange={(e) => setQuestionForm({ ...questionForm, score: Number(e.target.value) || 0 })} />
                </div>
                <div className="form-group">
                  <label>状态</label>
                  <select value={questionForm.status} onChange={(e) => setQuestionForm({ ...questionForm, status: e.target.value })}>
                    <option value="0">启用</option>
                    <option value="1">禁用</option>
                  </select>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowQuestionModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showSeriesModal && (
        <div className="modal-overlay" onClick={() => setShowSeriesModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{seriesModalType === 'add' ? '新增数据' : '编辑数据'}</h3>
              <button className="modal-close" onClick={() => setShowSeriesModal(false)}>×</button>
            </div>
            <form onSubmit={handleSeriesSubmit}>
              <div className="modal-body">
                {seriesModalType === 'edit' && (
                  <div className="form-group">
                    <label>列表 Key</label>
                    <input value={editingSeries ? String((editingSeries as { listKey?: string }).listKey || '') : ''} disabled style={{ background: '#f3f4f6' }} />
                  </div>
                )}
                <div className="form-group">
                  <label>数据项</label>
                  {seriesItems.map((item, idx) => (
                    <div key={idx} style={{ display: 'flex', gap: 8, marginBottom: 8, alignItems: 'flex-start' }}>
                      <input
                        value={item.name}
                        onChange={(e) => {
                          const newItems = [...seriesItems]
                          newItems[idx].name = e.target.value
                          setSeriesItems(newItems)
                        }}
                        placeholder="名称如: aqi"
                        style={{ width: 120 }}
                      />
                      <input
                        value={item.data}
                        onChange={(e) => {
                          const newItems = [...seriesItems]
                          newItems[idx].data = e.target.value
                          setSeriesItems(newItems)
                        }}
                        placeholder="数据如: 25, 45, 23, 456"
                        style={{ flex: 1 }}
                      />
                      <button
                        type="button"
                        className="btn btn-danger btn-sm"
                        onClick={() => {
                          const newItems = seriesItems.filter((_, i) => i !== idx)
                          setSeriesItems(newItems.length ? newItems : [{ name: '', data: '' }])
                        }}
                      >
                        ×
                      </button>
                    </div>
                  ))}
                  <button
                    type="button"
                    className="btn btn-secondary btn-sm"
                    onClick={() => setSeriesItems([...seriesItems, { name: '', data: '' }])}
                  >
                    + 添加数据项
                  </button>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowSeriesModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">提交</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

function BizOpsPage() {
  const { state, run } = useRunner()

  const [newsId, setNewsId] = useState('')
  const [newsDetail, setNewsDetail] = useState<Record<string, unknown> | null>(null)
  const [newsComments, setNewsComments] = useState<Array<Record<string, unknown>>>([])
  const [newsCommentUser, setNewsCommentUser] = useState('管理员')
  const [newsCommentText, setNewsCommentText] = useState('')
  const [newsCategoryId, setNewsCategoryId] = useState('')
  const [categoryNews, setCategoryNews] = useState<Array<Record<string, unknown>>>([])

  const [noticeId, setNoticeId] = useState('')
  const [noticeDetail, setNoticeDetail] = useState<Record<string, unknown> | null>(null)

  const [neighborId, setNeighborId] = useState('')
  const [neighborDetail, setNeighborDetail] = useState<Record<string, unknown> | null>(null)

  const [activityId, setActivityId] = useState('')
  const [activityDetail, setActivityDetail] = useState<Record<string, unknown> | null>(null)
  const [activityTopList, setActivityTopList] = useState<Array<Record<string, unknown>>>([])
  const [activitySearchWords, setActivitySearchWords] = useState('')
  const [activitySearchList, setActivitySearchList] = useState<Array<Record<string, unknown>>>([])
  const [activityCategoryId, setActivityCategoryId] = useState('')
  const [activityCategoryList, setActivityCategoryList] = useState<Array<Record<string, unknown>>>([])
  const [activityEvaluate, setActivityEvaluate] = useState('')
  const [activityStar, setActivityStar] = useState(5)

  const [regPageNum, setRegPageNum] = useState(1)
  const [regTotal, setRegTotal] = useState(0)
  const [regFilterActivityId, setRegFilterActivityId] = useState('')
  const [regFilterUserId, setRegFilterUserId] = useState('')
  const [registrations, setRegistrations] = useState<Array<Record<string, unknown>>>([])

  const loadNewsDetail = () => {
    if (!newsId.trim()) return
    run('获取新闻详情', async () => {
      const res = await api.pressNewsDetail(newsId.trim())
      const data = res.data
      if (data?.code === 200) setNewsDetail((data.data || null) as Record<string, unknown> | null)
      return data
    })
  }

  const loadNewsComments = () => {
    if (!newsId.trim()) return
    run('获取新闻评论列表', async () => {
      const res = await api.commentList(newsId.trim(), { pageNum: 1, pageSize: 50 })
      const data = res.data
      if (data?.code === 200 && data?.data) setNewsComments(data.data as Array<Record<string, unknown>>)
      return data
    })
  }

  const likeNews = () => {
    if (!newsId.trim()) return
    run('新闻点赞', async () => (await api.pressLike(newsId.trim())).data)
  }

  const publishNewsComment = () => {
    if (!newsId.trim() || !newsCommentText.trim() || !newsCommentUser.trim()) return
    run('发布新闻评论', async () => {
      const res = await api.pressComment({
        content: newsCommentText.trim(),
        newsId: newsId.trim(),
        userName: newsCommentUser.trim(),
      })
      if (res.data?.code === 200) {
        setNewsCommentText('')
        loadNewsComments()
      }
      return res.data
    })
  }

  const likeNewsComment = (id: number) => {
    run('评论点赞', async () => {
      const res = await api.commentLike(String(id))
      if (res.data?.code === 200) loadNewsComments()
      return res.data
    })
  }

  const loadCategoryNews = () => {
    if (!newsCategoryId.trim()) return
    run('分类新闻列表', async () => {
      const res = await api.pressCategoryNewsList({ pageNum: 1, pageSize: 20, id: newsCategoryId.trim() })
      const data = res.data
      if (data?.code === 200 && data?.data) setCategoryNews(data.data as Array<Record<string, unknown>>)
      return data
    })
  }

  const loadNoticeDetail = () => {
    if (!noticeId.trim()) return
    run('获取公告详情', async () => {
      const res = await api.noticeDetail(noticeId.trim())
      const data = res.data
      if (data?.code === 200) setNoticeDetail((data.data || null) as Record<string, unknown> | null)
      return data
    })
  }

  const markNoticeRead = () => {
    if (!noticeId.trim()) return
    run('公告标记已读', async () => (await api.readNotice(noticeId.trim())).data)
  }

  const loadNeighborDetail = () => {
    if (!neighborId.trim()) return
    run('获取友邻详情', async () => {
      const res = await api.neighborDetail(neighborId.trim())
      const data = res.data
      if (data?.code === 200) setNeighborDetail((data.data || null) as Record<string, unknown> | null)
      return data
    })
  }

  const loadActivityTop = () => {
    run('获取热门活动', async () => {
      const res = await api.activityTopList({ pageNum: 1, pageSize: 20 })
      const data = res.data
      if (data?.code === 200 && data?.data) setActivityTopList(data.data as Array<Record<string, unknown>>)
      return data
    })
  }

  const searchActivities = () => {
    if (!activitySearchWords.trim()) return
    run('活动搜索', async () => {
      const res = await api.activitySearch({ words: activitySearchWords.trim() }, { pageNum: 1, pageSize: 20 })
      const data = res.data
      if (data?.code === 200 && data?.data) setActivitySearchList(data.data as Array<Record<string, unknown>>)
      return data
    })
  }

  const loadActivityDetail = () => {
    if (!activityId.trim()) return
    run('活动详情', async () => {
      const res = await api.activityDetail(activityId.trim())
      const data = res.data
      if (data?.code === 200) setActivityDetail((data.data || null) as Record<string, unknown> | null)
      return data
    })
  }

  const loadActivityCategoryList = () => {
    if (!activityCategoryId.trim()) return
    run('分类活动列表', async () => {
      const res = await api.activityCategoryList(activityCategoryId.trim(), { pageNum: 1, pageSize: 20 })
      const data = res.data
      if (data?.code === 200 && data?.data) setActivityCategoryList(data.data as Array<Record<string, unknown>>)
      return data
    })
  }

  const doRegistration = () => {
    if (!activityId.trim()) return
    run('活动报名', async () => (await api.registration({ activityId: Number(activityId) })).data)
  }

  const doCheckin = () => {
    if (!activityId.trim()) return
    run('活动签到', async () => (await api.checkin(activityId.trim())).data)
  }

  const doRegistrationComment = () => {
    if (!activityId.trim() || !activityEvaluate.trim()) return
    run('活动评价', async () => (await api.registrationComment(activityId.trim(), { evaluate: activityEvaluate.trim(), star: activityStar })).data)
  }

  const loadRegistrationList = () => {
    run('报名记录列表', async () => {
      const res = await api.registrationList({
        pageNum: regPageNum,
        pageSize: 10,
        activityId: regFilterActivityId.trim() || undefined,
        userId: regFilterUserId.trim() || undefined,
      })
      const data = res.data
      if (data?.code === 200 && data?.data) {
        setRegistrations(sortById(data.data as Array<Record<string, unknown>>))
        setRegTotal(Number(data.total || 0))
      }
      return data
    })
  }

  useEffect(() => {
    loadRegistrationList()
  }, [regPageNum])

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>业务操作</h2>
        <p>非CRUD接口的可视化管理模块</p>
      </div>

      <div className="card">
        <div className="card-header"><h3>新闻互动管理</h3></div>
        <div className="card-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={newsId} onChange={(e) => setNewsId(e.target.value)} placeholder="新闻ID" />
            <button className="btn btn-secondary" onClick={loadNewsDetail}>详情</button>
            <button className="btn btn-secondary" onClick={likeNews}>点赞</button>
            <button className="btn btn-secondary" onClick={loadNewsComments}>评论列表</button>
          </div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={newsCommentUser} onChange={(e) => setNewsCommentUser(e.target.value)} placeholder="评论用户" />
            <input value={newsCommentText} onChange={(e) => setNewsCommentText(e.target.value)} placeholder="评论内容" style={{ flex: 1 }} />
            <button className="btn btn-primary" onClick={publishNewsComment}>发布评论</button>
          </div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={newsCategoryId} onChange={(e) => setNewsCategoryId(e.target.value)} placeholder="分类ID" />
            <button className="btn btn-secondary" onClick={loadCategoryNews}>分类新闻列表</button>
          </div>
          {newsDetail && <p style={{ fontSize: 13, color: '#374151' }}>新闻: {(newsDetail.title as string) || '-'} | 点赞: {String(newsDetail.likeNum || 0)} | 评论: {String(newsDetail.commentNum || 0)}</p>}
          {newsComments.length > 0 && (
            <table className="data-table" style={{ marginTop: 8 }}>
              <thead><tr><th>评论人</th><th>内容</th><th>点赞</th><th>操作</th></tr></thead>
              <tbody>
                {newsComments.map((c) => (
                  <tr key={String(c.id)}>
                    <td>{String(c.userName || '-')}</td>
                    <td>{String(c.content || '-')}</td>
                    <td>{String(c.likeNum || 0)}</td>
                    <td><button className="btn btn-secondary btn-sm" onClick={() => likeNewsComment(Number(c.id || 0))}>点赞评论</button></td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {categoryNews.length > 0 && <p style={{ marginTop: 8, color: '#6b7280' }}>分类新闻数量: {categoryNews.length}</p>}
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3>公告阅读管理</h3></div>
        <div className="card-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={noticeId} onChange={(e) => setNoticeId(e.target.value)} placeholder="公告ID" />
            <button className="btn btn-secondary" onClick={loadNoticeDetail}>详情</button>
            <button className="btn btn-primary" onClick={markNoticeRead}>标记已读</button>
          </div>
          {noticeDetail && <p style={{ fontSize: 13, color: '#374151' }}>公告: {(noticeDetail.title as string) || '-'} | 状态: {String(noticeDetail.noticeStatus || '-')}</p>}
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3>友邻详情管理</h3></div>
        <div className="card-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={neighborId} onChange={(e) => setNeighborId(e.target.value)} placeholder="友邻帖子ID" />
            <button className="btn btn-secondary" onClick={loadNeighborDetail}>详情</button>
          </div>
          {neighborDetail && <p style={{ fontSize: 13, color: '#374151' }}>用户: {String(neighborDetail.publishName || '-')} | 评论数: {String(neighborDetail.commentNum || 0)}</p>}
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3>活动业务管理</h3></div>
        <div className="card-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={activityId} onChange={(e) => setActivityId(e.target.value)} placeholder="活动ID" />
            <button className="btn btn-secondary" onClick={loadActivityDetail}>活动详情</button>
            <button className="btn btn-secondary" onClick={doRegistration}>报名</button>
            <button className="btn btn-secondary" onClick={doCheckin}>签到</button>
          </div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={activitySearchWords} onChange={(e) => setActivitySearchWords(e.target.value)} placeholder="搜索关键词" />
            <button className="btn btn-secondary" onClick={searchActivities}>活动搜索</button>
            <input value={activityCategoryId} onChange={(e) => setActivityCategoryId(e.target.value)} placeholder="分类ID" />
            <button className="btn btn-secondary" onClick={loadActivityCategoryList}>分类活动列表</button>
            <button className="btn btn-secondary" onClick={loadActivityTop}>热门活动</button>
          </div>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={activityEvaluate} onChange={(e) => setActivityEvaluate(e.target.value)} placeholder="评价内容" style={{ flex: 1 }} />
            <select value={activityStar} onChange={(e) => setActivityStar(Number(e.target.value))} style={{ width: 90 }}>
              <option value={5}>5 星</option><option value={4}>4 星</option><option value={3}>3 星</option><option value={2}>2 星</option><option value={1}>1 星</option>
            </select>
            <button className="btn btn-primary" onClick={doRegistrationComment}>提交评价</button>
          </div>
          {activityDetail && <p style={{ fontSize: 13, color: '#374151' }}>活动: {String(activityDetail.title || '-')} | 地址: {String(activityDetail.address || '-')}</p>}
          <p style={{ margin: '8px 0 0', color: '#6b7280' }}>热门: {activityTopList.length} 条，搜索: {activitySearchList.length} 条，分类: {activityCategoryList.length} 条</p>
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3>报名记录管理</h3></div>
        <div className="card-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <input value={regFilterActivityId} onChange={(e) => setRegFilterActivityId(e.target.value)} placeholder="筛选活动ID" />
            <input value={regFilterUserId} onChange={(e) => setRegFilterUserId(e.target.value)} placeholder="筛选用户ID" />
            <button className="btn btn-secondary" onClick={() => { setRegPageNum(1); loadRegistrationList() }}>筛选</button>
          </div>
          <table className="data-table">
            <thead><tr><th>ID</th><th>活动ID</th><th>用户</th><th>签到</th><th>评分</th><th>评价</th></tr></thead>
            <tbody>
              {registrations.map((r) => (
                <tr key={String(r.id)}>
                  <td>{String(r.id)}</td>
                  <td>{String(r.activityId)}</td>
                  <td>{String(r.nickName || r.userName || '-')}</td>
                  <td>{String(r.checkinStatus || '0')}</td>
                  <td>{String(r.star || 0)}</td>
                  <td style={{ maxWidth: 260, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{String(r.comment || '-')}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {regTotal > 0 && (
            <div className="pagination" style={{ marginTop: 12, display: 'flex', gap: 8, alignItems: 'center' }}>
              <button className="btn btn-secondary btn-sm" disabled={regPageNum === 1} onClick={() => setRegPageNum(regPageNum - 1)}>上一页</button>
              <span style={{ color: '#6b7280' }}>第 {regPageNum} / {Math.max(1, Math.ceil(regTotal / 10))} 页，共 {regTotal} 条</span>
              <button className="btn btn-secondary btn-sm" disabled={regPageNum >= Math.ceil(regTotal / 10)} onClick={() => setRegPageNum(regPageNum + 1)}>下一页</button>
            </div>
          )}
        </div>
      </div>

      <ResponsePane state={state} />
    </div>
  )
}

function PlaygroundPage() {
  const { state, run } = useRunner()
  const [selectedKey, setSelectedKey] = useState(endpointCatalog[0].key)
  const [pathParamsText, setPathParamsText] = useState('{}')
  const [queryText, setQueryText] = useState('{}')
  const [bodyText, setBodyText] = useState('{}')
  const [uploadFile, setUploadFile] = useState<File | null>(null)

  const endpoint = useMemo(
    () => endpointCatalog.find((item) => item.key === selectedKey) || endpointCatalog[0],
    [selectedKey],
  )

  const invokeEndpoint = () => {
    run(`测试台: ${endpoint.description}`, async () => {
      if (endpoint.key === 'upload') {
        if (!uploadFile) {
          throw new Error('请先选择文件')
        }
        return (await api.upload(uploadFile)).data
      }

      const pathParams = JSON.parse(pathParamsText || '{}') as Record<string, string>
      const query = JSON.parse(queryText || '{}') as Record<string, unknown>
      const body = JSON.parse(bodyText || '{}') as Record<string, unknown>

      const url = endpoint.path.replace(/\{(\w+)\}/g, (_: string, key: string) => {
        const value = pathParams[key]
        if (!value) throw new Error(`缺少路径参数: ${key}`)
        return encodeURIComponent(value)
      })

      const response = await apiClient.request({
        method: endpoint.method,
        url,
        params: endpoint.method === 'GET' ? query : query,
        data: endpoint.method === 'GET' ? undefined : body,
      })
      return response.data
    })
  }

  return (
    <div className="page-grid">
      <div className="page-header">
        <h2>API 测试台</h2>
        <p>测试所有接口</p>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>接口选择</h3>
          <span className={`status-badge ${endpoint.auth ? 'status-loading' : 'status-idle'}`}>
            {endpoint.auth ? '需要token' : '无需token'}
          </span>
        </div>
        <div className="card-body">
          <div className="form-group">
            <label>选择接口</label>
            <select value={selectedKey} onChange={(e) => setSelectedKey(e.target.value)}>
              {endpointCatalog.map((item) => (
                <option key={item.key} value={item.key}>
                  [{item.method}] {item.path} - {item.description}
                </option>
              ))}
            </select>
          </div>
          <p style={{ color: '#6b7280', fontSize: 14 }}>
            当前接口: <code style={{ background: '#f3f4f6', padding: '2px 6px', borderRadius: 4 }}>{endpoint.path}</code>
          </p>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3>请求参数</h3>
        </div>
        <div className="card-body">
          {endpoint.key === 'upload' ? (
            <div className="form-group">
              <label>选择文件</label>
              <input
                type="file"
                onChange={(e) => setUploadFile(e.target.files && e.target.files[0] ? e.target.files[0] : null)}
              />
            </div>
          ) : (
            <>
              <div className="form-group">
                <label>路径参数 (Path Parameters)</label>
                <textarea
                  value={pathParamsText}
                  onChange={(e) => setPathParamsText(e.target.value)}
                  rows={3}
                  placeholder='{"id": "123"}'
                />
              </div>
              <div className="form-group">
                <label>查询参数 (Query Parameters)</label>
                <textarea
                  value={queryText}
                  onChange={(e) => setQueryText(e.target.value)}
                  rows={3}
                  placeholder='{"pageNum": 1, "pageSize": 10}'
                />
              </div>
              {endpoint.method !== 'GET' && (
                <div className="form-group">
                  <label>请求体 (Request Body)</label>
                  <textarea
                    value={bodyText}
                    onChange={(e) => setBodyText(e.target.value)}
                    rows={4}
                    placeholder='{"key": "value"}'
                  />
                </div>
              )}
            </>
          )}
          <button className="btn btn-primary" onClick={invokeEndpoint}>
            发送请求
          </button>
        </div>
      </div>

      <ResponsePane state={state} />
    </div>
  )
}

const navItems = [
  { to: '/dashboard', label: '仪表盘', icon: '📊' },
  { to: '/user', label: '用户管理', icon: '👤' },
  { to: '/news', label: '新闻管理', icon: '📰' },
  { to: '/notice', label: '公告管理', icon: '📢' },
  { to: '/neighbor', label: '友邻帖子', icon: '💬' },
  { to: '/activity', label: '社区活动', icon: '🎉' },
  { to: '/upload', label: '文件上传', icon: '📁' },
  { to: '/green', label: '绿动未来', icon: '🌿' },
  { to: '/images', label: '图片管理', icon: '🖼️' },
  { to: '/files', label: '文件管理', icon: '📂' },
  { to: '/biz', label: '业务操作', icon: '🧩' },
  { to: '/playground', label: 'API测试台', icon: '🧪' },
]

function AdminLayout() {
  const navigate = useNavigate()
  const token = useAuthStore((s) => s.token)
  const clearToken = useAuthStore((s) => s.clearToken)

  if (!token) {
    return <Navigate to="/login" replace />
  }

  const handleLogout = async () => {
    try {
      await api.logout()
    } catch {
      // ignore logout API error
    }
    clearToken()
    navigate('/login')
  }

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <div className="brand">
          <p>Digital Community</p>
          <h1>Admin Panel</h1>
        </div>
        <nav>
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            >
              <span className="nav-icon">{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <main className="content">
        <header className="topbar">
          <div className="topbar-left">
            <h2>数字社区管理后台</h2>
            <p>当前已登录（token长度 {token.length}）</p>
          </div>
          <button className="logout-btn" onClick={handleLogout}>
            退出登录
          </button>
        </header>

        <Routes>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/user" element={<UserPage />} />
          <Route path="/news" element={<NewsPage />} />
          <Route path="/notice" element={<NoticePage />} />
          <Route path="/neighbor" element={<NeighborPage />} />
          <Route path="/activity" element={<ActivityPage />} />
          <Route path="/upload" element={<UploadPage />} />
          <Route path="/green" element={<GreenFuturePage />} />
          <Route path="/images" element={<ImagePage />} />
          <Route path="/files" element={<FilePage />} />
          <Route path="/biz" element={<BizOpsPage />} />
          <Route path="/playground" element={<PlaygroundPage />} />
        </Routes>
      </main>
    </div>
  )
}

function App() {
  const token = useAuthStore((s) => s.token)

  return (
    <Routes>
      <Route
        path="/login"
        element={token ? <Navigate to="/dashboard" replace /> : <LoginPage />}
      />
      <Route path="/*" element={<AdminLayout />} />
    </Routes>
  )
}

export default App

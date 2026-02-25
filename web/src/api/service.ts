import { apiClient, type ApiResponse } from './client'

export const api = {
  logout: () => apiClient.post<ApiResponse>('/logout'),

  phoneLogin: (payload: { phone: string; SMSCode: string }) =>
    apiClient.post<ApiResponse>('/prod-api/api/phone/login', payload),
  login: (payload: { userName: string; passWord: string }) =>
    apiClient.post<ApiResponse>('/prod-api/api/login', payload),
  smsCode: (phone: string) =>
    apiClient.get<ApiResponse<string>>('/prod-api/api/smsCode', { params: { phone } }),
  register: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/register', payload),

  dataCard: () => apiClient.get<ApiResponse>('/prod-api/api/common/datacard'),
  dataCardCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/common/datacard', payload),
  dataCardUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/common/datacard/${id}`, payload),
  dataCardDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/common/datacard/${id}`),
  questionList: (id: string, level: string) =>
    apiClient.get<ApiResponse>(`/prod-api/api/question/questionList/${id}/${level}`),
  questionAdminList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/question/list', { params }),
  questionCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/question', payload),
  questionUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/question/${id}`, payload),
  questionDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/question/${id}`),
  savePaper: (payload: { score: number | string; answer: Array<{ qid: number; answer: string }> }) =>
    apiClient.post<ApiResponse>('/prod-api/api/question/savePaper', payload),
  dataSeriesList: () => apiClient.get<ApiResponse>('/prod-api/api/data/list'),
  dataSeriesByKey: (listKey: string) => apiClient.get<ApiResponse>(`/prod-api/api/data/${listKey}`),
  dataSeriesCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/data/list', payload),
  dataSeriesUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/data/list/${id}`, payload),
  dataSeriesDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/data/list/${id}`),

  getUserInfo: () => apiClient.get<ApiResponse>('/prod-api/api/user/getUserInfo'),
  updateUserInfo: (payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>('/prod-api/api/user/updateUserInfo', payload),
  resetPwd: (payload: { oldPassword: string; newPassword: string }) =>
    apiClient.put<ApiResponse>('/prod-api/api/user/resetPwd', payload),

  userList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/user/list', { params }),
  userCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/user', payload),
  userUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/user/${id}`, payload),
  userDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/user/${id}`),

  rotationList: (params: { pageNum: number; pageSize: number; type: string }) =>
    apiClient.get<ApiResponse>('/prod-api/api/rotation/list', { params }),
  rotationCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/rotation', payload),
  rotationUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/rotation/${id}`, payload),
  rotationDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/rotation/${id}`),

  pressCategoryList: () => apiClient.get<ApiResponse>('/prod-api/api/press/category/list'),
  pressCategoryCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/press/category', payload),
  pressCategoryUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/press/category/${id}`, payload),
  pressCategoryDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/press/category/${id}`),

  pressNewsList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/press/newsList', { params }),
  pressNewsCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/press/news', payload),
  pressNewsUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/press/news/${id}`, payload),
  pressNewsDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/press/news/${id}`),
  pressCategoryNewsList: (params: { pageNum: number; pageSize: number; id: string }) =>
    apiClient.get<ApiResponse>('/prod-api/api/press/category/newsList', { params }),
  pressNewsDetail: (id: string) => apiClient.get<ApiResponse>(`/prod-api/api/press/news/${id}`),
  pressLike: (id: string) => apiClient.put<ApiResponse>(`/prod-api/api/press/like/${id}`),

  pressComment: (payload: { content: string; newsId: string; userName: string }) =>
    apiClient.post<ApiResponse>('/prod-api/api/comment/pressComment', payload),
  commentList: (id: string, params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>(`/prod-api/api/comment/comment/${id}`, { params }),
  commentLike: (id: string) => apiClient.put<ApiResponse>(`/prod-api/api/comment/like/${id}`),

  upload: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return apiClient.post<ApiResponse>('/prod-api/api/common/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  imageList: (params?: { pageNum?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/common/images', { params }),
  imageDelete: (url: string) => apiClient.delete<ApiResponse>('/prod-api/api/common/images', { params: { url } }),

  fileList: (params?: { pageNum?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/common/files', { params }),
  fileDelete: (url: string) => apiClient.delete<ApiResponse>('/prod-api/api/common/files', { params: { url } }),

  noticeList: (params: { pageNum: number; pageSize: number; noticeStatus: string }) =>
    apiClient.get<ApiResponse>('/prod-api/api/notice/list', { params }),
  noticeCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/notice', payload),
  noticeUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/notice/${id}`, payload),
  noticeDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/notice/${id}`),
  noticeDetail: (id: string) => apiClient.get<ApiResponse>(`/prod-api/api/notice/${id}`),
  readNotice: (id: string) => apiClient.put<ApiResponse>(`/prod-api/api/readNotice/${id}`),

  neighborList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/friendly_neighborhood/list', { params }),
  neighborCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/friendly_neighborhood', payload),
  neighborUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/friendly_neighborhood/${id}`, payload),
  neighborDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/friendly_neighborhood/${id}`),
  neighborAddComment: (payload: { content: string; neighborhoodId: number }) =>
    apiClient.post<ApiResponse>('/prod-api/api/friendly_neighborhood/add/comment', payload),
  neighborDetail: (id: string) => apiClient.get<ApiResponse>(`/prod-api/api/friendly_neighborhood/${id}`),

  activityTopList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/activity/topList', { params }),
  activityList: (params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>('/prod-api/api/activity/list', { params }),
  activitySearch: (payload: { words: string }, params: { pageNum: number; pageSize: number }) =>
    apiClient.post<ApiResponse>('/prod-api/api/activity/search', payload, { params }),
  activityDetail: (id: string) => apiClient.get<ApiResponse>(`/prod-api/api/activity/${id}`),
  activityCategoryList: (id: string, params: { pageNum: number; pageSize: number }) =>
    apiClient.get<ApiResponse>(`/prod-api/api/activity/category/list/${id}`, { params }),
  activityCreate: (payload: Record<string, unknown>) =>
    apiClient.post<ApiResponse>('/prod-api/api/activity', payload),
  activityUpdate: (id: number, payload: Record<string, unknown>) =>
    apiClient.put<ApiResponse>(`/prod-api/api/activity/${id}`, payload),
  activityDelete: (id: number) =>
    apiClient.delete<ApiResponse>(`/prod-api/api/activity/${id}`),

  registrationList: (params: { pageNum: number; pageSize: number; activityId?: string; userId?: string }) =>
    apiClient.get<ApiResponse>('/prod-api/api/registration/list', { params }),
  registration: (payload: { activityId: number }) =>
    apiClient.post<ApiResponse>('/prod-api/api/registration', payload),
  checkin: (id: string) => apiClient.put<ApiResponse>(`/prod-api/api/checkin/${id}`),
  registrationComment: (id: string, payload: { evaluate: string; star: number }) =>
    apiClient.put<ApiResponse>(`/prod-api/api/registration/comment/${id}`, payload),
}

export type ApiEndpoint = {
  key: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE'
  path: string
  auth: boolean
  description: string
}

export const endpointCatalog: ApiEndpoint[] = [
  { key: 'logout', method: 'POST', path: '/logout', auth: false, description: '退出登录' },
  { key: 'phoneLogin', method: 'POST', path: '/prod-api/api/phone/login', auth: false, description: '手机登录' },
  { key: 'login', method: 'POST', path: '/prod-api/api/login', auth: false, description: '用户名登录' },
  { key: 'smsCode', method: 'GET', path: '/prod-api/api/smsCode', auth: false, description: '获取短信验证码' },
  { key: 'register', method: 'POST', path: '/prod-api/api/register', auth: false, description: '注册' },
  { key: 'dataCard', method: 'GET', path: '/prod-api/api/common/datacard', auth: false, description: '绿动未来数据卡片' },
  { key: 'dataCardCreate', method: 'POST', path: '/prod-api/api/common/datacard', auth: true, description: '创建数据卡片' },
  { key: 'dataCardUpdate', method: 'PUT', path: '/prod-api/api/common/datacard/{id}', auth: true, description: '更新数据卡片' },
  { key: 'dataCardDelete', method: 'DELETE', path: '/prod-api/api/common/datacard/{id}', auth: true, description: '删除数据卡片' },
  { key: 'questionList', method: 'GET', path: '/prod-api/api/question/questionList/{id}/{level}', auth: true, description: '随机抽题' },
  { key: 'questionAdminList', method: 'GET', path: '/prod-api/api/question/list', auth: true, description: '题库列表' },
  { key: 'questionCreate', method: 'POST', path: '/prod-api/api/question', auth: true, description: '创建题目' },
  { key: 'questionUpdate', method: 'PUT', path: '/prod-api/api/question/{id}', auth: true, description: '更新题目' },
  { key: 'questionDelete', method: 'DELETE', path: '/prod-api/api/question/{id}', auth: true, description: '删除题目' },
  { key: 'savePaper', method: 'POST', path: '/prod-api/api/question/savePaper', auth: true, description: '提交答案' },
  { key: 'dataSeriesList', method: 'GET', path: '/prod-api/api/data/list', auth: true, description: '数据系列列表' },
  { key: 'dataSeriesCreate', method: 'POST', path: '/prod-api/api/data/list', auth: true, description: '创建数据系列' },
  { key: 'dataSeriesUpdate', method: 'PUT', path: '/prod-api/api/data/list/{id}', auth: true, description: '更新数据系列' },
  { key: 'dataSeriesDelete', method: 'DELETE', path: '/prod-api/api/data/list/{id}', auth: true, description: '删除数据系列' },
  { key: 'userList', method: 'GET', path: '/prod-api/api/user/list', auth: true, description: '用户列表' },
  { key: 'userCreate', method: 'POST', path: '/prod-api/api/user', auth: true, description: '创建用户' },
  { key: 'userUpdate', method: 'PUT', path: '/prod-api/api/user/{id}', auth: true, description: '更新用户' },
  { key: 'userDelete', method: 'DELETE', path: '/prod-api/api/user/{id}', auth: true, description: '删除用户' },
  { key: 'getUserInfo', method: 'GET', path: '/prod-api/api/user/getUserInfo', auth: true, description: '获取用户信息' },
  { key: 'updateUserInfo', method: 'PUT', path: '/prod-api/api/user/updateUserInfo', auth: true, description: '更新用户信息' },
  { key: 'resetPwd', method: 'PUT', path: '/prod-api/api/user/resetPwd', auth: true, description: '重置密码' },
  { key: 'rotationList', method: 'GET', path: '/prod-api/api/rotation/list', auth: false, description: '轮播列表' },
  { key: 'rotationCreate', method: 'POST', path: '/prod-api/api/rotation', auth: true, description: '创建轮播图' },
  { key: 'rotationUpdate', method: 'PUT', path: '/prod-api/api/rotation/{id}', auth: true, description: '更新轮播图' },
  { key: 'rotationDelete', method: 'DELETE', path: '/prod-api/api/rotation/{id}', auth: true, description: '删除轮播图' },
  { key: 'pressCategoryList', method: 'GET', path: '/prod-api/api/press/category/list', auth: true, description: '新闻分类列表' },
  { key: 'pressCategoryCreate', method: 'POST', path: '/prod-api/api/press/category', auth: true, description: '创建新闻分类' },
  { key: 'pressCategoryUpdate', method: 'PUT', path: '/prod-api/api/press/category/{id}', auth: true, description: '更新新闻分类' },
  { key: 'pressCategoryDelete', method: 'DELETE', path: '/prod-api/api/press/category/{id}', auth: true, description: '删除新闻分类' },
  { key: 'pressNewsList', method: 'GET', path: '/prod-api/api/press/newsList', auth: true, description: '新闻列表' },
  { key: 'pressNewsCreate', method: 'POST', path: '/prod-api/api/press/news', auth: true, description: '创建新闻' },
  { key: 'pressNewsUpdate', method: 'PUT', path: '/prod-api/api/press/news/{id}', auth: true, description: '更新新闻' },
  { key: 'pressNewsDelete', method: 'DELETE', path: '/prod-api/api/press/news/{id}', auth: true, description: '删除新闻' },
  { key: 'pressCategoryNewsList', method: 'GET', path: '/prod-api/api/press/category/newsList', auth: true, description: '分类新闻列表' },
  { key: 'pressNewsDetail', method: 'GET', path: '/prod-api/api/press/news/{id}', auth: true, description: '新闻详情' },
  { key: 'pressLike', method: 'PUT', path: '/prod-api/api/press/like/{id}', auth: true, description: '新闻点赞' },
  { key: 'pressComment', method: 'POST', path: '/prod-api/api/comment/pressComment', auth: true, description: '新闻评论' },
  { key: 'commentList', method: 'GET', path: '/prod-api/api/comment/comment/{id}', auth: true, description: '新闻评论列表' },
  { key: 'commentLike', method: 'PUT', path: '/prod-api/api/comment/like/{id}', auth: true, description: '评论点赞' },
  { key: 'upload', method: 'POST', path: '/prod-api/api/common/upload', auth: true, description: '通用上传' },
  { key: 'imageList', method: 'GET', path: '/prod-api/api/common/images', auth: false, description: '图片列表' },
  { key: 'imageDelete', method: 'DELETE', path: '/prod-api/api/common/images', auth: true, description: '删除图片' },
  { key: 'fileList', method: 'GET', path: '/prod-api/api/common/files', auth: false, description: '文件列表' },
  { key: 'fileDelete', method: 'DELETE', path: '/prod-api/api/common/files', auth: true, description: '删除文件' },
  { key: 'noticeList', method: 'GET', path: '/prod-api/api/notice/list', auth: false, description: '通知列表' },
  { key: 'noticeCreate', method: 'POST', path: '/prod-api/api/notice', auth: true, description: '创建公告' },
  { key: 'noticeUpdate', method: 'PUT', path: '/prod-api/api/notice/{id}', auth: true, description: '更新公告' },
  { key: 'noticeDelete', method: 'DELETE', path: '/prod-api/api/notice/{id}', auth: true, description: '删除公告' },
  { key: 'noticeDetail', method: 'GET', path: '/prod-api/api/notice/{id}', auth: true, description: '通知详情' },
  { key: 'readNotice', method: 'PUT', path: '/prod-api/api/readNotice/{id}', auth: true, description: '通知已读' },
  { key: 'neighborList', method: 'GET', path: '/prod-api/api/friendly_neighborhood/list', auth: false, description: '友邻圈列表' },
  { key: 'neighborCreate', method: 'POST', path: '/prod-api/api/friendly_neighborhood', auth: true, description: '创建友邻帖子' },
  { key: 'neighborUpdate', method: 'PUT', path: '/prod-api/api/friendly_neighborhood/{id}', auth: true, description: '更新友邻帖子' },
  { key: 'neighborDelete', method: 'DELETE', path: '/prod-api/api/friendly_neighborhood/{id}', auth: true, description: '删除友邻帖子' },
  { key: 'neighborAddComment', method: 'POST', path: '/prod-api/api/friendly_neighborhood/add/comment', auth: false, description: '友邻圈评论' },
  { key: 'neighborDetail', method: 'GET', path: '/prod-api/api/friendly_neighborhood/{id}', auth: false, description: '友邻圈详情' },
  { key: 'activityTopList', method: 'GET', path: '/prod-api/api/activity/topList', auth: false, description: '热门活动' },
  { key: 'activityList', method: 'GET', path: '/prod-api/api/activity/list', auth: false, description: '活动列表' },
  { key: 'activitySearch', method: 'POST', path: '/prod-api/api/activity/search', auth: false, description: '活动搜索' },
  { key: 'activityDetail', method: 'GET', path: '/prod-api/api/activity/{id}', auth: false, description: '活动详情' },
  { key: 'activityCategoryList', method: 'GET', path: '/prod-api/api/activity/category/list/{id}', auth: false, description: '分类活动列表' },
  { key: 'activityCreate', method: 'POST', path: '/prod-api/api/activity', auth: true, description: '创建活动' },
  { key: 'activityUpdate', method: 'PUT', path: '/prod-api/api/activity/{id}', auth: true, description: '更新活动' },
  { key: 'activityDelete', method: 'DELETE', path: '/prod-api/api/activity/{id}', auth: true, description: '删除活动' },
  { key: 'registrationList', method: 'GET', path: '/prod-api/api/registration/list', auth: true, description: '活动报名列表' },
  { key: 'registration', method: 'POST', path: '/prod-api/api/registration', auth: true, description: '活动报名' },
  { key: 'checkin', method: 'PUT', path: '/prod-api/api/checkin/{id}', auth: true, description: '活动签到' },
  { key: 'registrationComment', method: 'PUT', path: '/prod-api/api/registration/comment/{id}', auth: true, description: '活动评论' },
]

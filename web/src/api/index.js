import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 用户相关
export const userApi = {
  register: (data) => api.post('/user/register', data),
  login: (data) => api.post('/user/login', data),
  getUser: (id) => api.get(`/user/${id}`),
  getUsers: (params) => api.get('/users', { params })
}

// 任务相关
export const taskApi = {
  create: (data) => api.post('/task', data),
  getTask: (id) => api.get(`/task/${id}`),
  getUserTasks: (userId, params) => api.get(`/task/user/${userId}`, { params }),
  getTodayTasks: (userId) => api.get(`/task/today/${userId}`),
  update: (id, data) => api.put(`/task/${id}`, data),
  delete: (id) => api.delete(`/task/${id}`)
}

// 打卡相关
export const checkinApi = {
  checkin: (taskId, data) => api.post(`/checkin/${taskId}`, data),
  getUserCheckins: (userId, params) => api.get(`/checkin/user/${userId}`, { params }),
  getTodayChecked: (userId) => api.get(`/checkin/today/${userId}`)
}

// 积分钱包相关
export const walletApi = {
  getWallet: (userId, params) => api.get(`/wallet/${userId}`, { params }),
  getBalance: (userId) => api.get(`/wallet/${userId}/balance`),
  spend: (data) => api.post('/wallet/spend', data),
  addRecord: (data) => api.post('/wallet/add', data),
  delete: (id, userId) => api.delete(`/wallet/delete/${id}`, { data: { user_id: userId } })
}

// 积分历史
export const pointsApi = {
  getHistory: (userId, params) => api.get(`/points/${userId}`, { params })
}

export default api
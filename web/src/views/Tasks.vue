<template>
  <div class="tasks-container">
    <header>
      <h1>每日打卡</h1>
      <div class="user-info">
        <span>{{ user?.username }}</span>
        <span class="points">积分: {{ user?.points || 0 }}</span>
        <router-link to="/points" class="wallet-btn">钱包</router-link>
      </div>
    </header>

    <div class="create-task">
      <h3>创建新任务</h3>
      <form @submit.prevent="createTask">
        <input v-model="newTask.title" placeholder="任务标题" required />
        <select v-model="newTask.circle_mode">
          <option value="once">单次</option>
          <option value="weekly">每周</option>
          <option value="workday">工作日</option>
          <option value="weekend">周末</option>
          <option value="custom">自定义</option>
        </select>
        <select v-model="newTask.level">
          <option value="1">低风险</option>
          <option value="2">中风险</option>
          <option value="3">高风险</option>
        </select>
        <input v-model.number="newTask.points" type="number" placeholder="积分" min="1" />
        <button type="submit" :disabled="creating">{{ creating ? '创建中...' : '创建' }}</button>
      </form>
    </div>

    <div class="task-list">
      <h3>我的任务（按风险等级排序）</h3>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="tasks.length === 0" class="empty">暂无任务</div>
      <div v-else class="tasks">
        <div class="task-item" v-for="task in tasks" :key="task.id" :class="'level-' + task.level">
          <div class="task-info">
            <h4>
              <span class="level-badge" :class="'badge-' + task.level">{{ levelText(task.level) }}</span>
              {{ task.title }}
            </h4>
            <p class="task-meta">
              <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
              <span class="points">+{{ task.points }}积分</span>
            </p>
          </div>
          <div class="task-actions">
            <button class="delete-btn" @click="deleteTask(task.id)">删除</button>
            <button @click="checkin(task.id)" :disabled="task.checked">
              {{ task.checked ? '已打卡' : '打卡' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi, checkinApi, userApi } from '../api'

const router = useRouter()
const user = ref(null)
const tasks = ref([])
const loading = ref(true)
const creating = ref(false)

const newTask = ref({
  title: '',
  circle_mode: 'workday',
  level: 1,
  points: 10
})

const circleModeText = (mode) => {
  const texts = {
    once: '单次',
    weekly: '每周',
    workday: '工作日',
    weekend: '周末',
    custom: '自定义'
  }
  return texts[mode] || mode
}

const levelText = (level) => {
  const texts = { 1: '低', 2: '中', 3: '高' }
  return texts[level] || '低'
}

const loadUser = async () => {
  const userId = parseInt(localStorage.getItem('userId'))
  if (!userId) {
    router.push('/login')
    return
  }
  try {
    const res = await userApi.getUser(userId)
    user.value = res.data
  } catch (e) {
    router.push('/login')
  }
}

const loadTasks = async () => {
  const userId = parseInt(localStorage.getItem('userId'))
  loading.value = true
  try {
    const res = await taskApi.getUserTasks(userId, { limit: 20, offset: 0 })
    const checkedRes = await checkinApi.getTodayChecked(userId)
    const checkedIds = checkedRes.data || []
    tasks.value = res.data.map(task => ({
      ...task,
      checked: checkedIds.includes(task.id)
    }))
  } catch (e) {
    console.error('加载任务失败', e)
  } finally {
    loading.value = false
  }
}

const createTask = async () => {
  creating.value = true
  const userId = parseInt(localStorage.getItem('userId'))
  try {
    const res = await taskApi.create({
      user_id: userId,
      title: newTask.value.title,
      circle_mode: newTask.value.circle_mode,
      level: newTask.value.level,
      points: newTask.value.points || 10
    })
    tasks.value.unshift({ ...res.data, checked: false })
    newTask.value.title = ''
  } catch (e) {
    alert('创建失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    creating.value = false
  }
}

const deleteTask = async (taskId) => {
  if (!confirm('确定删除该任务吗？')) return
  try {
    await taskApi.delete(taskId)
    tasks.value = tasks.value.filter(t => t.id !== taskId)
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const checkin = async (taskId) => {
  const userId = parseInt(localStorage.getItem('userId'))
  try {
    await checkinApi.checkin(taskId, { user_id: userId })
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = true
    const res = await userApi.getUser(userId)
    user.value = res.data
    alert('打卡成功！获得积分')
  } catch (e) {
    alert('打卡失败：' + (e.response?.data?.error || '未知错误'))
  }
}

onMounted(async () => {
  await loadUser()
  await loadTasks()
})
</script>

<style scoped>
.tasks-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}
h1 { color: #333; }
.user-info {
  display: flex;
  gap: 15px;
  align-items: center;
}
.points { color: #42b883; font-weight: bold; }
.wallet-btn {
  padding: 8px 16px;
  background: #42b883;
  color: white;
  border-radius: 4px;
  text-decoration: none;
}
.create-task {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}
.create-task h3 { margin-bottom: 15px; }
.create-task form {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.create-task input, .create-task select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.create-task input:first-child { flex: 1; min-width: 200px; }
.create-task button {
  padding: 10px 20px;
  background: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.task-list h3 { margin-bottom: 15px; }
.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #666;
}
.tasks {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border-left: 4px solid #42b883;
}
.task-item.level-1 { border-left-color: #42b883; }
.task-item.level-2 { border-left-color: #e6a23c; }
.task-item.level-3 { border-left-color: #f56c6c; }
.task-info h4 { margin: 0 0 5px 0; display: flex; align-items: center; gap: 8px; }
.level-badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}
.badge-1 { background: #e8f5e9; color: #42b883; }
.badge-2 { background: #fdf6ec; color: #e6a23c; }
.badge-3 { background: #fef0f0; color: #f56c6c; }
.task-meta {
  display: flex;
  gap: 10px;
  color: #666;
  font-size: 14px;
}
.task-meta .points { color: #42b883; }
.task-actions { display: flex; gap: 10px; }
.delete-btn {
  padding: 10px 15px;
  background: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.task-item button:not(.delete-btn) {
  padding: 10px 20px;
  background: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.task-item button:disabled { background: #ccc; }
</style>
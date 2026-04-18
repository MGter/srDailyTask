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
        <input v-model.number="newTask.points" type="number" placeholder="积分奖励" min="1" />
        <button type="submit" :disabled="creating">{{ creating ? '创建中...' : '创建' }}</button>
      </form>
    </div>

    <div class="task-list">
      <h3>我的任务</h3>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="tasks.length === 0" class="empty">暂无任务</div>
      <div v-else class="tasks">
        <div class="task-item" v-for="task in tasks" :key="task.id">
          <div class="task-info">
            <h4>{{ task.title }}</h4>
            <p class="task-meta">
              <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
              <span class="points">+{{ task.points }}积分</span>
            </p>
          </div>
          <button @click="checkin(task.id)" :disabled="task.checked">
            {{ task.checked ? '已打卡' : '打卡' }}
          </button>
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
    tasks.value = res.data.map(task => ({ ...task, checked: false }))
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

const checkin = async (taskId) => {
  const userId = parseInt(localStorage.getItem('userId'))
  try {
    await checkinApi.checkin(taskId, { user_id: userId })
    // 更新任务状态
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = true
    // 更新用户积分
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
h1 {
  color: #333;
}
.user-info {
  display: flex;
  gap: 15px;
  align-items: center;
}
.points {
  color: #42b883;
  font-weight: bold;
}
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
.create-task h3 {
  margin-bottom: 15px;
}
.create-task form {
  display: flex;
  gap: 10px;
}
.create-task input,
.create-task select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.create-task input:first-child {
  flex: 1;
}
.create-task button {
  padding: 10px 20px;
  background: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.task-list h3 {
  margin-bottom: 15px;
}
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
}
.task-info h4 {
  margin: 0 0 5px 0;
}
.task-meta {
  display: flex;
  gap: 10px;
  color: #666;
  font-size: 14px;
}
.task-meta .points {
  color: #42b883;
}
.task-item button {
  padding: 10px 20px;
  background: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.task-item button:disabled {
  background: #ccc;
}
</style>
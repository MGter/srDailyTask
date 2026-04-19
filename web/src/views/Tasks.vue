<template>
  <div class="dashboard-container">
    <header>
      <h1>每日打卡</h1>
      <div class="user-info">
        <span>{{ user?.username }}</span>
        <span class="points">积分: {{ user?.points || 0 }}</span>
        <button class="logout-btn" @click="logout">退出</button>
      </div>
    </header>

    <div class="today-info">
      <p>今天是 {{ todayText }}，以下任务需要打卡</p>
    </div>

    <!-- 左右布局 -->
    <div class="main-layout">
      <!-- 左侧：打卡任务 + 钱包 -->
      <div class="left-column">
        <!-- 打卡任务区域 -->
        <div class="section checkin-section">
          <div class="section-header">
            <h2>今日打卡任务</h2>
          </div>

          <div class="create-task">
            <form @submit.prevent="createTask">
              <input v-model="newTask.title" placeholder="任务标题" required />
              <select v-model="newTask.circle_mode">
                <option value="once">单次</option>
                <option value="weekly">每周一</option>
                <option value="workday">工作日</option>
                <option value="weekend">周末</option>
                <option value="custom">每天</option>
              </select>
              <select v-model.number="newTask.level">
                <option :value="1">● 低</option>
                <option :value="2">● 中</option>
                <option :value="3">● 高</option>
              </select>
              <input v-model.number="newTask.points" type="number" placeholder="积分" min="1" />
              <button type="submit" :disabled="creating">{{ creating ? '创建...' : '创建' }}</button>
            </form>
          </div>

          <div v-if="loadingTasks" class="loading">加载中...</div>
          <div v-else-if="todayTasks.length === 0" class="empty">今天没有需要打卡的任务</div>
          <div v-else class="tasks">
            <div class="task-item" v-for="task in todayTasks" :key="task.id" :class="'level-' + task.level">
              <div class="task-info">
                <h4>
                  <span class="level-dot" :class="'dot-' + task.level"></span>
                  {{ task.title }}
                </h4>
                <p class="task-meta">
                  <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
                  <span class="task-points">+{{ task.points }}</span>
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

        <!-- 积分钱包区域 -->
        <div class="section wallet-section">
          <div class="section-header">
            <h2>积分钱包</h2>
          </div>

          <div class="action-form">
            <form @submit.prevent="addRecord">
              <select v-model="addForm.type">
                <option value="earn">收入</option>
                <option value="spend">支出</option>
              </select>
              <input v-model.number="addForm.amount" type="number" placeholder="金额" min="1" required />
              <input v-model="addForm.description" placeholder="描述" required />
              <input v-model="addForm.record_time" type="datetime-local" required />
              <button type="submit" :disabled="adding">{{ adding ? '添加...' : '添加' }}</button>
            </form>
          </div>

          <div v-if="loadingHistory" class="loading">加载中...</div>
          <div v-else-if="history.length === 0" class="empty">暂无记录</div>
          <div v-else class="history-list">
            <div class="history-item" v-for="item in history" :key="item.id">
              <div class="item-info">
                <span :class="['type', item.type]">{{ item.type === 'earn' ? '收入' : '支出' }}</span>
                <span class="desc">{{ item.description }}</span>
                <span class="time">{{ formatDate(item.record_time) }}</span>
              </div>
              <div class="item-right">
                <span :class="['amount', item.type]">
                  {{ item.type === 'earn' ? '+' : '-' }}{{ item.amount }}
                </span>
                <button class="delete-btn small" @click="deleteRecord(item.id)">删除</button>
              </div>
            </div>
          </div>
        </div>

        <!-- 其他任务 -->
        <div class="section other-tasks-section" v-if="otherTasks.length > 0">
          <div class="section-header">
            <h2>其他任务</h2>
            <span class="sub-title">今天不需要打卡</span>
          </div>
          <div class="tasks">
            <div class="task-item inactive" v-for="task in otherTasks" :key="task.id" :class="'level-' + task.level">
              <div class="task-info">
                <h4>
                  <span class="level-dot" :class="'dot-' + task.level"></span>
                  {{ task.title }}
                </h4>
                <p class="task-meta">
                  <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
                  <span class="task-points">+{{ task.points }}</span>
                </p>
              </div>
              <div class="task-actions">
                <button class="delete-btn" @click="deleteTask(task.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：统计图表 -->
      <div class="right-column">
        <div class="section chart-section">
          <div class="section-header">
            <h2>积分统计</h2>
            <select v-model="chartDays" @change="loadDailyStats" class="days-select">
              <option value="7">7天</option>
              <option value="14">14天</option>
              <option value="30">30天</option>
              <option value="180">180天</option>
              <option value="360">360天</option>
            </select>
          </div>

          <div class="chart-wrapper">
            <div v-if="loadingStats" class="loading">加载中...</div>
            <canvas v-show="!loadingStats" ref="chartCanvas"></canvas>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi, checkinApi, userApi, walletApi, pointsApi } from '../api'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

const router = useRouter()
const user = ref(null)
const tasks = ref([])
const history = ref([])
const dailyStats = ref([])
const loadingTasks = ref(true)
const loadingHistory = ref(false)
const loadingStats = ref(false)
const creating = ref(false)
const adding = ref(false)
const chartDays = ref(7)
const chartCanvas = ref(null)
let chartInstance = null

const userId = parseInt(localStorage.getItem('userId'))

const newTask = ref({
  title: '',
  circle_mode: 'workday',
  level: 1,
  points: 10
})

const getDefaultTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

const addForm = ref({
  type: 'spend',
  amount: 10,
  description: '',
  record_time: getDefaultTime()
})

const todayText = computed(() => {
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const now = new Date()
  return `${now.getFullYear()}年${now.getMonth()+1}月${now.getDate()}日 ${days[now.getDay()]}`
})

const todayTasks = computed(() => {
  return tasks.value.filter(t => t.should_checkin_today && !t.is_expired)
})

const otherTasks = computed(() => {
  return tasks.value.filter(t => !t.should_checkin_today && !t.is_expired)
})

const circleModeText = (mode) => {
  const texts = {
    once: '单次',
    weekly: '每周一',
    workday: '工作日',
    weekend: '周末',
    custom: '每天'
  }
  return texts[mode] || mode
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const logout = () => {
  localStorage.removeItem('userId')
  router.push('/login')
}

const loadUser = async () => {
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
  loadingTasks.value = true
  try {
    const res = await taskApi.getUserTasks(userId, { limit: 50, offset: 0 })
    const checkedRes = await checkinApi.getTodayChecked(userId)
    const checkedIds = checkedRes.data || []
    tasks.value = res.data.map(task => ({
      ...task,
      checked: checkedIds.includes(task.id)
    }))
  } catch (e) {
    console.error('加载任务失败', e)
  } finally {
    loadingTasks.value = false
  }
}

const createTask = async () => {
  creating.value = true
  try {
    const res = await taskApi.create({
      user_id: userId,
      title: newTask.value.title,
      circle_mode: newTask.value.circle_mode,
      level: Number(newTask.value.level),
      points: Number(newTask.value.points) || 10
    })
    const checkedRes = await checkinApi.getTodayChecked(userId)
    const checkedIds = checkedRes.data || []
    tasks.value.unshift({
      ...res.data,
      checked: checkedIds.includes(res.data.id)
    })
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
  try {
    await checkinApi.checkin(taskId, { user_id: userId })
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = true
    if (task && task.circle_mode === 'once') {
      task.is_expired = true
    }
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await loadDailyStats()
  } catch (e) {
    alert('打卡失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const loadHistory = async () => {
  loadingHistory.value = true
  try {
    const res = await pointsApi.getHistory(userId, { limit: 30, offset: 0 })
    history.value = res.data
  } catch (e) {
    console.error('加载记录失败', e)
  } finally {
    loadingHistory.value = false
  }
}

const addRecord = async () => {
  adding.value = true
  try {
    await walletApi.addRecord({
      user_id: userId,
      type: addForm.value.type,
      amount: addForm.value.amount,
      description: addForm.value.description,
      record_time: new Date(addForm.value.record_time).toISOString()
    })
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await loadDailyStats()
    addForm.value.description = ''
    addForm.value.record_time = getDefaultTime()
  } catch (e) {
    alert('添加失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    adding.value = false
  }
}

const deleteRecord = async (id) => {
  if (!confirm('确定删除该记录吗？删除后会相应调整积分余额。')) return
  try {
    await walletApi.delete(id, userId)
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await loadDailyStats()
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const loadDailyStats = async () => {
  loadingStats.value = true
  try {
    const res = await pointsApi.getDailyStats(userId, { days: chartDays.value })
    dailyStats.value = res.data
    // 等待DOM更新后再渲染图表
    await nextTick()
    // 再次等待确保canvas已渲染
    setTimeout(() => {
      renderChart()
    }, 100)
  } catch (e) {
    console.error('加载统计失败', e)
  } finally {
    loadingStats.value = false
  }
}

const renderChart = () => {
  console.log('renderChart called, canvas:', chartCanvas.value, 'data:', dailyStats.value.length)
  if (!chartCanvas.value) {
    console.log('No canvas found')
    return
  }
  if (dailyStats.value.length === 0) {
    console.log('No data')
    return
  }

  if (chartInstance) {
    chartInstance.destroy()
  }

  const ctx = chartCanvas.value.getContext('2d')
  const labels = dailyStats.value.map(s => s.date.slice(5))
  const earnData = dailyStats.value.map(s => s.earn)
  const spendData = dailyStats.value.map(s => s.spend)
  const balanceData = dailyStats.value.map(s => s.balance)

  console.log('Chart labels:', labels, 'earnData:', earnData)

  chartInstance = new Chart(ctx, {
    type: 'bar',
    data: {
      labels,
      datasets: [
        {
          label: '当日积累',
          data: earnData,
          backgroundColor: '#34c759',
          borderRadius: 6
        },
        {
          label: '当日消耗',
          data: spendData,
          backgroundColor: '#ff3b30',
          borderRadius: 6
        },
        {
          label: '累计结余',
          data: balanceData,
          backgroundColor: '#ff9500',
          borderRadius: 6
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top',
          labels: {
            color: '#1d1d1f'
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            color: '#1d1d1f'
          },
          grid: {
            color: '#e5e5e5'
          }
        },
        x: {
          ticks: {
            color: '#1d1d1f'
          },
          grid: {
            display: false
          }
        }
      }
    }
  })
}

onMounted(async () => {
  await loadUser()
  await loadTasks()
  await loadHistory()
  await loadDailyStats()
})
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  padding: 20px;
  background-image: url('/assets/kita.png');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding: 12px 20px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

h1 {
  color: #1d1d1f;
  margin: 0;
  font-weight: 600;
  font-size: 20px;
}

.user-info {
  display: flex;
  gap: 15px;
  align-items: center;
}

.points {
  color: #34c759;
  font-weight: 600;
}

.logout-btn {
  padding: 8px 16px;
  background: #ff3b30;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
}

.today-info {
  background: #34c759;
  padding: 10px 16px;
  border-radius: 10px;
  margin-bottom: 15px;
  color: white;
  font-weight: 500;
}

/* 左右布局 */
.main-layout {
  display: flex;
  gap: 20px;
}

.left-column {
  flex: 1;
  min-width: 0;
}

.right-column {
  width: 380px;
  flex-shrink: 0;
}

.section {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  padding: 18px;
  margin-bottom: 15px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-header h2 {
  margin: 0;
  color: #1d1d1f;
  font-weight: 600;
  font-size: 16px;
}

.sub-title {
  color: #86868b;
  font-size: 13px;
}

/* 创建任务表单 */
.create-task form, .action-form form {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.create-task input, .create-task select,
.action-form input, .action-form select {
  padding: 8px 12px;
  border: 1px solid #d2d2d7;
  border-radius: 8px;
  font-size: 14px;
  color: #1d1d1f;
  background: #f5f5f7;
}

.create-task input:focus, .create-task select:focus,
.action-form input:focus, .action-form select:focus {
  border-color: #007aff;
  outline: none;
  background: #ffffff;
}

.create-task input:first-child { flex: 1; min-width: 140px; }
.action-form input[placeholder="描述"] { flex: 1; min-width: 100px; }

.create-task button, .action-form button {
  padding: 8px 16px;
  background: #007aff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
}

.create-task button:disabled, .action-form button:disabled {
  background: #c7c7cc;
}

/* 任务列表 */
.tasks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: #f5f5f7;
  border-radius: 10px;
  border-left: 3px solid #34c759;
}

.task-item.inactive {
  opacity: 0.6;
}

.task-item.level-1 { border-left-color: #34c759; }
.task-item.level-2 { border-left-color: #ff9500; }
.task-item.level-3 { border-left-color: #ff3b30; }

.task-info h4 {
  margin: 0 0 4px 0;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #1d1d1f;
  font-size: 14px;
}

.level-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot-1 { background: #34c759; }
.dot-2 { background: #ff9500; }
.dot-3 { background: #ff3b30; }

.task-meta {
  display: flex;
  gap: 8px;
  color: #86868b;
  font-size: 12px;
}

.task-points {
  color: #34c759;
  font-weight: 500;
}

.task-actions { display: flex; gap: 6px; }

.delete-btn {
  padding: 6px 12px;
  background: #ff3b30;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
}

.delete-btn.small { padding: 5px 8px; }

.task-item button:not(.delete-btn) {
  padding: 6px 14px;
  background: #007aff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
}

.task-item button:disabled { background: #c7c7cc; }

/* 积分历史 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 250px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f5f5f7;
  border-radius: 8px;
}

.item-info {
  display: flex;
  gap: 8px;
  align-items: center;
}

.type {
  padding: 3px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.type.earn {
  background: #34c759;
  color: white;
}

.type.spend {
  background: #ff3b30;
  color: white;
}

.desc { color: #1d1d1f; font-size: 13px; }
.time { color: #86868b; font-size: 11px; }

.item-right {
  display: flex;
  gap: 8px;
  align-items: center;
}

.amount { font-weight: 600; font-size: 13px; }
.amount.earn { color: #34c759; }
.amount.spend { color: #ff3b30; }

/* 图表 */
.days-select {
  padding: 5px 10px;
  border: 1px solid #d2d2d7;
  border-radius: 6px;
  font-size: 13px;
  color: #1d1d1f;
  background: #f5f5f7;
}

.chart-wrapper {
  height: 300px;
  position: relative;
  border-radius: 10px;
  padding: 15px;
}

.chart-wrapper canvas {
  width: 100% !important;
  height: 100% !important;
}

.loading, .empty {
  text-align: center;
  padding: 20px;
  color: #86868b;
}

/* 响应式 */
@media (max-width: 900px) {
  .main-layout {
    flex-direction: column;
  }
  .right-column {
    width: 100%;
  }
}
</style>
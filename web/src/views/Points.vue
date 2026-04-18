<template>
  <div class="points-container">
    <header>
      <h1>积分钱包</h1>
      <router-link to="/tasks" class="back-btn">返回任务</router-link>
    </header>

    <div class="balance-card">
      <h2>当前积分</h2>
      <p class="balance">{{ balance }}</p>
    </div>

    <div class="action-section">
      <h3>新增记录</h3>
      <form @submit.prevent="addRecord">
        <select v-model="addForm.type">
          <option value="earn">收入</option>
          <option value="spend">支出</option>
        </select>
        <input v-model.number="addForm.amount" type="number" placeholder="金额" min="1" required />
        <input v-model="addForm.description" placeholder="描述" required />
        <input v-model="addForm.record_time" type="datetime-local" required />
        <button type="submit" class="add-btn">{{ adding ? '添加中...' : '添加' }}</button>
      </form>
    </div>

    <div class="history-section">
      <h3>积分记录</h3>
      <div v-if="loading" class="loading">加载中...</div>
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
            <button class="delete-btn" @click="deleteRecord(item.id)">删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { walletApi, pointsApi } from '../api'

const router = useRouter()
const balance = ref(0)
const history = ref([])
const loading = ref(true)
const adding = ref(false)

const userId = parseInt(localStorage.getItem('userId'))

// 默认时间为当前时间
const getDefaultTime = () => {
  const now = new Date()
  return now.toISOString().slice(0, 16)
}

const addForm = ref({
  type: 'earn',
  amount: 10,
  description: '',
  record_time: getDefaultTime()
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const loadBalance = async () => {
  try {
    const res = await walletApi.getBalance(userId)
    balance.value = res.data.balance || 0
  } catch (e) {
    console.error('加载余额失败', e)
  }
}

const loadHistory = async () => {
  loading.value = true
  try {
    const res = await pointsApi.getHistory(userId, { limit: 50, offset: 0 })
    history.value = res.data
  } catch (e) {
    console.error('加载记录失败', e)
  } finally {
    loading.value = false
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
    await loadBalance()
    await loadHistory()
    addForm.value.description = ''
    addForm.value.record_time = getDefaultTime()
    alert('添加成功')
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
    await loadBalance()
    await loadHistory()
    alert('删除成功')
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

onMounted(async () => {
  if (!userId) {
    router.push('/login')
    return
  }
  await loadBalance()
  await loadHistory()
})
</script>

<style scoped>
.points-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}
.back-btn {
  padding: 8px 16px;
  background: #666;
  color: white;
  border-radius: 4px;
  text-decoration: none;
}
.balance-card {
  background: linear-gradient(135deg, #42b883, #35495e);
  color: white;
  padding: 30px;
  border-radius: 12px;
  text-align: center;
  margin-bottom: 30px;
}
.balance-card h2 { margin-bottom: 10px; }
.balance { font-size: 48px; font-weight: bold; }
.action-section {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}
.action-section h3 { margin-bottom: 15px; }
.action-section form {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.action-section input, .action-section select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.action-section input[type="number"] { width: 80px; }
.action-section input[type="datetime-local"] { width: 180px; }
.action-section input[placeholder="描述"] { flex: 1; min-width: 150px; }
.add-btn {
  padding: 10px 20px;
  background: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.history-section h3 { margin-bottom: 15px; }
.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #666;
}
.history-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}
.item-info {
  display: flex;
  gap: 15px;
  align-items: center;
}
.type {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.type.earn { background: #e8f5e9; color: #42b883; }
.type.spend { background: #ffebee; color: #f56c6c; }
.desc { color: #333; }
.time { color: #999; font-size: 12px; }
.item-right {
  display: flex;
  gap: 15px;
  align-items: center;
}
.amount { font-weight: bold; font-size: 18px; }
.amount.earn { color: #42b883; }
.amount.spend { color: #f56c6c; }
.delete-btn {
  padding: 6px 12px;
  background: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}
</style>
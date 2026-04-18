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

    <div class="spend-section">
      <h3>消费积分</h3>
      <form @submit.prevent="spendPoints">
        <input v-model.number="spendForm.amount" type="number" placeholder="消费金额" min="1" required />
        <input v-model="spendForm.description" placeholder="消费原因" required />
        <button type="submit" :disabled="spending">{{ spending ? '处理中...' : '确认消费' }}</button>
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
            <span class="time">{{ formatDate(item.created_at) }}</span>
          </div>
          <span :class="['amount', item.type]">
            {{ item.type === 'earn' ? '+' : '-' }}{{ item.amount }}
          </span>
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
const spending = ref(false)

const spendForm = ref({
  amount: 1,
  description: ''
})

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
}

const userId = localStorage.getItem('userId')

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
    const res = await pointsApi.getHistory(userId, { limit: 20, offset: 0 })
    history.value = res.data
  } catch (e) {
    console.error('加载记录失败', e)
  } finally {
    loading.value = false
  }
}

const spendPoints = async () => {
  if (spendForm.value.amount > balance.value) {
    alert('积分不足')
    return
  }
  spending.value = true
  try {
    await walletApi.spend({
      user_id: userId,
      amount: spendForm.value.amount,
      description: spendForm.value.description
    })
    await loadBalance()
    await loadHistory()
    spendForm.value.description = ''
    alert('消费成功')
  } catch (e) {
    alert('消费失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    spending.value = false
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
.balance-card h2 {
  margin-bottom: 10px;
}
.balance {
  font-size: 48px;
  font-weight: bold;
}
.spend-section {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}
.spend-section h3 {
  margin-bottom: 15px;
}
.spend-section form {
  display: flex;
  gap: 10px;
}
.spend-section input {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.spend-section input:first-child {
  width: 100px;
}
.spend-section input:nth-child(2) {
  flex: 1;
}
.spend-section button {
  padding: 10px 20px;
  background: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.history-section h3 {
  margin-bottom: 15px;
}
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
.type.earn {
  background: #e8f5e9;
  color: #42b883;
}
.type.spend {
  background: #ffebee;
  color: #f56c6c;
}
.desc {
  color: #333;
}
.time {
  color: #999;
  font-size: 12px;
}
.amount {
  font-weight: bold;
  font-size: 18px;
}
.amount.earn {
  color: #42b883;
}
.amount.spend {
  color: #f56c6c;
}
</style>
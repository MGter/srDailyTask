<template>
  <div class="login-container">
    <div class="login-box">
      <h2>登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" type="text" required placeholder="请输入用户名" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="form.password" type="password" required placeholder="请输入密码" />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        <p class="error" v-if="error">{{ error }}</p>
        <p class="link">
          还没有账号？<router-link to="/register">注册</router-link>
        </p>
      </form>
    </div>
    <footer class="footer">
      <span>© 2026 srDailyTask</span>
      <a href="https://github.com/MGter/srDailyTask" target="_blank">GitHub</a>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { userApi } from '../api'

const router = useRouter()
const form = ref({ username: '', password: '' })
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await userApi.login(form.value)
    localStorage.setItem('user', JSON.stringify(res.data))
    localStorage.setItem('userId', res.data.id)
    router.push('/tasks')
  } catch (e) {
    error.value = e.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-image: url('/assets/kita.webp');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}
.login-box {
  background: #ffffff;
  padding: 40px 35px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 340px;
}
h2 {
  text-align: center;
  color: #1d1d1f;
  margin-bottom: 30px;
  font-weight: 600;
  font-size: 24px;
}
.form-group {
  margin-bottom: 20px;
}
label {
  display: block;
  margin-bottom: 6px;
  color: #1d1d1f;
  font-weight: 500;
  font-size: 14px;
}
input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #d2d2d7;
  border-radius: 10px;
  font-size: 15px;
  color: #1d1d1f;
  background: #f5f5f7;
  box-sizing: border-box;
}
input:focus {
  border-color: #007aff;
  outline: none;
  background: #ffffff;
}
button {
  width: 100%;
  padding: 14px;
  background: #007aff;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}
button:hover {
  background: #0066d6;
}
button:disabled {
  background: #c7c7cc;
}
.error {
  color: #ff3b30;
  text-align: center;
  margin-top: 15px;
  font-weight: 500;
}
.link {
  text-align: center;
  margin-top: 25px;
  color: #86868b;
}
.link a {
  color: #007aff;
  font-weight: 500;
}
.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  text-align: center;
  padding: 15px 20px;
  color: #86868b;
  font-size: 12px;
  display: flex;
  justify-content: center;
  gap: 20px;
}
.footer a {
  color: #007aff;
  text-decoration: none;
}
</style>
<template>
  <div class="register-container">
    <div class="register-box">
      <h2>注册</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" type="text" required placeholder="请输入用户名" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="form.password" type="password" required placeholder="请输入密码" />
        </div>
        <div class="form-group">
          <label>邮箱</label>
          <input v-model="form.email" type="email" placeholder="请输入邮箱（可选）" />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
        <p class="error" v-if="error">{{ error }}</p>
        <p class="success" v-if="success">{{ success }}</p>
        <p class="link">
          已有账号？<router-link to="/login">登录</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { userApi } from '../api'

const router = useRouter()
const form = ref({ username: '', password: '', email: '' })
const loading = ref(false)
const error = ref('')
const success = ref('')

const handleRegister = async () => {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    const res = await userApi.register(form.value)
    success.value = '注册成功！'
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (e) {
    error.value = e.response?.data?.error || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-image: url('/assets/kita.webp');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}
.register-box {
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
.success {
  color: #34c759;
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
</style>
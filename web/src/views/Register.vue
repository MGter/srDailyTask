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
  background-image: url('@/assets/kita.png');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}
.register-box {
  background: rgba(255, 255, 255, 0.92);
  padding: 40px 35px;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  width: 320px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}
h2 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 30px;
  font-weight: 600;
}
.form-group {
  margin-bottom: 20px;
}
label {
  display: block;
  margin-bottom: 8px;
  color: #2c3e50;
  font-weight: 500;
}
input {
  width: 100%;
  padding: 12px 14px;
  border: 2px solid #dce6e9;
  border-radius: 8px;
  font-size: 14px;
  background: rgba(255, 255, 255, 0.9);
  transition: border-color 0.2s;
}
input:focus {
  border-color: #42b883;
  outline: none;
}
button {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #42b883, #35495e);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}
button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 184, 131, 0.4);
}
button:disabled {
  background: #a0a0a0;
  transform: none;
  box-shadow: none;
}
.error {
  color: #f56c6c;
  text-align: center;
  margin-top: 15px;
  font-weight: 500;
}
.success {
  color: #42b883;
  text-align: center;
  margin-top: 15px;
  font-weight: 500;
}
.link {
  text-align: center;
  margin-top: 25px;
  color: #7f8c8d;
}
.link a {
  color: #42b883;
  font-weight: 500;
}
</style>
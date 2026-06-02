<template>
  <div class="login-container" :style="{ backgroundImage: `url(/images/login-bg.png)` }">
    <div class="login-content">
      <div class="login-left">
        <img src="/images/login-illustration.png" alt="login" class="login-illustration" />
      </div>
      <el-card class="login-card">
        <div class="login-header">
          <img src="/images/logo.png" alt="logo" class="login-logo" />
          <h2>题小助管理后台</h2>
        </div>
        <el-form ref="formRef" :model="form" :rules="rules" @keyup.enter="handleLogin">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" size="large" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" show-password size="large" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" size="large" style="width: 100%" @click="handleLogin">
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useUserStore } from "@/stores/user";
import type { FormInstance } from "element-plus";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const formRef = ref<FormInstance>();
const loading = ref(false);

const form = reactive({
  username: "",
  password: "",
});

const rules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

async function handleLogin() {
  await formRef.value?.validate();
  loading.value = true;
  try {
    await userStore.login(form.username, form.password);
    const redirect = (route.query.redirect as string) || "/";
    router.push(redirect);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-size: cover;
  background-position: center;
}
.login-content {
  display: flex;
  align-items: center;
  gap: 60px;
}
.login-left {
  flex-shrink: 0;
}
.login-illustration {
  width: 300px;
  height: auto;
}
.login-card {
  width: 400px;
}
.login-header {
  text-align: center;
  margin-bottom: 30px;
}
.login-logo {
  width: 64px;
  height: 64px;
  margin-bottom: 12px;
}
h2 {
  margin: 0;
  font-size: 22px;
  color: #303133;
}
</style>

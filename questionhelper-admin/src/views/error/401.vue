<template>
  <div class="error-page">
    <h1>401</h1>
    <p>未授权，请登录</p>
    <el-button type="primary" @click="$router.push('/login')">去登录</el-button>
  </div>
</template>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}
h1 {
  font-size: 96px;
  color: #f56c6c;
}
</style>

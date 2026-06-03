<template>
  <div class="error-page">
    <div class="error-content">
      <img src="/images/404.png" alt="404" class="error-img" />
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: "Page404" });
</script>

<style scoped>
.error-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #f0f2f5;
}
.error-content {
  text-align: center;
}
.error-img {
  width: 300px;
  margin-bottom: 20px;
}
h1 {
  font-size: 72px;
  color: #409eff;
  margin-bottom: 10px;
}
p {
  font-size: 18px;
  color: #606266;
  margin-bottom: 30px;
}
</style>

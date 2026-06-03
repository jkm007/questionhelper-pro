<template>
  <div class="error-page">
    <div class="error-content">
      <img src="/images/401.png" alt="401" class="error-img" />
      <h1>401</h1>
      <p>抱歉，您没有权限访问此页面</p>
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: "Page401" });
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
  color: #e6a23c;
  margin-bottom: 10px;
}
p {
  font-size: 18px;
  color: #606266;
  margin-bottom: 30px;
}
</style>

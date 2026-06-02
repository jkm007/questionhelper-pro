<template>
  <div class="error-page">
    <img src="/images/404.png" alt="404" class="error-image" />
    <p class="error-desc">抱歉，您访问的页面不存在</p>
    <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
  </div>
</template>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: #f5f7fa;
}
.error-image {
  width: 300px;
  height: auto;
  margin-bottom: 24px;
}
.error-desc {
  font-size: 16px;
  color: #606266;
  margin-bottom: 24px;
}
</style>

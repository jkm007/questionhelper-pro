<template>
  <div class="dashboard-container">
    <h2>仪表盘</h2>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总用户数</template>
          <div class="stat-value">{{ stats.totalUsers || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总题目数</template>
          <div class="stat-value">{{ stats.totalQuestions || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总考试数</template>
          <div class="stat-value">{{ stats.totalExams || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>今日活跃</template>
          <div class="stat-value">{{ stats.todayActive || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getOverviewApi } from "@/api/statistics";

const stats = ref<any>({});

onMounted(async () => {
  try {
    stats.value = await getOverviewApi();
  } catch {
    // ignore
  }
});
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
}
</style>

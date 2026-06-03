<template>
  <el-config-provider :locale="locale" :size="size">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import { useAppStore } from "@/stores";
import zhCn from "element-plus/es/locale/lang/zh-cn";

const appStore = useAppStore();

const locale = computed(() => zhCn);
const size = computed(() => appStore.size as any);
</script>

import { store } from "@/stores";

export const useDictStore = defineStore("dict", () => {
  // 字典数据缓存
  const dictCache = ref<Record<string, any[]>>({});

  /**
   * 加载字典数据
   */
  const loadDictItems = async (dictCode: string) => {
    // 暂时返回空数组
    return [];
  };

  /**
   * 获取字典项列表
   */
  const getDictItems = (dictCode: string) => {
    return dictCache.value[dictCode] || [];
  };

  /**
   * 清空字典缓存
   */
  const clearDictCache = () => {
    dictCache.value = {};
  };

  return {
    dictCache,
    loadDictItems,
    getDictItems,
    clearDictCache,
  };
});

export function useDictStoreHook() {
  return useDictStore(store);
}

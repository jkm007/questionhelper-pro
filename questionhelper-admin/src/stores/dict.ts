import { defineStore } from "pinia";
import { ref } from "vue";

export const useDictStore = defineStore("dict", () => {
  const dictMap = ref<Map<string, any[]>>(new Map());

  function setDict(code: string, data: any[]) {
    dictMap.value.set(code, data);
  }

  function getDict(code: string): any[] | undefined {
    return dictMap.value.get(code);
  }

  function clearDict() {
    dictMap.value.clear();
  }

  return { dictMap, setDict, getDict, clearDict };
});

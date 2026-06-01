import { defineStore } from "pinia";
import { ref } from "vue";

interface TagView {
  path: string;
  title: string;
  name: string;
}

export const useTagsViewStore = defineStore("tags-view", () => {
  const visitedViews = ref<TagView[]>([]);
  const cachedViews = ref<string[]>([]);

  function addView(view: TagView) {
    if (!visitedViews.value.find((v) => v.path === view.path)) {
      visitedViews.value.push(view);
    }
    if (!cachedViews.value.includes(view.name)) {
      cachedViews.value.push(view.name);
    }
  }

  function removeView(path: string) {
    const index = visitedViews.value.findIndex((v) => v.path === path);
    if (index > -1) {
      const view = visitedViews.value[index];
      visitedViews.value.splice(index, 1);
      const cacheIndex = cachedViews.value.indexOf(view.name);
      if (cacheIndex > -1) {
        cachedViews.value.splice(cacheIndex, 1);
      }
    }
  }

  function clearViews() {
    visitedViews.value = [];
    cachedViews.value = [];
  }

  return { visitedViews, cachedViews, addView, removeView, clearViews };
});

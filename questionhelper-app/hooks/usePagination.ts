import { ref, reactive } from 'vue'

interface PaginationOptions {
  pageSize?: number
  immediate?: boolean
}

export const usePagination = <T = any>(
  fetchFn: (params: any) => Promise<any>,
  options: PaginationOptions = {}
) => {
  const { pageSize = 10, immediate = false } = options

  const list = ref<T[]>([]) as any
  const loading = ref(false)
  const finished = ref(false)
  const refreshing = ref(false)

  const pagination = reactive({
    page: 1,
    pageSize,
    total: 0
  })

  const loadData = async (isRefresh = false) => {
    if (loading.value) return

    if (isRefresh) {
      pagination.page = 1
      finished.value = false
    }

    loading.value = true

    try {
      const res = await fetchFn({
        page: pagination.page,
        page_size: pagination.pageSize
      })

      const { list: newList, total } = res.data

      if (isRefresh) {
        list.value = newList || []
      } else {
        list.value = [...list.value, ...(newList || [])]
      }

      pagination.total = total || 0

      if (list.value.length >= pagination.total) {
        finished.value = true
      } else {
        pagination.page++
      }
    } catch (e) {
      console.error('加载数据失败', e)
    } finally {
      loading.value = false
      refreshing.value = false
    }
  }

  const onRefresh = () => {
    refreshing.value = true
    loadData(true)
  }

  const onLoadMore = () => {
    if (!finished.value) {
      loadData()
    }
  }

  if (immediate) {
    loadData()
  }

  return {
    list,
    loading,
    finished,
    refreshing,
    pagination,
    loadData,
    onRefresh,
    onLoadMore
  }
}

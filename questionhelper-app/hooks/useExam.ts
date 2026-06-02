import { ref, computed, onBeforeUnmount } from 'vue'

export const useExam = () => {
  const duration = ref(0) // 考试时长（秒）
  const remaining = ref(0) // 剩余时间（秒）
  const isRunning = ref(false)
  let timer: ReturnType<typeof setInterval> | null = null
  let syncTimer: ReturnType<typeof setInterval> | null = null
  let serverTimestamp = 0
  let localTimestamp = 0

  // 格式化时间
  const formattedTime = computed(() => {
    const hours = Math.floor(remaining.value / 3600)
    const minutes = Math.floor((remaining.value % 3600) / 60)
    const seconds = remaining.value % 60

    if (hours > 0) {
      return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
    }
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
  })

  // 是否时间紧迫（小于5分钟）
  const isUrgent = computed(() => remaining.value < 300 && remaining.value > 0)

  // 同步服务端时间（预留接口）
  const syncWithServer = () => {
    // TODO: 调用后端获取服务端时间戳
  }

  // 开始计时
  const start = (seconds: number, serverTs?: number) => {
    duration.value = seconds
    remaining.value = seconds
    isRunning.value = true

    if (serverTs) {
      serverTimestamp = serverTs
      localTimestamp = Date.now()
      syncTimer = setInterval(() => {
        syncWithServer()
      }, 60000)
    }

    timer = setInterval(() => {
      if (remaining.value > 0) {
        remaining.value--
      } else {
        stop()
      }
    }, 1000)
  }

  // 暂停计时
  const pause = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    isRunning.value = false
  }

  // 继续计时
  const resume = () => {
    if (!isRunning.value && remaining.value > 0) {
      isRunning.value = true
      timer = setInterval(() => {
        if (remaining.value > 0) {
          remaining.value--
        } else {
          stop()
        }
      }, 1000)
    }
  }

  // 停止计时
  const stop = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    if (syncTimer) {
      clearInterval(syncTimer)
      syncTimer = null
    }
    isRunning.value = false
  }

  // 重置计时
  const reset = () => {
    stop()
    remaining.value = duration.value
  }

  // 获取已用时间
  const getElapsedTime = () => {
    return duration.value - remaining.value
  }

  // 获取剩余时间
  const getRemainingTime = () => {
    return remaining.value
  }

  // 组件销毁时清理
  onBeforeUnmount(() => {
    stop()
  })

  return {
    duration,
    remaining,
    isRunning,
    formattedTime,
    isUrgent,
    start,
    pause,
    resume,
    stop,
    reset,
    getElapsedTime,
    getRemainingTime
  }
}

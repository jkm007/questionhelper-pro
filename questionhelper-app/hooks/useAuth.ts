import { computed } from 'vue'
import { useUserStore } from '@/store/modules/user'
import { useRouter } from 'uni-mini-router'

export const useAuth = () => {
  const userStore = useUserStore()

  const isLoggedIn = computed(() => userStore.isLoggedIn)
  const userInfo = computed(() => userStore.userInfo)
  const userId = computed(() => userStore.userId)
  const nickname = computed(() => userStore.nickname)
  const avatar = computed(() => userStore.avatar)

  const requireAuth = (callback?: () => void) => {
    if (!isLoggedIn.value) {
      uni.navigateTo({ url: '/pages/login/index' })
      return false
    }
    callback?.()
    return true
  }

  const logout = async () => {
    uni.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: async (res) => {
        if (res.confirm) {
          await userStore.logout()
          uni.reLaunch({ url: '/pages/login/index' })
        }
      }
    })
  }

  return {
    isLoggedIn,
    userInfo,
    userId,
    nickname,
    avatar,
    requireAuth,
    logout
  }
}

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePracticeStore = defineStore('practice', () => {
  const currentSession = ref<any>(null)
  const currentIndex = ref<number>(0)
  const answers = ref<Record<number, { answer: string; isCorrect: boolean; duration: number }>>({})

  // 设置当前练习
  const setCurrentSession = (session: any) => {
    currentSession.value = session
    currentIndex.value = 0
    answers.value = {}
  }

  // 设置当前题目索引
  const setCurrentIndex = (index: number) => {
    currentIndex.value = index
  }

  // 记录答案
  const recordAnswer = (questionId: number, answer: string, isCorrect: boolean, duration: number) => {
    answers.value[questionId] = { answer, isCorrect, duration }
  }

  // 获取正确数量
  const getCorrectCount = (): number => {
    return Object.values(answers.value).filter(a => a.isCorrect).length
  }

  // 清空练习状态
  const clearPractice = () => {
    currentSession.value = null
    currentIndex.value = 0
    answers.value = {}
  }

  return {
    currentSession,
    currentIndex,
    answers,
    setCurrentSession,
    setCurrentIndex,
    recordAnswer,
    getCorrectCount,
    clearPractice
  }
})

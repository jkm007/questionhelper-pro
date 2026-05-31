import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useExamStore = defineStore('exam', () => {
  const currentExam = ref<any>(null)
  const answers = ref<Record<number, string>>({})
  const startTime = ref<number>(0)

  // 设置当前考试
  const setCurrentExam = (exam: any) => {
    currentExam.value = exam
    answers.value = {}
    startTime.value = Date.now()
  }

  // 设置答案
  const setAnswer = (questionId: number, answer: string) => {
    answers.value[questionId] = answer
  }

  // 获取答案
  const getAnswer = (questionId: number): string => {
    return answers.value[questionId] || ''
  }

  // 清空考试状态
  const clearExam = () => {
    currentExam.value = null
    answers.value = {}
    startTime.value = 0
  }

  return {
    currentExam,
    answers,
    startTime,
    setCurrentExam,
    setAnswer,
    getAnswer,
    clearExam
  }
})

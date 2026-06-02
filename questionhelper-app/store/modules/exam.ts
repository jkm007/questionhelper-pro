import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useExamStore = defineStore('exam', () => {
  const currentExam = ref<any>(null)
  const examId = ref<string>('')
  const answers = ref<Record<string, number[]>>({})
  const flaggedQuestions = ref<Set<string>>(new Set())
  const startTime = ref<number>(0)
  const currentIndex = ref<number>(0)

  const setCurrentExam = (exam: any) => {
    currentExam.value = exam
    examId.value = exam?.id?.toString() || ''
    answers.value = {}
    flaggedQuestions.value = new Set()
    startTime.value = Date.now()
    currentIndex.value = 0
  }

  const setAnswer = (questionId: string, selectedOptions: number[]) => {
    answers.value[questionId] = selectedOptions
    saveAnswersToLocal()
  }

  const getAnswer = (questionId: string): number[] => {
    return answers.value[questionId] || []
  }

  const setCurrentIndex = (index: number) => {
    currentIndex.value = index
  }

  const toggleFlag = (questionId: string) => {
    const newSet = new Set(flaggedQuestions.value)
    if (newSet.has(questionId)) {
      newSet.delete(questionId)
    } else {
      newSet.add(questionId)
    }
    flaggedQuestions.value = newSet
  }

  const saveAnswersToLocal = () => {
    if (!examId.value) return
    try {
      uni.setStorageSync(`exam_answers_${examId.value}`, JSON.stringify({
        examId: examId.value,
        answers: answers.value,
        currentIndex: currentIndex.value,
        flaggedQuestions: Array.from(flaggedQuestions.value),
        timestamp: Date.now()
      }))
    } catch (e) {
      console.error('保存答案到本地失败', e)
    }
  }

  const loadAnswersFromLocal = (id: string) => {
    try {
      const raw = uni.getStorageSync(`exam_answers_${id}`)
      if (!raw) return null
      const data = JSON.parse(raw)
      return data
    } catch {
      return null
    }
  }

  const clearLocalCache = (id?: string) => {
    const key = id || examId.value
    if (key) {
      uni.removeStorageSync(`exam_answers_${key}`)
    }
  }

  const clearExam = () => {
    clearLocalCache()
    currentExam.value = null
    examId.value = ''
    answers.value = {}
    flaggedQuestions.value = new Set()
    startTime.value = 0
    currentIndex.value = 0
  }

  return {
    currentExam, examId, answers, flaggedQuestions, startTime, currentIndex,
    setCurrentExam, setAnswer, getAnswer, setCurrentIndex, toggleFlag,
    saveAnswersToLocal, loadAnswersFromLocal, clearLocalCache, clearExam
  }
})

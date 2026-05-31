<template>
  <view class="edit-container">
    <!-- 题目类型 -->
    <view class="form-section">
      <text class="section-title">基本信息</text>
      <view class="form-item">
        <text class="form-label">题目类型</text>
        <view class="type-selector">
          <view
            v-for="item in typeOptions"
            :key="item.value"
            class="type-option"
            :class="{ active: form.type === item.value }"
            @tap="form.type = item.value"
          >
            <text>{{ item.label }}</text>
          </view>
        </view>
      </view>
      <view class="form-item">
        <text class="form-label">难度等级</text>
        <view class="difficulty-selector">
          <view
            v-for="item in difficultyOptions"
            :key="item.value"
            class="difficulty-option"
            :class="{ active: form.difficulty === item.value }"
            @tap="form.difficulty = item.value"
          >
            <text>{{ item.label }}</text>
          </view>
        </view>
      </view>
      <view class="form-item">
        <text class="form-label">所属分类</text>
        <view class="picker-wrap" @tap="showCategoryPicker = true">
          <text :class="{ placeholder: !form.categoryName }">
            {{ form.categoryName || '请选择分类' }}
          </text>
          <text class="picker-arrow">></text>
        </view>
      </view>
    </view>

    <!-- 题目内容 -->
    <view class="form-section">
      <text class="section-title">题目内容</text>
      <view class="form-item">
        <text class="form-label">题目标题</text>
        <input
          v-model="form.title"
          class="form-input"
          placeholder="请输入题目标题"
          maxlength="200"
        />
      </view>
      <view class="form-item">
        <text class="form-label">题目内容</text>
        <textarea
          v-model="form.content"
          class="form-textarea"
          placeholder="请输入题目内容"
          maxlength="5000"
          :auto-height="true"
        />
      </view>
    </view>

    <!-- 选项 (单选/多选/判断题) -->
    <view v-if="showOptions" class="form-section">
      <view class="section-header">
        <text class="section-title">选项设置</text>
        <text class="add-btn" @tap="addOption">+ 添加选项</text>
      </view>
      <view
        v-for="(option, index) in form.options"
        :key="index"
        class="option-item"
      >
        <view class="option-header">
          <text class="option-label">{{ optionLabels[index] }}</text>
          <view
            class="correct-toggle"
            :class="{ active: form.correctOptions.includes(index) }"
            @tap="toggleCorrect(index)"
          >
            <text>{{ form.correctOptions.includes(index) ? '✓ 正确' : '设为正确' }}</text>
          </view>
          <text v-if="form.options.length > 2" class="delete-btn" @tap="removeOption(index)">删除</text>
        </view>
        <input
          v-model="option.content"
          class="form-input"
          :placeholder="`请输入选项${optionLabels[index]}内容`"
        />
      </view>
    </view>

    <!-- 答案 & 解析 -->
    <view class="form-section">
      <text class="section-title">答案与解析</text>
      <view class="form-item">
        <text class="form-label">参考答案</text>
        <textarea
          v-model="form.answer"
          class="form-textarea"
          placeholder="请输入参考答案"
          maxlength="5000"
          :auto-height="true"
        />
      </view>
      <view class="form-item">
        <text class="form-label">题目解析</text>
        <textarea
          v-model="form.analysis"
          class="form-textarea"
          placeholder="请输入题目解析"
          maxlength="5000"
          :auto-height="true"
        />
      </view>
    </view>

    <!-- 知识点 -->
    <view class="form-section">
      <view class="section-header">
        <text class="section-title">知识点标签</text>
        <text class="add-btn" @tap="showKnowledgePicker = true">+ 添加知识点</text>
      </view>
      <view v-if="form.knowledgePoints.length > 0" class="knowledge-tags">
        <view
          v-for="(point, index) in form.knowledgePoints"
          :key="index"
          class="knowledge-tag"
        >
          <text>{{ point.name }}</text>
          <text class="tag-close" @tap="removeKnowledge(index)">✕</text>
        </view>
      </view>
      <view v-else class="empty-tip">
        <text>暂未添加知识点</text>
      </view>
    </view>

    <!-- 提交按钮 -->
    <view class="submit-section">
      <button class="submit-btn" :loading="submitting" @tap="handleSubmit">
        {{ isEdit ? '保存修改' : '提交题目' }}
      </button>
    </view>

    <!-- 分类选择弹窗 -->
    <view v-if="showCategoryPicker" class="picker-mask" @tap="showCategoryPicker = false">
      <view class="picker-popup" @tap.stop>
        <view class="picker-header">
          <text class="picker-cancel" @tap="showCategoryPicker = false">取消</text>
          <text class="picker-title">选择分类</text>
          <text class="picker-confirm" @tap="confirmCategory">确定</text>
        </view>
        <scroll-view class="picker-body" scroll-y>
          <view
            v-for="item in categoryList"
            :key="item.id"
            class="picker-item"
            :class="{ active: tempCategoryId === item.id }"
            @tap="tempCategoryId = item.id"
          >
            <text>{{ item.name }}</text>
          </view>
        </scroll-view>
      </view>
    </view>

    <!-- 知识点选择弹窗 -->
    <view v-if="showKnowledgePicker" class="picker-mask" @tap="showKnowledgePicker = false">
      <view class="picker-popup" @tap.stop>
        <view class="picker-header">
          <text class="picker-cancel" @tap="showKnowledgePicker = false">取消</text>
          <text class="picker-title">选择知识点</text>
          <text class="picker-confirm" @tap="confirmKnowledge">确定</text>
        </view>
        <scroll-view class="picker-body" scroll-y>
          <view
            v-for="item in knowledgeList"
            :key="item.id"
            class="picker-item picker-item--check"
            :class="{ active: selectedKnowledgeIds.includes(item.id) }"
            @tap="toggleKnowledge(item.id)"
          >
            <text>{{ item.name }}</text>
            <text v-if="selectedKnowledgeIds.includes(item.id)" class="check-icon">✓</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getQuestionDetail, getCategories, getKnowledgePoints } from '@/api/question'

const isEdit = ref(false)
const questionId = ref(0)
const submitting = ref(false)
const showCategoryPicker = ref(false)
const showKnowledgePicker = ref(false)
const tempCategoryId = ref(0)
const selectedKnowledgeIds = ref<number[]>([])

const categoryList = ref<any[]>([])
const knowledgeList = ref<any[]>([])

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']

const typeOptions = [
  { label: '单选题', value: 1 },
  { label: '多选题', value: 2 },
  { label: '判断题', value: 3 },
  { label: '填空题', value: 4 },
  { label: '简答题', value: 5 }
]

const difficultyOptions = [
  { label: '简单', value: 1 },
  { label: '中等', value: 2 },
  { label: '困难', value: 3 }
]

const form = reactive({
  type: 1,
  difficulty: 1,
  categoryId: 0,
  categoryName: '',
  title: '',
  content: '',
  options: [
    { content: '' },
    { content: '' },
    { content: '' },
    { content: '' }
  ] as { content: string }[],
  correctOptions: [] as number[],
  answer: '',
  analysis: '',
  knowledgePoints: [] as { id: number; name: string }[]
})

const showOptions = computed(() => {
  return form.type >= 1 && form.type <= 3
})

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const id = currentPage?.options?.id
  if (id) {
    isEdit.value = true
    questionId.value = Number(id)
    loadDetail(Number(id))
  }
  loadCategories()
  loadKnowledgePoints()
})

const loadDetail = async (id: number) => {
  try {
    uni.showLoading({ title: '加载中...' })
    const res = await getQuestionDetail(id)
    const data = res.data
    form.type = data.type || 1
    form.difficulty = data.difficulty || 1
    form.categoryId = data.categoryId || 0
    form.categoryName = data.category?.name || ''
    form.title = data.title || ''
    form.content = data.content || ''
    form.options = data.options?.length ? data.options.map((o: any) => ({ content: o.content })) : [{ content: '' }, { content: '' }, { content: '' }, { content: '' }]
    form.correctOptions = data.correctOptions || []
    form.answer = data.answer || ''
    form.analysis = data.analysis || ''
    form.knowledgePoints = data.knowledgePoints || []
    uni.hideLoading()
  } catch (e) {
    uni.hideLoading()
    console.error('加载题目详情失败', e)
  }
}

const loadCategories = async () => {
  try {
    const res = await getCategories()
    categoryList.value = res.data || []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

const loadKnowledgePoints = async () => {
  try {
    const res = await getKnowledgePoints()
    knowledgeList.value = res.data || []
  } catch (e) {
    console.error('加载知识点失败', e)
  }
}

const addOption = () => {
  if (form.options.length >= 8) {
    uni.showToast({ title: '最多8个选项', icon: 'none' })
    return
  }
  form.options.push({ content: '' })
}

const removeOption = (index: number) => {
  form.options.splice(index, 1)
  form.correctOptions = form.correctOptions
    .filter(i => i !== index)
    .map(i => i > index ? i - 1 : i)
}

const toggleCorrect = (index: number) => {
  if (form.type === 1 || form.type === 3) {
    form.correctOptions = [index]
  } else {
    const idx = form.correctOptions.indexOf(index)
    if (idx > -1) {
      form.correctOptions.splice(idx, 1)
    } else {
      form.correctOptions.push(index)
    }
  }
}

const confirmCategory = () => {
  const selected = categoryList.value.find(c => c.id === tempCategoryId.value)
  if (selected) {
    form.categoryId = selected.id
    form.categoryName = selected.name
  }
  showCategoryPicker.value = false
}

const toggleKnowledge = (id: number) => {
  const idx = selectedKnowledgeIds.value.indexOf(id)
  if (idx > -1) {
    selectedKnowledgeIds.value.splice(idx, 1)
  } else {
    selectedKnowledgeIds.value.push(id)
  }
}

const removeKnowledge = (index: number) => {
  const point = form.knowledgePoints[index]
  selectedKnowledgeIds.value = selectedKnowledgeIds.value.filter(id => id !== point.id)
  form.knowledgePoints.splice(index, 1)
}

const confirmKnowledge = () => {
  form.knowledgePoints = knowledgeList.value
    .filter(item => selectedKnowledgeIds.value.includes(item.id))
    .map(item => ({ id: item.id, name: item.name }))
  showKnowledgePicker.value = false
}

const handleSubmit = async () => {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入题目标题', icon: 'none' })
    return
  }
  if (!form.content.trim()) {
    uni.showToast({ title: '请输入题目内容', icon: 'none' })
    return
  }
  if (showOptions.value) {
    const hasEmpty = form.options.some(o => !o.content.trim())
    if (hasEmpty) {
      uni.showToast({ title: '请填写所有选项内容', icon: 'none' })
      return
    }
    if (form.correctOptions.length === 0) {
      uni.showToast({ title: '请选择正确答案', icon: 'none' })
      return
    }
  }

  submitting.value = true
  try {
    // const data = {
    //   type: form.type,
    //   difficulty: form.difficulty,
    //   categoryId: form.categoryId,
    //   title: form.title,
    //   content: form.content,
    //   options: showOptions.value ? form.options : undefined,
    //   correctOptions: showOptions.value ? form.correctOptions : undefined,
    //   answer: form.answer,
    //   analysis: form.analysis,
    //   knowledgePointIds: form.knowledgePoints.map(p => p.id)
    // }
    // if (isEdit.value) {
    //   await updateQuestion(questionId.value, data)
    // } else {
    //   await createQuestion(data)
    // }
    uni.showToast({ title: isEdit.value ? '修改成功' : '提交成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 500)
  } catch (e) {
    console.error('提交失败', e)
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.edit-container {
  min-height: 100vh;
  background-color: #F5F7FA;
  padding-bottom: 40rpx;
}

.form-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .section-title {
    display: block;
    font-size: 30rpx;
    font-weight: bold;
    color: #303133;
    margin-bottom: 24rpx;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24rpx;

    .section-title {
      margin-bottom: 0;
    }
  }

  .add-btn {
    font-size: 26rpx;
    color: #4A90D9;
  }
}

.form-item {
  margin-bottom: 24rpx;

  &:last-child {
    margin-bottom: 0;
  }

  .form-label {
    display: block;
    font-size: 28rpx;
    color: #606266;
    margin-bottom: 16rpx;
  }
}

.form-input {
  width: 100%;
  height: 88rpx;
  padding: 0 24rpx;
  background-color: #F5F7FA;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #303133;
  box-sizing: border-box;
}

.form-textarea {
  width: 100%;
  min-height: 200rpx;
  padding: 20rpx 24rpx;
  background-color: #F5F7FA;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #303133;
  box-sizing: border-box;
}

.type-selector,
.difficulty-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.type-option,
.difficulty-option {
  padding: 16rpx 28rpx;
  background-color: #F5F7FA;
  border-radius: 30rpx;
  font-size: 26rpx;
  color: #606266;

  &.active {
    background-color: #4A90D9;
    color: #ffffff;
  }
}

.picker-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
  background-color: #F5F7FA;
  border-radius: 12rpx;

  .placeholder {
    color: #909399;
  }

  .picker-arrow {
    font-size: 28rpx;
    color: #909399;
  }
}

.option-item {
  margin-bottom: 20rpx;
  padding: 20rpx;
  background-color: #FAFAFA;
  border-radius: 12rpx;

  .option-header {
    display: flex;
    align-items: center;
    margin-bottom: 16rpx;
  }

  .option-label {
    flex-shrink: 0;
    width: 48rpx;
    height: 48rpx;
    border-radius: 50%;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 26rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 16rpx;
  }

  .correct-toggle {
    flex-shrink: 0;
    padding: 8rpx 20rpx;
    border: 1rpx solid #DCDFE6;
    border-radius: 30rpx;
    font-size: 24rpx;
    color: #909399;
    margin-right: 16rpx;

    &.active {
      border-color: #67C23A;
      background-color: #F6FFED;
      color: #67C23A;
    }
  }

  .delete-btn {
    flex-shrink: 0;
    font-size: 24rpx;
    color: #F56C6C;
    margin-left: auto;
  }
}

.knowledge-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.knowledge-tag {
  display: flex;
  align-items: center;
  padding: 10rpx 24rpx;
  background-color: #E6F7FF;
  border-radius: 30rpx;
  font-size: 24rpx;
  color: #1890FF;

  .tag-close {
    margin-left: 12rpx;
    font-size: 20rpx;
    color: #1890FF;
  }
}

.empty-tip {
  text-align: center;
  padding: 40rpx 0;
  font-size: 26rpx;
  color: #909399;
}

.submit-section {
  padding: 30rpx;

  .submit-btn {
    width: 100%;
    height: 100rpx;
    line-height: 100rpx;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 32rpx;
    border-radius: 16rpx;
    border: none;

    &:active {
      opacity: 0.8;
    }
  }
}

.picker-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  align-items: flex-end;
}

.picker-popup {
  width: 100%;
  max-height: 70vh;
  background-color: #ffffff;
  border-radius: 24rpx 24rpx 0 0;
  display: flex;
  flex-direction: column;

  .picker-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 30rpx;
    border-bottom: 1rpx solid #EBEEF5;
  }

  .picker-cancel {
    font-size: 28rpx;
    color: #909399;
  }

  .picker-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #303133;
  }

  .picker-confirm {
    font-size: 28rpx;
    color: #4A90D9;
  }

  .picker-body {
    flex: 1;
    max-height: 600rpx;
    padding: 20rpx 0;
  }
}

.picker-item {
  padding: 24rpx 30rpx;
  font-size: 28rpx;
  color: #303133;
  display: flex;
  align-items: center;
  justify-content: space-between;

  &.active {
    color: #4A90D9;
    background-color: #F5F7FA;
  }

  &--check {
    padding: 20rpx 30rpx;
  }

  .check-icon {
    font-size: 28rpx;
    color: #4A90D9;
  }
}
</style>

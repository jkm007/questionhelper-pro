<template>
  <view class="page">
    <!-- Export Format Selection -->
    <view class="section">
      <text class="section-title">导出格式</text>
      <view class="format-options">
        <view
          :class="['format-item', selectedFormat === 'pdf' ? 'selected' : '']"
          @tap="selectedFormat = 'pdf'"
        >
          <view class="format-icon pdf-icon">
            <text class="format-icon-text">PDF</text>
          </view>
          <text class="format-name">PDF 文档</text>
          <text class="format-desc">适合打印复习</text>
          <view v-if="selectedFormat === 'pdf'" class="format-check">
            <text class="check-mark">&#xe623;</text>
          </view>
        </view>
        <view
          :class="['format-item', selectedFormat === 'excel' ? 'selected' : '']"
          @tap="selectedFormat = 'excel'"
        >
          <view class="format-icon excel-icon">
            <text class="format-icon-text">XLS</text>
          </view>
          <text class="format-name">Excel 表格</text>
          <text class="format-desc">适合数据分析</text>
          <view v-if="selectedFormat === 'excel'" class="format-check">
            <text class="check-mark">&#xe623;</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Filter Options -->
    <view class="section">
      <text class="section-title">筛选条件</text>

      <!-- Mastery Level -->
      <view class="filter-group">
        <text class="filter-label">掌握程度</text>
        <view class="filter-chips">
          <view
            v-for="(label, key) in masteryLabels"
            :key="key"
            :class="['chip', filters.masteryLevel.includes(key) ? 'chip-active' : '']"
            @tap="toggleFilter('masteryLevel', key)"
          >
            <text class="chip-text">{{ label }}</text>
          </view>
        </view>
      </view>

      <!-- Category -->
      <view class="filter-group">
        <text class="filter-label">学科分类</text>
        <view class="filter-chips">
          <view
            v-for="cat in categoryOptions"
            :key="cat"
            :class="['chip', filters.category.includes(cat) ? 'chip-active' : '']"
            @tap="toggleFilter('category', cat)"
          >
            <text class="chip-text">{{ cat }}</text>
          </view>
        </view>
      </view>

      <!-- Date Range -->
      <view class="filter-group">
        <text class="filter-label">时间范围</text>
        <view class="date-row">
          <view class="date-input" @tap="pickDate('start')">
            <text :class="['date-text', filters.startDate ? '' : 'placeholder']">
              {{ filters.startDate || '开始日期' }}
            </text>
          </view>
          <text class="date-separator">至</text>
          <view class="date-input" @tap="pickDate('end')">
            <text :class="['date-text', filters.endDate ? '' : 'placeholder']">
              {{ filters.endDate || '结束日期' }}
            </text>
          </view>
        </view>
        <view class="date-shortcuts">
          <view class="shortcut" @tap="setDateRange('week')">
            <text class="shortcut-text">最近一周</text>
          </view>
          <view class="shortcut" @tap="setDateRange('month')">
            <text class="shortcut-text">最近一月</text>
          </view>
          <view class="shortcut" @tap="setDateRange('quarter')">
            <text class="shortcut-text">最近三月</text>
          </view>
          <view class="shortcut" @tap="setDateRange('all')">
            <text class="shortcut-text">全部</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Preview -->
    <view class="section">
      <text class="section-title">导出预览</text>
      <view v-if="previewLoading" class="preview-loading">
        <text class="loading-text">正在生成预览...</text>
      </view>
      <view v-else class="preview-content">
        <view class="preview-header">
          <text class="preview-count">共 {{ previewCount }} 道错题</text>
          <text class="preview-hint">{{ previewMasterySummary }}</text>
        </view>
        <view class="preview-table">
          <view class="table-row table-header">
            <text class="table-cell cell-title">题目</text>
            <text class="table-cell cell-type">类型</text>
            <text class="table-cell cell-mastery">掌握度</text>
            <text class="table-cell cell-count">错误次数</text>
          </view>
          <view v-for="(item, idx) in previewList" :key="idx" class="table-row">
            <text class="table-cell cell-title">{{ item.title }}</text>
            <text class="table-cell cell-type">{{ typeLabels[item.questionType] }}</text>
            <text class="table-cell cell-mastery">{{ item.masteryLevel }}%</text>
            <text class="table-cell cell-count">{{ item.wrongCount }}</text>
          </view>
        </view>
        <view v-if="previewCount > 5" class="preview-more">
          <text class="more-text">还有 {{ previewCount - 5 }} 道题目...</text>
        </view>
        <view v-if="previewCount === 0" class="preview-empty">
          <text class="empty-text">暂无符合条件的错题</text>
        </view>
      </view>
    </view>

    <!-- Export Button -->
    <view class="bottom-bar">
      <view
        :class="['export-btn', canExport ? '' : 'export-btn-disabled']"
        @tap="onExport"
      >
        <text class="export-btn-text">
          {{ exporting ? '正在导出...' : '开始导出' }}
        </text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

interface PreviewItem {
  id: string
  title: string
  questionType: string
  masteryLevel: number
  wrongCount: number
}

const selectedFormat = ref<'pdf' | 'excel'>('pdf')
const previewLoading = ref(false)
const previewCount = ref(0)
const previewList = ref<PreviewItem[]>([])
const exporting = ref(false)

const filters = ref({
  masteryLevel: [] as string[],
  category: [] as string[],
  startDate: '',
  endDate: ''
})

const masteryLabels: Record<string, string> = {
  low: '未掌握',
  medium: '部分掌握',
  high: '已掌握'
}

const typeLabels: Record<string, string> = {
  single: '单选题',
  multiple: '多选题',
  judge: '判断题',
  fill: '填空题',
  essay: '简答题'
}

const categoryOptions = ['数学', '英语', '语文', '物理', '化学', '生物', '历史', '地理', '政治']

const canExport = computed(() => {
  return previewCount.value > 0 && !exporting.value
})

const previewMasterySummary = computed(() => {
  if (previewCount.value === 0) return ''
  const low = previewList.value.filter((i) => i.masteryLevel < 30).length
  const med = previewList.value.filter((i) => i.masteryLevel >= 30 && i.masteryLevel < 70).length
  const high = previewList.value.filter((i) => i.masteryLevel >= 70).length
  const parts = []
  if (low > 0) parts.push(`${low} 题未掌握`)
  if (med > 0) parts.push(`${med} 题部分掌握`)
  if (high > 0) parts.push(`${high} 题已掌握`)
  return parts.join('，')
})

function toggleFilter(field: 'masteryLevel' | 'category', value: string) {
  const arr = filters.value[field]
  const idx = arr.indexOf(value)
  if (idx >= 0) {
    arr.splice(idx, 1)
  } else {
    arr.push(value)
  }
}

function pickDate(type: 'start' | 'end') {
  const currentDate = type === 'start' ? filters.value.startDate : filters.value.endDate
  uni.showModal({
    title: type === 'start' ? '选择开始日期' : '选择结束日期',
    content: '请输入日期 (YYYY-MM-DD)',
    editable: true,
    placeholderText: currentDate || 'YYYY-MM-DD',
    success: (res) => {
      if (res.confirm && res.content) {
        const dateStr = res.content.trim()
        if (/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) {
          if (type === 'start') {
            filters.value.startDate = dateStr
          } else {
            filters.value.endDate = dateStr
          }
        } else {
          uni.showToast({ title: '日期格式不正确', icon: 'none' })
        }
      }
    }
  })
}

function setDateRange(range: string) {
  const now = new Date()
  const end = formatDate(now)
  let start = ''
  if (range === 'all') {
    filters.value.startDate = ''
    filters.value.endDate = ''
    return
  }
  if (range === 'week') {
    start = formatDate(new Date(now.getTime() - 7 * 86400000))
  } else if (range === 'month') {
    start = formatDate(new Date(now.getTime() - 30 * 86400000))
  } else if (range === 'quarter') {
    start = formatDate(new Date(now.getTime() - 90 * 86400000))
  }
  filters.value.startDate = start
  filters.value.endDate = end
}

function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

async function loadPreview() {
  previewLoading.value = true
  try {
    const params: Record<string, any> = {
      pageSize: 5
    }
    if (filters.value.masteryLevel.length > 0) {
      params.masteryLevel = filters.value.masteryLevel.join(',')
    }
    if (filters.value.category.length > 0) {
      params.category = filters.value.category.join(',')
    }
    if (filters.value.startDate) {
      params.startDate = filters.value.startDate
    }
    if (filters.value.endDate) {
      params.endDate = filters.value.endDate
    }

    // Simulate API call for preview
    await new Promise((r) => setTimeout(r, 500))
    previewList.value = [
      { id: '1', title: '二次函数图像的平移变换', questionType: 'single', masteryLevel: 20, wrongCount: 3 },
      { id: '2', title: '三角函数的基本性质', questionType: 'multiple', masteryLevel: 45, wrongCount: 2 },
      { id: '3', title: '概率与统计综合应用', questionType: 'essay', masteryLevel: 15, wrongCount: 5 },
      { id: '4', title: '立体几何-异面直线夹角', questionType: 'fill', masteryLevel: 60, wrongCount: 1 },
      { id: '5', title: '导数的几何意义', questionType: 'judge', masteryLevel: 35, wrongCount: 4 }
    ]
    previewCount.value = previewList.value.length > 0 ? previewList.value.length + 3 : 0
  } catch (e) {
    uni.showToast({ title: '加载预览失败', icon: 'none' })
    previewCount.value = 0
    previewList.value = []
  } finally {
    previewLoading.value = false
  }
}

async function onExport() {
  if (!canExport.value) return
  exporting.value = true
  try {
    const exportParams: Record<string, any> = {
      format: selectedFormat.value
    }
    if (filters.value.masteryLevel.length > 0) {
      exportParams.masteryLevel = filters.value.masteryLevel.join(',')
    }
    if (filters.value.category.length > 0) {
      exportParams.category = filters.value.category.join(',')
    }
    if (filters.value.startDate) {
      exportParams.startDate = filters.value.startDate
    }
    if (filters.value.endDate) {
      exportParams.endDate = filters.value.endDate
    }

    // Simulate export process
    await new Promise((r) => setTimeout(r, 2000))

    uni.showModal({
      title: '导出成功',
      content: `已成功导出 ${previewCount.value} 道错题，是否立即查看？`,
      confirmText: '查看',
      cancelText: '稍后',
      success: (res) => {
        if (res.confirm) {
          const formatName = selectedFormat.value === 'pdf' ? 'PDF' : 'Excel'
          uni.showToast({ title: `已保存至本地 (${formatName})`, icon: 'success', duration: 2000 })
        }
      }
    })
  } catch (e) {
    uni.showToast({ title: '导出失败，请重试', icon: 'none' })
  } finally {
    exporting.value = false
  }
}

watch(filters, () => {
  loadPreview()
}, { deep: true })

onMounted(() => {
  loadPreview()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding: 20rpx 24rpx;
  padding-bottom: 160rpx;
}

.section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
  display: block;
}

/* Format Options */
.format-options {
  display: flex;
  gap: 20rpx;
}

.format-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 28rpx 16rpx;
  border: 2rpx solid #e8e8e8;
  border-radius: 16rpx;
  background-color: #fafafa;
  position: relative;
}

.format-item.selected {
  border-color: #1677ff;
  background-color: #e8f3ff;
}

.format-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16rpx;
}

.pdf-icon {
  background-color: #ff4d4f;
}

.excel-icon {
  background-color: #52c41a;
}

.format-icon-text {
  font-size: 24rpx;
  font-weight: 700;
  color: #fff;
}

.format-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.format-desc {
  font-size: 22rpx;
  color: #999;
}

.format-check {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  width: 36rpx;
  height: 36rpx;
  background-color: #1677ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.check-mark {
  font-size: 20rpx;
  color: #fff;
}

/* Filter Options */
.filter-group {
  margin-bottom: 28rpx;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-label {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 16rpx;
  display: block;
}

.filter-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.chip {
  padding: 12rpx 24rpx;
  border: 2rpx solid #e0e0e0;
  border-radius: 32rpx;
  background-color: #fafafa;
}

.chip.chip-active {
  border-color: #1677ff;
  background-color: #e8f3ff;
}

.chip-text {
  font-size: 26rpx;
  color: #666;
}

.chip.chip-active .chip-text {
  color: #1677ff;
}

/* Date Range */
.date-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.date-input {
  flex: 1;
  padding: 16rpx 20rpx;
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  background-color: #fafafa;
}

.date-text {
  font-size: 26rpx;
  color: #333;
}

.date-text.placeholder {
  color: #ccc;
}

.date-separator {
  font-size: 26rpx;
  color: #999;
  flex-shrink: 0;
}

.date-shortcuts {
  display: flex;
  gap: 16rpx;
}

.shortcut {
  padding: 10rpx 20rpx;
  background-color: #f0f1f5;
  border-radius: 24rpx;
}

.shortcut-text {
  font-size: 24rpx;
  color: #666;
}

/* Preview */
.preview-loading {
  padding: 40rpx 0;
  text-align: center;
}

.loading-text {
  font-size: 26rpx;
  color: #999;
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.preview-count {
  font-size: 28rpx;
  font-weight: 600;
  color: #1677ff;
}

.preview-hint {
  font-size: 22rpx;
  color: #999;
}

.preview-table {
  border: 2rpx solid #f0f0f0;
  border-radius: 12rpx;
  overflow: hidden;
}

.table-row {
  display: flex;
  align-items: center;
  border-bottom: 1rpx solid #f0f0f0;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row.table-header {
  background-color: #fafafa;
}

.table-cell {
  padding: 16rpx 12rpx;
  font-size: 24rpx;
  color: #666;
  text-align: center;
}

.table-header .table-cell {
  font-weight: 600;
  color: #333;
}

.cell-title {
  flex: 3;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cell-type {
  flex: 1.5;
}

.cell-mastery {
  flex: 1.5;
}

.cell-count {
  flex: 1.5;
}

.preview-more {
  padding: 16rpx 0;
  text-align: center;
}

.more-text {
  font-size: 24rpx;
  color: #999;
}

.preview-empty {
  padding: 40rpx 0;
  text-align: center;
}

.empty-text {
  font-size: 26rpx;
  color: #ccc;
}

/* Bottom Bar */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.export-btn {
  width: 100%;
  padding: 24rpx 0;
  background-color: #1677ff;
  border-radius: 44rpx;
  text-align: center;
}

.export-btn-disabled {
  background-color: #c0c4cc;
}

.export-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}
</style>

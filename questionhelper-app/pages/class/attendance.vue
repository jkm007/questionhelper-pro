<template>
  <view class="page">
    <!-- Tab Bar: Teacher / Student views -->
    <view class="tab-bar">
      <view
        :class="['tab-item', activeTab === 'sessions' ? 'active' : '']"
        @tap="activeTab = 'sessions'"
      >
        <text class="tab-text">考勤记录</text>
      </view>
      <view
        :class="['tab-item', activeTab === 'calendar' ? 'active' : '']"
        @tap="activeTab = 'calendar'"
      >
        <text class="tab-text">日历视图</text>
      </view>
      <view
        v-if="isTeacher"
        :class="['tab-item', activeTab === 'create' ? 'active' : '']"
        @tap="activeTab = 'create'"
      >
        <text class="tab-text">发起考勤</text>
      </view>
    </view>

    <!-- Sessions Tab -->
    <scroll-view
      v-if="activeTab === 'sessions'"
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <!-- Summary Stats -->
      <view class="summary-card">
        <view class="summary-item">
          <text class="summary-num">{{ summary.total }}</text>
          <text class="summary-label">总次数</text>
        </view>
        <view class="summary-divider"></view>
        <view class="summary-item">
          <text class="summary-num summary-success">{{ summary.checkedIn }}</text>
          <text class="summary-label">已签到</text>
        </view>
        <view class="summary-divider"></view>
        <view class="summary-item">
          <text class="summary-num summary-warn">{{ summary.late }}</text>
          <text class="summary-label">迟到</text>
        </view>
        <view class="summary-divider"></view>
        <view class="summary-item">
          <text class="summary-num summary-error">{{ summary.absent }}</text>
          <text class="summary-label">未签到</text>
        </view>
      </view>

      <!-- Session List -->
      <view
        v-for="session in sessions"
        :key="session.id"
        class="session-card"
        @tap="onTapSession(session)"
      >
        <view class="session-header">
          <text class="session-title">{{ session.title }}</text>
          <view :class="['session-status', 'status-' + session.status]">
            <text class="session-status-text">{{ statusLabels[session.status] }}</text>
          </view>
        </view>
        <view class="session-body">
          <view class="session-info-row">
            <text class="session-info-label">发起时间</text>
            <text class="session-info-value">{{ session.createdAt }}</text>
          </view>
          <view class="session-info-row">
            <text class="session-info-label">截止时间</text>
            <text class="session-info-value">{{ session.deadline }}</text>
          </view>
          <view v-if="isTeacher" class="session-info-row">
            <text class="session-info-label">签到率</text>
            <text class="session-info-value rate">{{ session.checkInRate }}%</text>
          </view>
          <view v-else class="session-info-row">
            <text class="session-info-label">我的状态</text>
            <view :class="['my-status-badge', 'my-' + session.myStatus]">
              <text class="my-status-text">{{ attendanceLabels[session.myStatus] }}</text>
            </view>
          </view>
        </view>
        <view v-if="isTeacher" class="session-footer">
          <text class="session-stats">
            已签 {{ session.checkedInCount }} / 迟到 {{ session.lateCount }} / 未签 {{ session.absentCount }}
          </text>
          <view class="session-action-btn" @tap.stop="onExportSession(session)">
            <text class="session-action-text">导出</text>
          </view>
        </view>
        <view v-else-if="session.status === 'active' && session.myStatus === 'absent'" class="session-footer">
          <view class="checkin-btn" @tap.stop="onCheckIn(session)">
            <text class="checkin-btn-text">签到</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && sessions.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/default.png" mode="aspectFit" />
        <text class="empty-text">暂无考勤记录</text>
      </view>

      <!-- Load More -->
      <view v-if="sessions.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <!-- Calendar Tab -->
    <view v-if="activeTab === 'calendar'" class="calendar-container">
      <!-- Month Navigator -->
      <view class="month-nav">
        <view class="month-arrow" @tap="changeMonth(-1)">
          <text class="arrow-text">&#xe61e;</text>
        </view>
        <text class="month-title">{{ currentYear }}年{{ currentMonth }}月</text>
        <view class="month-arrow" @tap="changeMonth(1)">
          <text class="arrow-text">&#xe61e;</text>
        </view>
      </view>

      <!-- Week Header -->
      <view class="week-header">
        <text v-for="day in weekDays" :key="day" class="week-day">{{ day }}</text>
      </view>

      <!-- Calendar Grid -->
      <view class="calendar-grid">
        <view
          v-for="(cell, index) in calendarCells"
          :key="index"
          :class="['calendar-cell', cell.isEmpty ? 'empty' : '', cell.isToday ? 'today' : '', 'att-' + cell.status]"
          @tap="cell.day && onCalendarDay(cell)"
        >
          <text class="cell-day">{{ cell.day || '' }}</text>
          <view v-if="cell.status && cell.status !== 'none'" class="cell-dot"></view>
        </view>
      </view>

      <!-- Legend -->
      <view class="legend">
        <view class="legend-item">
          <view class="legend-dot legend-checked"></view>
          <text class="legend-text">已签到</text>
        </view>
        <view class="legend-item">
          <view class="legend-dot legend-late"></view>
          <text class="legend-text">迟到</text>
        </view>
        <view class="legend-item">
          <view class="legend-dot legend-early"></view>
          <text class="legend-text">早退</text>
        </view>
        <view class="legend-item">
          <view class="legend-dot legend-absent"></view>
          <text class="legend-text">未签到</text>
        </view>
      </view>

      <!-- Selected Day Detail -->
      <view v-if="selectedDayRecords.length > 0" class="day-detail">
        <text class="day-detail-title">{{ selectedDayText }} 考勤详情</text>
        <view
          v-for="record in selectedDayRecords"
          :key="record.id"
          class="day-record-item"
        >
          <text class="day-record-title">{{ record.sessionTitle }}</text>
          <view :class="['day-record-badge', 'my-' + record.status]">
            <text class="day-record-text">{{ attendanceLabels[record.status] }}</text>
          </view>
          <text class="day-record-time">{{ record.time }}</text>
        </view>
      </view>
    </view>

    <!-- Create Tab (Teacher Only) -->
    <view v-if="activeTab === 'create' && isTeacher" class="create-container">
      <view class="form-section">
        <view class="form-item">
          <text class="form-label">考勤标题</text>
          <input
            class="form-input"
            v-model="createForm.title"
            placeholder="请输入考勤标题"
            maxlength="50"
          />
        </view>
        <view class="form-item">
          <text class="form-label">截止时间</text>
          <picker mode="time" :value="createForm.deadlineTime" @change="onTimeChange">
            <view class="picker-display">
              <text :class="['picker-text', createForm.deadlineTime ? '' : 'placeholder']">
                {{ createForm.deadlineTime || '请选择截止时间' }}
              </text>
            </view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">截止时长(分钟)</text>
          <input
            class="form-input"
            v-model="createForm.duration"
            type="number"
            placeholder="签到持续时长，默认30分钟"
          />
        </view>
        <view class="form-item">
          <text class="form-label">备注说明</text>
          <textarea
            class="form-textarea"
            v-model="createForm.remark"
            placeholder="可选，填写备注信息"
            maxlength="200"
          />
        </view>
      </view>
      <view class="create-btn" @tap="onCreateSession">
        <text class="create-btn-text">发起考勤</text>
      </view>
    </view>

    <!-- Session Detail Modal (Teacher: record list) -->
    <view v-if="showDetail" class="modal-mask" @tap="showDetail = false"></view>
    <view :class="['detail-modal', showDetail ? 'open' : '']">
      <view class="modal-header">
        <text class="modal-title">{{ detailSession?.title }}</text>
        <text class="modal-close" @tap="showDetail = false">&#xe621;</text>
      </view>
      <scroll-view class="modal-body" scroll-y>
        <view class="detail-stats">
          <view class="detail-stat">
            <text class="detail-stat-num">{{ detailRecords.length }}</text>
            <text class="detail-stat-label">总人数</text>
          </view>
          <view class="detail-stat">
            <text class="detail-stat-num detail-success">{{ detailCheckedCount }}</text>
            <text class="detail-stat-label">已签到</text>
          </view>
          <view class="detail-stat">
            <text class="detail-stat-num detail-error">{{ detailAbsentCount }}</text>
            <text class="detail-stat-label">未签到</text>
          </view>
        </view>
        <view
          v-for="record in detailRecords"
          :key="record.id"
          class="detail-record"
        >
          <image class="detail-avatar" :src="record.avatar" mode="aspectFill" />
          <view class="detail-info">
            <text class="detail-name">{{ record.userName }}</text>
            <text class="detail-time">{{ record.checkInTime || '未签到' }}</text>
          </view>
          <view :class="['detail-badge', 'my-' + record.status]">
            <text class="detail-badge-text">{{ attendanceLabels[record.status] }}</text>
          </view>
        </view>
        <view v-if="detailRecords.length === 0" class="modal-empty">
          <text class="modal-empty-text">暂无记录</text>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import {
  getAttendanceSessions,
  getAttendanceRecords,
  createAttendanceSession,
  studentCheckIn,
  exportAttendance,
  getAttendanceCalendar
} from '@/api/class'

// ---------- Types ----------

interface AttendanceSession {
  id: string
  title: string
  status: 'active' | 'closed' | 'expired'
  createdAt: string
  deadline: string
  checkInRate?: number
  checkedInCount?: number
  lateCount?: number
  absentCount?: number
  myStatus?: AttendanceStatus
}

type AttendanceStatus = 'checked' | 'absent' | 'late' | 'early'

interface AttendanceRecord {
  id: string
  userId: string
  userName: string
  avatar: string
  status: AttendanceStatus
  checkInTime: string
  sessionTitle?: string
  time?: string
}

interface CalendarCell {
  day: number | null
  isEmpty: boolean
  isToday: boolean
  status: string
  date?: string
}

interface CalendarRecord {
  date: string
  status: AttendanceStatus
  sessionTitle: string
  id: string
  time: string
}

// ---------- Constants ----------

const statusLabels: Record<string, string> = {
  active: '进行中',
  closed: '已结束',
  expired: '已过期'
}

const attendanceLabels: Record<string, string> = {
  checked: '已签到',
  absent: '未签到',
  late: '迟到',
  early: '早退'
}

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

// ---------- State ----------

const classId = ref('')
const isTeacher = ref(false)
const activeTab = ref<'sessions' | 'calendar' | 'create'>('sessions')

// Sessions tab
const sessions = ref<AttendanceSession[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const summary = reactive({ total: 0, checkedIn: 0, late: 0, absent: 0 })

// Calendar tab
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth() + 1)
const calendarData = ref<CalendarRecord[]>([])
const selectedDayRecords = ref<CalendarRecord[]>([])
const selectedDayText = ref('')

// Create tab
const createForm = reactive({
  title: '',
  deadlineTime: '',
  duration: '30',
  remark: ''
})

// Detail modal
const showDetail = ref(false)
const detailSession = ref<AttendanceSession | null>(null)
const detailRecords = ref<AttendanceRecord[]>([])

const detailCheckedCount = computed(() =>
  detailRecords.value.filter((r) => r.status === 'checked' || r.status === 'late').length
)
const detailAbsentCount = computed(() =>
  detailRecords.value.filter((r) => r.status === 'absent').length
)

// ---------- Calendar computed ----------

const calendarCells = computed<CalendarCell[]>(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month - 1, 1).getDay()
  const daysInMonth = new Date(year, month, 0).getDate()
  const today = new Date()
  const isCurrentMonth = today.getFullYear() === year && today.getMonth() + 1 === month

  const cells: CalendarCell[] = []

  // Leading empty cells
  for (let i = 0; i < firstDay; i++) {
    cells.push({ day: null, isEmpty: true, isToday: false, status: '' })
  }

  // Day cells
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${year}-${String(month).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const record = calendarData.value.find((r) => r.date === dateStr)
    cells.push({
      day: d,
      isEmpty: false,
      isToday: isCurrentMonth && today.getDate() === d,
      status: record ? record.status : 'none',
      date: dateStr
    })
  }

  return cells
})

// ---------- Data loading ----------

async function loadSessions(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    sessions.value = []
  }
  loading.value = true
  try {
    const res = await getAttendanceSessions(classId.value, {
      page: page.value,
      pageSize
    })
    const data = res.data || {}
    const newList: AttendanceSession[] = data.list || []
    if (reset) {
      sessions.value = newList
    } else {
      sessions.value = [...sessions.value, ...newList]
    }
    if (data.summary) {
      summary.total = data.summary.total || 0
      summary.checkedIn = data.summary.checkedIn || 0
      summary.late = data.summary.late || 0
      summary.absent = data.summary.absent || 0
    }
    noMore.value = newList.length < pageSize
    page.value++
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function loadCalendar() {
  try {
    const res = await getAttendanceCalendar(classId.value, {
      year: currentYear.value,
      month: currentMonth.value
    })
    calendarData.value = res.data?.list || []
    selectedDayRecords.value = []
    selectedDayText.value = ''
  } catch (e) {
    uni.showToast({ title: '加载日历数据失败', icon: 'none' })
  }
}

async function loadSessionDetail(session: AttendanceSession) {
  detailSession.value = session
  showDetail.value = true
  try {
    const res = await getAttendanceRecords(classId.value, session.id)
    detailRecords.value = res.data?.list || []
  } catch (e) {
    uni.showToast({ title: '加载详情失败', icon: 'none' })
  }
}

// ---------- Event handlers ----------

async function onRefresh() {
  refreshing.value = true
  await loadSessions(true)
  refreshing.value = false
}

function onLoadMore() {
  if (!noMore.value && !loading.value) {
    loadSessions()
  }
}

function onTapSession(session: AttendanceSession) {
  if (isTeacher.value) {
    loadSessionDetail(session)
  }
}

function changeMonth(delta: number) {
  let m = currentMonth.value + delta
  let y = currentYear.value
  if (m < 1) {
    m = 12
    y--
  } else if (m > 12) {
    m = 1
    y++
  }
  currentYear.value = y
  currentMonth.value = m
  loadCalendar()
}

function onCalendarDay(cell: CalendarCell) {
  if (!cell.date) return
  selectedDayText.value = `${currentMonth.value}月${cell.day}日`
  selectedDayRecords.value = calendarData.value.filter((r) => r.date === cell.date)
}

async function onCheckIn(session: AttendanceSession) {
  uni.showModal({
    title: '确认签到',
    content: `确定签到「${session.title}」吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await studentCheckIn(classId.value, session.id)
          uni.showToast({ title: '签到成功', icon: 'success' })
          loadSessions(true)
        } catch (e) {
          uni.showToast({ title: '签到失败', icon: 'none' })
        }
      }
    }
  })
}

function onTimeChange(e: any) {
  createForm.deadlineTime = e.detail.value
}

async function onCreateSession() {
  if (!createForm.title.trim()) {
    uni.showToast({ title: '请输入考勤标题', icon: 'none' })
    return
  }
  try {
    await createAttendanceSession(classId.value, {
      title: createForm.title.trim(),
      deadlineTime: createForm.deadlineTime || undefined,
      duration: Number(createForm.duration) || 30,
      remark: createForm.remark.trim() || undefined
    })
    uni.showToast({ title: '发起成功', icon: 'success' })
    createForm.title = ''
    createForm.deadlineTime = ''
    createForm.duration = '30'
    createForm.remark = ''
    activeTab.value = 'sessions'
    loadSessions(true)
  } catch (e) {
    uni.showToast({ title: '发起失败', icon: 'none' })
  }
}

async function onExportSession(session: AttendanceSession) {
  try {
    uni.showLoading({ title: '导出中...' })
    await exportAttendance(classId.value, session.id)
    uni.hideLoading()
    uni.showToast({ title: '导出成功', icon: 'success' })
  } catch (e) {
    uni.hideLoading()
    uni.showToast({ title: '导出失败', icon: 'none' })
  }
}

// ---------- Lifecycle ----------

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.classId || ''
  // Role is passed via route param or determined from class detail
  const role = currentPage?.options?.role || ''
  isTeacher.value = role === 'teacher' || role === 'creator' || role === 'admin'

  if (!classId.value) {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
    return
  }
  loadSessions(true)
  loadCalendar()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
}

/* ---- Tab Bar ---- */
.tab-bar {
  display: flex;
  background-color: #fff;
  padding: 0 24rpx;
  border-bottom: 1rpx solid #eee;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  position: relative;
}

.tab-item.active .tab-text {
  color: #1677ff;
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 48rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background-color: #1677ff;
}

.tab-text {
  font-size: 28rpx;
  color: #666;
}

/* ---- Summary Card ---- */
.summary-card {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx 0;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.summary-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.summary-num {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
}

.summary-success {
  color: #52c41a;
}

.summary-warn {
  color: #fa8c16;
}

.summary-error {
  color: #ff4d4f;
}

.summary-label {
  font-size: 24rpx;
  color: #999;
}

.summary-divider {
  width: 1rpx;
  height: 48rpx;
  background-color: #e8e8e8;
}

/* ---- Sessions List ---- */
.list-container {
  height: calc(100vh - 90rpx);
}

.session-card {
  background-color: #fff;
  border-radius: 16rpx;
  margin: 0 24rpx 20rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.session-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.session-title {
  flex: 1;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-right: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-status {
  padding: 6rpx 16rpx;
  border-radius: 16rpx;
  flex-shrink: 0;
}

.status-active {
  background-color: #e8f3ff;
}

.status-active .session-status-text {
  color: #1677ff;
}

.status-closed {
  background-color: #f0f0f0;
}

.status-closed .session-status-text {
  color: #999;
}

.status-expired {
  background-color: #fff1f0;
}

.status-expired .session-status-text {
  color: #ff4d4f;
}

.session-status-text {
  font-size: 22rpx;
}

.session-body {
  margin-bottom: 16rpx;
}

.session-info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6rpx 0;
}

.session-info-label {
  font-size: 26rpx;
  color: #999;
}

.session-info-value {
  font-size: 26rpx;
  color: #333;
}

.session-info-value.rate {
  color: #1677ff;
  font-weight: 600;
}

.my-status-badge {
  padding: 4rpx 14rpx;
  border-radius: 12rpx;
}

.my-checked {
  background-color: #f6ffed;
}

.my-checked .my-status-text,
.my-checked .detail-badge-text,
.my-checked .day-record-text {
  color: #52c41a;
}

.my-late {
  background-color: #fff0e6;
}

.my-late .my-status-text,
.my-late .detail-badge-text,
.my-late .day-record-text {
  color: #fa8c16;
}

.my-early {
  background-color: #f9f0ff;
}

.my-early .my-status-text,
.my-early .detail-badge-text,
.my-early .day-record-text {
  color: #722ed1;
}

.my-absent {
  background-color: #fff1f0;
}

.my-absent .my-status-text,
.my-absent .detail-badge-text,
.my-absent .day-record-text {
  color: #ff4d4f;
}

.my-status-text,
.detail-badge-text,
.day-record-text {
  font-size: 22rpx;
}

.session-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16rpx;
  border-top: 1rpx solid #f5f5f5;
}

.session-stats {
  font-size: 24rpx;
  color: #999;
}

.session-action-btn {
  padding: 10rpx 24rpx;
  border: 2rpx solid #1677ff;
  border-radius: 28rpx;
}

.session-action-text {
  font-size: 24rpx;
  color: #1677ff;
}

.checkin-btn {
  flex: 1;
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 20rpx 0;
  text-align: center;
}

.checkin-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}

/* ---- Calendar ---- */
.calendar-container {
  padding: 20rpx 24rpx;
}

.month-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 0;
}

.month-arrow {
  padding: 12rpx 20rpx;
}

.arrow-text {
  font-size: 28rpx;
  color: #1677ff;
}

.month-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.week-header {
  display: flex;
  padding: 16rpx 0;
}

.week-day {
  flex: 1;
  text-align: center;
  font-size: 24rpx;
  color: #999;
}

.calendar-grid {
  display: flex;
  flex-wrap: wrap;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.calendar-cell {
  width: calc(100% / 7);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16rpx 0;
  position: relative;
}

.calendar-cell.empty {
  visibility: hidden;
}

.cell-day {
  font-size: 28rpx;
  color: #333;
}

.calendar-cell.today .cell-day {
  color: #1677ff;
  font-weight: 700;
}

.cell-dot {
  width: 10rpx;
  height: 10rpx;
  border-radius: 50%;
  margin-top: 6rpx;
}

.att-checked .cell-dot {
  background-color: #52c41a;
}

.att-late .cell-dot {
  background-color: #fa8c16;
}

.att-early .cell-dot {
  background-color: #722ed1;
}

.att-absent .cell-dot {
  background-color: #ff4d4f;
}

.att-none .cell-dot {
  display: none;
}

/* Legend */
.legend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32rpx;
  padding: 24rpx 0;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.legend-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
}

.legend-checked {
  background-color: #52c41a;
}

.legend-late {
  background-color: #fa8c16;
}

.legend-early {
  background-color: #722ed1;
}

.legend-absent {
  background-color: #ff4d4f;
}

.legend-text {
  font-size: 22rpx;
  color: #999;
}

/* Day Detail */
.day-detail {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-top: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.day-detail-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.day-record-item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.day-record-item:last-child {
  border-bottom: none;
}

.day-record-title {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  margin-right: 16rpx;
}

.day-record-badge {
  padding: 4rpx 14rpx;
  border-radius: 12rpx;
  margin-right: 16rpx;
}

.day-record-time {
  font-size: 22rpx;
  color: #999;
  flex-shrink: 0;
}

/* ---- Create Form ---- */
.create-container {
  padding: 20rpx 24rpx;
}

.form-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 12rpx 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.form-item {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.form-input {
  width: 100%;
  font-size: 28rpx;
  color: #333;
  padding: 16rpx 0;
}

.form-textarea {
  width: 100%;
  font-size: 28rpx;
  color: #333;
  padding: 16rpx 0;
  height: 160rpx;
}

.picker-display {
  padding: 16rpx 0;
}

.picker-text {
  font-size: 28rpx;
  color: #333;
}

.picker-text.placeholder {
  color: #ccc;
}

.create-btn {
  margin-top: 40rpx;
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 28rpx 0;
  text-align: center;
}

.create-btn-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
}

/* ---- Empty / Load More ---- */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 0;
}

.empty-img {
  width: 240rpx;
  height: 240rpx;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.load-more {
  padding: 24rpx 0;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #999;
}

/* ---- Detail Modal ---- */
.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

.detail-modal {
  position: fixed;
  left: 0;
  right: 0;
  bottom: -100%;
  background-color: #f5f6fa;
  border-radius: 24rpx 24rpx 0 0;
  z-index: 101;
  max-height: 80vh;
  transition: bottom 0.3s ease;
}

.detail-modal.open {
  bottom: 0;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  background-color: #fff;
  border-radius: 24rpx 24rpx 0 0;
  border-bottom: 1rpx solid #eee;
}

.modal-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.modal-close {
  font-size: 32rpx;
  color: #999;
  padding: 8rpx;
}

.modal-body {
  max-height: calc(80vh - 100rpx);
  padding: 0 32rpx 32rpx;
}

.detail-stats {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background-color: #fff;
  border-radius: 16rpx;
  margin: 20rpx 0;
  padding: 24rpx 0;
}

.detail-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.detail-stat-num {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
}

.detail-success {
  color: #52c41a;
}

.detail-error {
  color: #ff4d4f;
}

.detail-stat-label {
  font-size: 24rpx;
  color: #999;
}

.detail-record {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 20rpx;
  margin-bottom: 12rpx;
}

.detail-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  flex-shrink: 0;
  margin-right: 20rpx;
}

.detail-info {
  flex: 1;
  min-width: 0;
}

.detail-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 4rpx;
}

.detail-time {
  font-size: 22rpx;
  color: #999;
}

.detail-badge {
  padding: 6rpx 16rpx;
  border-radius: 16rpx;
  flex-shrink: 0;
}

.modal-empty {
  padding: 60rpx 0;
  text-align: center;
}

.modal-empty-text {
  font-size: 28rpx;
  color: #999;
}
</style>

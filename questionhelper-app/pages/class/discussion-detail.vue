<template>
  <view class="page">
    <scroll-view class="content-area" scroll-y>
      <!-- Discussion Content -->
      <view class="discuss-header">
        <text class="discuss-title">{{ discussion.title }}</text>
        <view class="discuss-info">
          <text class="info-author">{{ discussion.authorName }}</text>
          <text class="info-time">{{ discussion.createdAt }}</text>
        </view>
        <text class="discuss-body">{{ discussion.content }}</text>
      </view>

      <!-- Replies -->
      <view class="replies-section">
        <text class="replies-title">全部回复 ({{ replies.length }})</text>
        <view v-for="reply in replies" :key="reply.id" class="reply-card">
          <view class="reply-header">
            <text class="reply-author">{{ reply.authorName }}</text>
            <text class="reply-time">{{ reply.createdAt }}</text>
          </view>
          <text class="reply-content">{{ reply.content }}</text>
          <view class="reply-actions">
            <text class="action-btn" @tap="onLike(reply)">
              {{ reply.liked ? '已赞' : '赞' }} {{ reply.likeCount }}
            </text>
            <text class="action-btn" @tap="onReply(reply)">回复</text>
          </view>
        </view>

        <view v-if="replies.length === 0" class="empty-state">
          <text class="empty-text">暂无回复，快来抢沙发</text>
        </view>
      </view>
    </scroll-view>

    <!-- Reply Input -->
    <view class="reply-input-bar">
      <input
        class="reply-input"
        v-model="replyText"
        :placeholder="replyTo ? `回复 ${replyTo.authorName}...` : '写回复...'"
        confirm-type="send"
        @confirm="onSubmitReply"
      />
      <text class="send-btn" @tap="onSubmitReply">发送</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface DiscussionInfo {
  id: string
  title: string
  content: string
  authorName: string
  createdAt: string
}

interface ReplyItem {
  id: string
  content: string
  authorName: string
  createdAt: string
  liked: boolean
  likeCount: number
}

const discussionId = ref('')
const discussion = ref<DiscussionInfo>({ id: '', title: '', content: '', authorName: '', createdAt: '' })
const replies = ref<ReplyItem[]>([])
const replyText = ref('')
const replyTo = ref<ReplyItem | null>(null)

async function loadDetail() {
  try {
    // TODO: import { getDiscussionDetail, getDiscussionReplies } from '@/api/class'
    // const res = await getDiscussionDetail(discussionId.value)
    // discussion.value = res.data
    // const repRes = await getDiscussionReplies(discussionId.value)
    // replies.value = repRes.data.list || []
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  }
}

function onLike(reply: ReplyItem) {
  reply.liked = !reply.liked
  reply.likeCount += reply.liked ? 1 : -1
  // TODO: API call to toggle like
}

function onReply(reply: ReplyItem) {
  replyTo.value = reply
}

async function onSubmitReply() {
  if (!replyText.value.trim()) return
  try {
    // TODO: import { postDiscussionReply } from '@/api/class'
    // await postDiscussionReply({ discussionId: discussionId.value, content: replyText.value, replyToId: replyTo.value?.id })
    uni.showToast({ title: '回复成功', icon: 'success' })
    replyText.value = ''
    replyTo.value = null
    loadDetail()
  } catch (e) {
    uni.showToast({ title: '回复失败', icon: 'none' })
  }
}

onMounted(() => {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1] as any
  discussionId.value = current.options?.id || ''
  loadDetail()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  display: flex;
  flex-direction: column;
}

.content-area {
  flex: 1;
  padding-bottom: 120rpx;
}

.discuss-header {
  background-color: #fff;
  padding: 32rpx 24rpx;
  margin-bottom: 20rpx;
}

.discuss-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.discuss-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.info-author {
  font-size: 26rpx;
  color: #1677ff;
}

.info-time {
  font-size: 24rpx;
  color: #999;
}

.discuss-body {
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}

.replies-section {
  background-color: #fff;
  padding: 24rpx;
}

.replies-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.reply-card {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.reply-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.reply-author {
  font-size: 26rpx;
  color: #1677ff;
}

.reply-time {
  font-size: 22rpx;
  color: #999;
}

.reply-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  margin-bottom: 12rpx;
  display: block;
}

.reply-actions {
  display: flex;
  gap: 40rpx;
}

.action-btn {
  font-size: 24rpx;
  color: #999;
}

.empty-state {
  padding: 60rpx 0;
  text-align: center;
}

.empty-text {
  font-size: 26rpx;
  color: #999;
}

.reply-input-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #fff;
  border-top: 1rpx solid #eee;
  gap: 16rpx;
}

.reply-input {
  flex: 1;
  height: 72rpx;
  background-color: #f5f6fa;
  border-radius: 36rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
}

.send-btn {
  font-size: 28rpx;
  color: #1677ff;
  font-weight: 600;
  padding: 0 12rpx;
}
</style>

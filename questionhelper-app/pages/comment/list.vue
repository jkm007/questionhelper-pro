<template>
  <view class="comment-page">
    <!-- Comment List -->
    <scroll-view
      class="comment-list"
      scroll-y
      :scroll-top="scrollTop"
      @scrolltolower="loadMore"
    >
      <view v-if="comments.length === 0 && !loading" class="empty-state">
        <text class="empty-text">No comments yet</text>
      </view>

      <view
        v-for="comment in comments"
        :key="comment.id"
        class="comment-item"
      >
        <view class="comment-main" @tap="onReply(comment)">
          <image class="avatar" :src="comment.avatar || '/static/default-avatar.png'" mode="aspectFill" />
          <view class="comment-body">
            <view class="comment-header">
              <text class="nickname">{{ comment.nickname }}</text>
              <text class="time">{{ comment.createTime }}</text>
            </view>
            <text class="content">{{ comment.content }}</text>
            <view v-if="comment.images && comment.images.length > 0" class="comment-images">
              <image
                v-for="(img, idx) in comment.images"
                :key="idx"
                class="comment-image"
                :src="img"
                mode="aspectFill"
                @tap.stop="previewImage(comment.images, idx)"
              />
            </view>
            <view class="comment-actions">
              <view class="action-btn" @tap.stop="likeComment(comment)">
                <image
                  class="action-icon"
                  :src="comment.isLiked ? '/static/icon-liked.png' : '/static/icon-like.png'"
                  mode="aspectFit"
                />
                <text class="action-text">{{ comment.likeCount || 0 }}</text>
              </view>
              <view class="action-btn" @tap.stop="onReply(comment)">
                <image class="action-icon" src="/static/icon-reply.png" mode="aspectFit" />
                <text class="action-text">{{ comment.replyCount || 0 }}</text>
              </view>
            </view>

            <!-- Nested Replies (1 level) -->
            <view v-if="comment.replies && comment.replies.length > 0" class="replies-wrap">
              <view
                v-for="reply in comment.replies"
                :key="reply.id"
                class="reply-item"
              >
                <image class="reply-avatar" :src="reply.avatar || '/static/default-avatar.png'" mode="aspectFill" />
                <view class="reply-body">
                  <view class="reply-header">
                    <text class="reply-nickname">{{ reply.nickname }}</text>
                    <text class="reply-time">{{ reply.createTime }}</text>
                  </view>
                  <text class="reply-content">
                    <text v-if="reply.replyNickname" class="reply-to">@{{ reply.replyNickname }} </text>
                    {{ reply.content }}
                  </text>
                </view>
              </view>
              <view
                v-if="comment.replyCount > (comment.replies?.length || 0)"
                class="view-more-replies"
                @tap.stop="loadReplies(comment)"
              >
                <text class="view-more-text">View {{ comment.replyCount - (comment.replies?.length || 0) }} more replies</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">Loading...</text>
      </view>
      <view v-if="noMore && comments.length > 0" class="no-more-wrap">
        <text class="no-more-text">No more comments</text>
      </view>
    </scroll-view>

    <!-- Comment Input -->
    <view class="input-bar">
      <view class="input-wrap">
        <input
          class="comment-input"
          v-model="inputText"
          :placeholder="replyTarget ? `Reply to ${replyTarget.nickname}` : 'Write a comment...'"
          :adjust-position="true"
          confirm-type="send"
          @confirm="submitComment"
        />
      </view>
      <view class="send-btn" :class="{ active: inputText.trim() }" @tap="submitComment">
        <text class="send-text">Send</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  getCommentList,
  createComment,
  likeComment as apiLikeComment,
  getReplyList
} from '@/api/comment'

interface Comment {
  id: string
  avatar: string
  nickname: string
  content: string
  createTime: string
  images?: string[]
  likeCount: number
  replyCount: number
  isLiked: boolean
  replies?: Comment[]
}

const targetId = ref('')
const targetType = ref('')
const comments = ref<Comment[]>([])
const loading = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const scrollTop = ref(0)
const inputText = ref('')
const replyTarget = ref<Comment | null>(null)

onLoad((options: any) => {
  targetId.value = options?.targetId || ''
  targetType.value = options?.targetType || ''
  fetchComments()
})

async function fetchComments() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getCommentList({
      targetId: targetId.value,
      targetType: targetType.value,
      page: page.value,
      pageSize
    })
    const list = res.data?.list || []
    if (page.value === 1) {
      comments.value = list
    } else {
      comments.value.push(...list)
    }
    if (list.length < pageSize) {
      noMore.value = true
    }
  } catch (e) {
    console.error('Failed to load comments', e)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  if (noMore.value || loading.value) return
  page.value++
  fetchComments()
}

async function submitComment() {
  const text = inputText.value.trim()
  if (!text) return

  try {
    await createComment({
      targetId: targetId.value,
      targetType: targetType.value,
      content: text,
      parentId: replyTarget.value?.id || undefined,
      replyUserId: replyTarget.value?.id || undefined
    })
    inputText.value = ''
    replyTarget.value = null
    page.value = 1
    noMore.value = false
    await fetchComments()
    scrollTop.value = Math.random()
  } catch (e) {
    console.error('Failed to submit comment', e)
    uni.showToast({ title: 'Failed to send', icon: 'none' })
  }
}

async function likeComment(comment: Comment) {
  try {
    await apiLikeComment(comment.id)
    comment.isLiked = !comment.isLiked
    comment.likeCount += comment.isLiked ? 1 : -1
  } catch (e) {
    console.error('Failed to like', e)
  }
}

function onReply(comment: Comment) {
  replyTarget.value = comment
}

async function loadReplies(comment: Comment) {
  try {
    const res = await getReplyList({
      commentId: comment.id,
      page: 1,
      pageSize: 50
    })
    comment.replies = res.data?.list || []
  } catch (e) {
    console.error('Failed to load replies', e)
  }
}

function previewImage(images: string[], index: number) {
  uni.previewImage({
    urls: images,
    current: index
  })
}
</script>

<style scoped>
.comment-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

.comment-list {
  flex: 1;
  padding: 20rpx;
  padding-bottom: 120rpx;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.comment-item {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  padding: 24rpx;
}

.comment-main {
  display: flex;
  flex-direction: row;
}

.avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 36rpx;
  flex-shrink: 0;
}

.comment-body {
  flex: 1;
  margin-left: 20rpx;
}

.comment-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.nickname {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
}

.time {
  font-size: 22rpx;
  color: #999;
}

.content {
  font-size: 28rpx;
  color: #333;
  margin-top: 12rpx;
  line-height: 1.6;
}

.comment-images {
  display: flex;
  flex-wrap: wrap;
  margin-top: 16rpx;
  gap: 12rpx;
}

.comment-image {
  width: 180rpx;
  height: 180rpx;
  border-radius: 8rpx;
}

.comment-actions {
  display: flex;
  flex-direction: row;
  margin-top: 16rpx;
  gap: 40rpx;
}

.action-btn {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.action-icon {
  width: 36rpx;
  height: 36rpx;
}

.action-text {
  font-size: 24rpx;
  color: #999;
  margin-left: 8rpx;
}

.replies-wrap {
  background-color: #f8f8f8;
  border-radius: 12rpx;
  padding: 16rpx;
  margin-top: 16rpx;
}

.reply-item {
  display: flex;
  flex-direction: row;
  margin-bottom: 16rpx;
}

.reply-item:last-child {
  margin-bottom: 0;
}

.reply-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 24rpx;
  flex-shrink: 0;
}

.reply-body {
  flex: 1;
  margin-left: 12rpx;
}

.reply-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.reply-nickname {
  font-size: 24rpx;
  color: #666;
  font-weight: 600;
}

.reply-time {
  font-size: 20rpx;
  color: #bbb;
}

.reply-content {
  font-size: 26rpx;
  color: #333;
  margin-top: 6rpx;
  line-height: 1.5;
}

.reply-to {
  color: #4a90d9;
}

.view-more-replies {
  padding-top: 12rpx;
}

.view-more-text {
  font-size: 24rpx;
  color: #4a90d9;
}

.loading-wrap,
.no-more-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx 0;
}

.loading-text,
.no-more-text {
  font-size: 24rpx;
  color: #999;
}

.input-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #fff;
  border-top: 1rpx solid #eee;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
}

.input-wrap {
  flex: 1;
  background-color: #f5f5f5;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
}

.comment-input {
  width: 100%;
  font-size: 28rpx;
  color: #333;
}

.send-btn {
  margin-left: 16rpx;
  padding: 0 24rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ccc;
  border-radius: 36rpx;
}

.send-btn.active {
  background-color: #4a90d9;
}

.send-text {
  font-size: 28rpx;
  color: #fff;
}
</style>

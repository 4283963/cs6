<template>
  <el-card class="log-card">
    <template #header>
      <div class="card-header">
        <span>📡 实时状态日志</span>
        <el-tag :type="connected ? 'success' : 'danger'">
          {{ connected ? '已连接' : '未连接' }}
        </el-tag>
      </div>
    </template>

    <div class="log-container" ref="logContainer">
      <div
        v-for="(log, index) in logs"
        :key="index"
        class="log-item"
        :class="{ 'log-on': log.includes('开启'), 'log-off': log.includes('关闭') }"
      >
        <span class="log-time">{{ getCurrentTime() }}</span>
        <span class="log-content">{{ log }}</span>
      </div>
      <div v-if="logs.length === 0" class="empty-log">
        等待定时任务触发...
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const connected = ref(false)
const logs = ref([])
const logContainer = ref(null)
let eventSource = null

const getCurrentTime = () => {
  const now = new Date()
  return now.toLocaleTimeString('zh-CN', { hour12: false })
}

const scrollToBottom = () => {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const connectSSE = () => {
  eventSource = new EventSource('/api/status/stream')

  eventSource.onopen = () => {
    connected.value = true
    logs.value.push('✅ 已连接到服务器，实时监听状态变化')
    scrollToBottom()
  }

  eventSource.onmessage = (event) => {
    logs.value.push(event.data)
    if (logs.value.length > 100) {
      logs.value = logs.value.slice(-50)
    }
    scrollToBottom()
  }

  eventSource.onerror = () => {
    connected.value = false
    logs.value.push('❌ 连接断开，5秒后重试...')
    scrollToBottom()
    setTimeout(() => {
      if (eventSource) {
        eventSource.close()
      }
      connectSSE()
    }, 5000)
  }
}

onMounted(() => {
  connectSSE()
})

onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
  }
})
</script>

<style scoped>
.log-card {
  background: #1e1e1e;
  border: none;
}

.log-card :deep(.el-card__header) {
  background: #2d2d2d;
  border-bottom: 1px solid #444;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
  color: #fff;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-container {
  height: 250px;
  overflow-y: auto;
  background: #1e1e1e;
  padding: 15px;
  font-family: 'Consolas', 'Monaco', monospace;
}

.log-item {
  padding: 6px 10px;
  margin-bottom: 5px;
  border-radius: 4px;
  background: #2d2d2d;
  color: #d4d4d4;
  font-size: 13px;
  border-left: 3px solid #666;
}

.log-on {
  border-left-color: #67c23a;
  background: rgba(103, 194, 58, 0.1);
}

.log-off {
  border-left-color: #909399;
  background: rgba(144, 147, 153, 0.1);
}

.log-time {
  color: #888;
  margin-right: 10px;
}

.log-content {
  color: #d4d4d4;
}

.empty-log {
  color: #666;
  text-align: center;
  padding: 50px;
  font-style: italic;
}

.log-container::-webkit-scrollbar {
  width: 6px;
}

.log-container::-webkit-scrollbar-track {
  background: #2d2d2d;
}

.log-container::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 3px;
}
</style>

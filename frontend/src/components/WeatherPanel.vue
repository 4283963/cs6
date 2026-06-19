<template>
  <el-card class="weather-card">
    <template #header>
      <div class="card-header">
        <span>🌤️ 天气与光照度监控</span>
        <el-tag :type="weatherType" size="small">
          {{ weatherStatus.condition === '暴雨' ? '🌧️ 暴雨' : '☀️ 正常' }}
        </el-tag>
      </div>
    </template>

    <div class="weather-content">
      <div class="illuminance-display">
        <div class="illuminance-value" :class="{ 'low-light': isLowLight }">
          {{ weatherStatus.illuminance?.toFixed(1) || '--' }}
        </div>
        <div class="illuminance-unit">lux</div>
      </div>

      <el-progress
        :percentage="illuminancePercent"
        :color="progressColor"
        :show-text="false"
        class="illuminance-bar"
      />

      <div class="threshold-info">
        <span>低光照阈值: <el-tag type="danger" size="small">20 lux</el-tag></span>
        <span v-if="isLowLight" class="warning-text">
          ⚠️ 光照不足，路灯将提前30分钟开启
        </span>
      </div>

      <el-divider />

      <div class="control-section">
        <div class="control-label">模拟天气场景：</div>
        <div class="control-buttons">
          <el-button type="success" @click="handleSimulateNormal" :loading="loading">
            ☀️ 正常天气
          </el-button>
          <el-button type="primary" @click="handleSimulateStorm" :loading="loading">
            🌧️ 暴雨天气
          </el-button>
        </div>
      </div>

      <el-divider />

      <div class="custom-input">
        <el-form :inline="true" :model="customForm" @submit.prevent>
          <el-form-item label="自定义光照度">
            <el-input-number
              v-model="customForm.illuminance"
              :min="0"
              :max="10000"
              :step="10"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="天气状况">
            <el-select v-model="customForm.condition" style="width: 120px">
              <el-option label="正常" value="normal" />
              <el-option label="阴天" value="阴天" />
              <el-option label="暴雨" value="暴雨" />
              <el-option label="雾霾" value="雾霾" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleUpdateWeather" :loading="loading">
              应用
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="update-time" v-if="weatherStatus.updated_at">
        最后更新: {{ formatDate(weatherStatus.updated_at) }}
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getWeather, updateWeather, simulateStorm, simulateNormal } from '../api'

const loading = ref(false)
const weatherStatus = ref({
  illuminance: 1000,
  condition: 'normal',
  updated_at: null
})

const customForm = ref({
  illuminance: 1000,
  condition: 'normal'
})

let eventSource = null

const isLowLight = computed(() => weatherStatus.value.illuminance < 20)

const weatherType = computed(() => {
  if (weatherStatus.value.condition === '暴雨') return 'info'
  if (isLowLight.value) return 'danger'
  return 'success'
})

const illuminancePercent = computed(() => {
  const max = 1000
  return Math.min(100, (weatherStatus.value.illuminance / max) * 100)
})

const progressColor = computed(() => {
  if (isLowLight.value) return '#f56c6c'
  if (weatherStatus.value.illuminance < 100) return '#e6a23c'
  return '#67c23a'
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const fetchWeather = async () => {
  try {
    const res = await getWeather()
    weatherStatus.value = res.data.data
    customForm.value.illuminance = res.data.data.illuminance
    customForm.value.condition = res.data.data.condition
  } catch (err) {
    ElMessage.error('获取天气状态失败')
  }
}

const handleSimulateStorm = async () => {
  loading.value = true
  try {
    await simulateStorm()
    ElMessage.success('已模拟暴雨天气，光照度 15 lux')
    await fetchWeather()
  } catch (err) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}

const handleSimulateNormal = async () => {
  loading.value = true
  try {
    await simulateNormal()
    ElMessage.success('已恢复正常天气，光照度 1000 lux')
    await fetchWeather()
  } catch (err) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}

const handleUpdateWeather = async () => {
  loading.value = true
  try {
    await updateWeather(customForm.value)
    ElMessage.success('天气状态更新成功')
    await fetchWeather()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '更新失败')
  } finally {
    loading.value = false
  }
}

const connectSSE = () => {
  eventSource = new EventSource('/api/weather/stream')
  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      weatherStatus.value = data
      customForm.value.illuminance = data.illuminance
      customForm.value.condition = data.condition
    } catch (e) {
      console.error('Parse weather SSE error:', e)
    }
  }
}

onMounted(() => {
  fetchWeather()
  connectSSE()
})

onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
  }
})
</script>

<style scoped>
.weather-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.weather-content {
  text-align: center;
}

.illuminance-display {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 8px;
  margin: 20px 0;
}

.illuminance-value {
  font-size: 48px;
  font-weight: bold;
  color: #67c23a;
  transition: color 0.3s;
}

.illuminance-value.low-light {
  color: #f56c6c;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.illuminance-unit {
  font-size: 20px;
  color: #909399;
}

.illuminance-bar {
  margin: 10px 0 20px;
}

.threshold-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
  font-size: 14px;
  color: #606266;
}

.warning-text {
  color: #f56c6c;
  font-weight: bold;
}

.control-section {
  text-align: left;
}

.control-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.control-buttons {
  display: flex;
  gap: 10px;
}

.custom-input {
  text-align: left;
}

.update-time {
  margin-top: 15px;
  font-size: 12px;
  color: #909399;
  text-align: right;
}
</style>

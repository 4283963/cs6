<template>
  <div class="app-container">
    <el-container>
      <el-header class="header">
        <h1>💡 路灯节能控制系统</h1>
        <span class="subtitle">产业园区智能照明管理平台</span>
      </el-header>

      <el-main>
        <el-row :gutter="20">
          <el-col :span="8">
            <ScheduleForm @refresh="fetchSchedules" />
          </el-col>
          <el-col :span="8">
            <WeatherPanel />
          </el-col>
          <el-col :span="8">
            <el-card class="info-card">
              <template #header>
                <div class="card-header">
                  <span>📊 系统说明</span>
                </div>
              </template>
              <div class="info-content">
                <el-descriptions :column="1" border size="small">
                  <el-descriptions-item label="光照阈值">
                    <el-tag type="danger">20 lux</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="提前时间">
                    <el-tag type="primary">30 分钟</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="工作模式">
                    <span>定时 + 动态修正</span>
                  </el-descriptions-item>
                </el-descriptions>
                <el-divider />
                <div class="logic-desc">
                  <p><strong>工作逻辑：</strong></p>
                  <ol>
                    <li>每天在开灯时间前 <b>30 分钟</b> 进行预检查</li>
                    <li>如果光照度 < <b>20 lux</b>（暴雨/极暗），立即开灯</li>
                    <li>否则按设定时间正常开灯</li>
                    <li>当天已提前开过灯的，到点不再重复开启</li>
                    <li>关灯时间不受天气影响，按设定执行</li>
                  </ol>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-row style="margin-top: 20px">
          <el-col :span="24">
            <ScheduleList :schedules="schedules" @refresh="fetchSchedules" />
          </el-col>
        </el-row>

        <el-row style="margin-top: 20px">
          <el-col :span="24">
            <StatusLog />
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSchedules } from './api'
import ScheduleForm from './components/ScheduleForm.vue'
import ScheduleList from './components/ScheduleList.vue'
import WeatherPanel from './components/WeatherPanel.vue'
import StatusLog from './components/StatusLog.vue'

const schedules = ref([])

const fetchSchedules = async () => {
  try {
    const res = await getSchedules()
    schedules.value = res.data.data
  } catch (err) {
    ElMessage.error('获取策略列表失败')
  }
}

onMounted(() => {
  fetchSchedules()
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.header {
  background: white;
  border-radius: 12px;
  text-align: center;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.header h1 {
  margin: 0;
  color: #333;
  font-size: 28px;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin-top: 5px;
}

.el-main {
  padding: 0;
}

.info-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
}

.info-content {
  font-size: 14px;
}

.logic-desc p {
  margin: 0 0 10px 0;
  color: #303133;
}

.logic-desc ol {
  margin: 0;
  padding-left: 20px;
  color: #606266;
}

.logic-desc li {
  margin-bottom: 5px;
  line-height: 1.6;
}
</style>

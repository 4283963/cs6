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
          <el-col :span="16">
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
</style>

<template>
  <el-card class="list-card">
    <template #header>
      <div class="card-header">
        <span>📋 定时策略列表</span>
        <el-button type="primary" link @click="$emit('refresh')">
          刷新
        </el-button>
      </div>
    </template>

    <el-table :data="schedules" stripe style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="group_name" label="群组名称" min-width="150" />
      <el-table-column label="开启时间" width="120">
        <template #default="{ row }">
          <el-tag type="success">{{ row.on_time }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="关闭时间" width="120">
        <template #default="{ row }">
          <el-tag type="info">{{ row.off_time }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '运行中' : '已停止' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" min-width="160">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="danger" link @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="schedules.length === 0 && !loading" description="暂无定时策略" />
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteSchedule } from '../api'

const props = defineProps({
  schedules: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['refresh'])

const loading = ref(false)

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除群组 [${row.group_name}] 的定时策略吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    loading.value = true
    await deleteSchedule(row.id)
    ElMessage.success('删除成功')
    emit('refresh')
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(err.response?.data?.error || '删除失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.list-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

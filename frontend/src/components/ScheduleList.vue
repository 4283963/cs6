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
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button type="danger" link @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="schedules.length === 0 && !loading" description="暂无定时策略" />

    <el-dialog
      v-model="editDialogVisible"
      title="✏️ 编辑定时策略"
      width="400px"
    >
      <el-form :model="editForm" label-width="80px" ref="editFormRef">
        <el-form-item label="群组名称">
          <el-input v-model="editForm.group_name" />
        </el-form-item>
        <el-form-item
          label="开启时间"
          :rules="[{ required: true, message: '请选择开启时间', trigger: 'change' }]"
        >
          <el-time-picker
            v-model="editForm.on_time"
            format="HH:mm"
            value-format="HH:mm"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item
          label="关闭时间"
          :rules="[{ required: true, message: '请选择关闭时间', trigger: 'change' }]"
        >
          <el-time-picker
            v-model="editForm.off_time"
            format="HH:mm"
            value-format="HH:mm"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updating">
          确定更新
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteSchedule, updateSchedule } from '../api'

const props = defineProps({
  schedules: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['refresh'])

const loading = ref(false)
const updating = ref(false)
const editDialogVisible = ref(false)
const editFormRef = ref(null)
const editingId = ref(null)
const editForm = ref({
  group_name: '',
  on_time: '',
  off_time: ''
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const handleEdit = (row) => {
  editingId.value = row.id
  editForm.value = {
    group_name: row.group_name,
    on_time: row.on_time,
    off_time: row.off_time
  }
  editDialogVisible.value = true
}

const handleUpdate = async () => {
  if (!editingId.value) return

  updating.value = true
  try {
    await updateSchedule(editingId.value, editForm.value)
    ElMessage.success('更新成功，旧定时任务已销毁')
    editDialogVisible.value = false
    emit('refresh')
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '更新失败')
  } finally {
    updating.value = false
  }
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
    ElMessage.success('删除成功，定时任务已销毁')
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

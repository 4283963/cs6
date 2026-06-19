<template>
  <el-card class="form-card">
    <template #header>
      <div class="card-header">
        <span>➕ 创建定时策略</span>
      </div>
    </template>

    <el-form :model="form" label-width="80px" ref="formRef">
      <el-form-item
        label="群组名称"
        prop="group_name"
        :rules="[{ required: true, message: '请输入群组名称', trigger: 'blur' }]"
      >
        <el-input
          v-model="form.group_name"
          placeholder="例如：A区路灯"
          clearable
        />
      </el-form-item>

      <el-form-item
        label="开启时间"
        prop="on_time"
        :rules="[{ required: true, message: '请选择开启时间', trigger: 'change' }]"
      >
        <el-time-picker
          v-model="form.on_time"
          placeholder="选择开启时间"
          format="HH:mm"
          value-format="HH:mm"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item
        label="关闭时间"
        prop="off_time"
        :rules="[{ required: true, message: '请选择关闭时间', trigger: 'change' }]"
      >
        <el-time-picker
          v-model="form.off_time"
          placeholder="选择关闭时间"
          format="HH:mm"
          value-format="HH:mm"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="loading" style="width: 100%">
          创建策略
        </el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createSchedule } from '../api'

const emit = defineEmits(['refresh'])

const formRef = ref(null)
const loading = ref(false)
const form = ref({
  group_name: '',
  on_time: '',
  off_time: ''
})

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await createSchedule(form.value)
      ElMessage.success('策略创建成功')
      form.value = { group_name: '', on_time: '', off_time: '' }
      formRef.value.resetFields()
      emit('refresh')
    } catch (err) {
      ElMessage.error(err.response?.data?.error || '创建失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.form-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
}
</style>

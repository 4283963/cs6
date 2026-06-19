import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 5000
})

export const createSchedule = (data) => {
  return request.post('/schedule/create', data)
}

export const getSchedules = () => {
  return request.get('/schedule/list')
}

export const deleteSchedule = (id) => {
  return request.delete(`/schedule/${id}`)
}

export default request

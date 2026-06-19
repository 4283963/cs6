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

export const updateSchedule = (id, data) => {
  return request.put(`/schedule/${id}`, data)
}

export const getWeather = () => {
  return request.get('/weather/current')
}

export const updateWeather = (data) => {
  return request.put('/weather/update', data)
}

export const simulateStorm = () => {
  return request.post('/weather/simulate/storm')
}

export const simulateNormal = () => {
  return request.post('/weather/simulate/normal')
}

export default request

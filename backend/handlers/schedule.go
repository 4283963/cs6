package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"streetlight-controller/database"
	"streetlight-controller/models"
	"streetlight-controller/scheduler"

	"github.com/gin-gonic/gin"
)

func isValidTimeFormat(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return false
	}
	min, err := strconv.Atoi(parts[1])
	if err != nil || min < 0 || min > 59 {
		return false
	}
	return true
}

func CreateSchedule(c *gin.Context) {
	var req models.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	if !isValidTimeFormat(req.OnTime) || !isValidTimeFormat(req.OffTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "时间格式错误，请使用 HH:mm 格式（如 18:00）",
		})
		return
	}

	var existing models.Schedule
	result := database.DB.Where("group_name = ?", req.GroupName).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "该群组已存在定时策略",
		})
		return
	}

	schedule := models.Schedule{
		GroupName: req.GroupName,
		OnTime:    req.OnTime,
		OffTime:   req.OffTime,
		Status:    "active",
	}

	result = database.DB.Create(&schedule)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建策略失败: " + result.Error.Error(),
		})
		return
	}

	err := scheduler.GlobalScheduler.AddSchedule(&schedule)
	if err != nil {
		database.DB.Delete(&schedule)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "启动定时任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "定时策略创建成功",
		"data":    schedule,
	})
}

func GetSchedules(c *gin.Context) {
	var schedules []models.Schedule
	result := database.DB.Find(&schedules)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": schedules,
	})
}

func DeleteSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var schedule models.Schedule
	result := database.DB.First(&schedule, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "策略不存在",
		})
		return
	}

	scheduler.GlobalScheduler.RemoveSchedule(uint(id))

	result = database.DB.Delete(&schedule)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "策略删除成功",
	})
}

func GetStatusStream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	c.SSEvent("message", "连接成功，等待状态更新...")
	flusher.Flush()

	for {
		select {
		case msg := <-scheduler.GlobalScheduler.StatusCh:
			c.SSEvent("message", msg)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

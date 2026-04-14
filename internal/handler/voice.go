package handler

import (
	"log"
	"net/http"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// VoiceRecognize 语音识别
func VoiceRecognize(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有音频文件"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取音频失败"})
		return
	}
	defer f.Close()

	text, err := service.Voice.Recognize(f)
	if err != nil {
		log.Printf("VoiceRecognize 失败: err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"text": text})
}

// VoiceStatus 语音识别状态
func VoiceStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"available": service.Voice.IsAvailable(),
	})
}

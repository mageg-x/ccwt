package service

import (
	"fmt"
	"io"
	"log"

	"github.com/ccwt/ccwt/internal/config"
)

// VoiceManager 语音识别管理（百度在线 ASR）
type VoiceManager struct{}

var Voice = &VoiceManager{}

// Recognize 语音识别，将 WAV(16k/mono) 音频转换为文本
func (v *VoiceManager) Recognize(audioReader io.Reader) (string, error) {
	if !config.Cfg.Voice.Enabled {
		return "", fmt.Errorf("语音识别未启用")
	}

	audio, err := io.ReadAll(audioReader)
	if err != nil {
		return "", fmt.Errorf("读取音频失败: %v", err)
	}
	if len(audio) == 0 {
		return "", fmt.Errorf("音频内容为空")
	}

	text, err := baiduClient.recognizeWav16k(audio)
	if err != nil {
		log.Printf("百度语音识别失败: err=%v", err)
		return "", err
	}
	return text, nil
}

// IsAvailable 检查语音识别是否可用
func (v *VoiceManager) IsAvailable() bool {
	ok, _ := v.Status()
	return ok
}

// Status 返回语音可用状态与原因
func (v *VoiceManager) Status() (bool, string) {
	if !config.Cfg.Voice.Enabled {
		return false, "voice.disabled"
	}
	if _, err := baiduClient.getToken(); err != nil {
		return false, "baidu.token.failed"
	}
	return true, ""
}

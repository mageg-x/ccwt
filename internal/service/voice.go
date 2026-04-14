package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ccwt/ccwt/internal/config"
)

// VoiceManager 语音识别管理
type VoiceManager struct{}

var Voice = &VoiceManager{}

// Recognize 语音识别，将音频转换为文本
func (v *VoiceManager) Recognize(audioReader io.Reader) (string, error) {
	if !config.Cfg.Voice.Enabled {
		return "", fmt.Errorf("语音识别未启用")
	}

	modelPath := config.Cfg.Voice.ModelPath
	if modelPath == "" {
		modelPath = filepath.Join(config.DataDir(), "models", "ggml-base.bin")
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("语音模型不存在: %s", modelPath)
	}

	// 保存音频到临时文件
	tmp, err := os.CreateTemp("", "ccwt-voice-*.wav")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := io.Copy(tmp, audioReader); err != nil {
		return "", fmt.Errorf("保存音频失败: %v", err)
	}
	tmp.Close()

	// 调用 whisper.cpp
	cmd := exec.Command("whisper-cpp",
		"--model", modelPath,
		"--language", "zh",
		"--output-txt",
		"--no-timestamps",
		"--file", tmp.Name(),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("whisper-cpp 执行失败: err=%v output=%s", err, string(out))
		return "", fmt.Errorf("语音识别失败: %v", err)
	}

	text := strings.TrimSpace(string(out))
	return text, nil
}

// IsAvailable 检查语音识别是否可用
func (v *VoiceManager) IsAvailable() bool {
	if !config.Cfg.Voice.Enabled {
		return false
	}
	_, err := exec.LookPath("whisper-cpp")
	return err == nil
}

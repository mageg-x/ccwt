package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	baiduAppID    = "7634735"
	baiduAPIKey   = "1QLOGTEIMqjo40Bq6yt7yEcO"
	baiduSecret   = "yDVvrYuDgBCqPya18k4uESE6qvtg58wA"
	baiduTokenURL = "https://aip.baidubce.com/oauth/2.0/token"
	baiduASRURL   = "https://vop.baidu.com/server_api"
)

type baiduASR struct {
	client *http.Client

	mu      sync.Mutex
	token   string
	expires time.Time
}

var baiduClient = &baiduASR{
	client: &http.Client{Timeout: 45 * time.Second},
}

func (b *baiduASR) getToken() (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.token != "" && time.Now().Before(b.expires.Add(-60*time.Second)) {
		return b.token, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", baiduAPIKey)
	form.Set("client_secret", baiduSecret)

	req, _ := http.NewRequest(http.MethodPost, baiduTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("百度 token 请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("解析百度 token 响应失败: %w", err)
	}
	if tr.AccessToken == "" {
		if tr.Error != "" {
			return "", fmt.Errorf("百度 token 失败: %s %s", tr.Error, tr.ErrorDesc)
		}
		return "", fmt.Errorf("百度 token 响应无 access_token")
	}

	if tr.ExpiresIn <= 0 {
		tr.ExpiresIn = 3600
	}
	b.token = tr.AccessToken
	b.expires = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return b.token, nil
}

func (b *baiduASR) recognizeWav16k(audio []byte) (string, error) {
	token, err := b.getToken()
	if err != nil {
		return "", err
	}

	payload := map[string]any{
		"format":  "wav",
		"rate":    16000,
		"channel": 1,
		"cuid":    "ccwt-" + baiduAppID,
		"token":   token,
		"dev_pid": 1537,
		"speech":  base64.StdEncoding.EncodeToString(audio),
		"len":     len(audio),
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baiduASRURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("百度语音识别请求失败: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var ar struct {
		ErrNo  int      `json:"err_no"`
		ErrMsg string   `json:"err_msg"`
		Result []string `json:"result"`
	}
	if err := json.Unmarshal(raw, &ar); err != nil {
		return "", fmt.Errorf("解析百度语音响应失败: %w", err)
	}
	if ar.ErrNo != 0 {
		return "", fmt.Errorf("百度语音识别失败: %d %s", ar.ErrNo, ar.ErrMsg)
	}
	if len(ar.Result) == 0 {
		return "", fmt.Errorf("百度语音识别结果为空")
	}
	return strings.TrimSpace(ar.Result[0]), nil
}

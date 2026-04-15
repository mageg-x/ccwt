package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/armon/go-socks5"
	"github.com/ccwt/ccwt/internal/config"
)

// ProxyManager SOCKS5 代理管理
type ProxyManager struct {
	running  bool
	listener net.Listener
	cancel   context.CancelFunc
	host     string
	port     int
	mu       sync.Mutex
}

var Proxy = &ProxyManager{}

// Status 获取代理状态
func (p *ProxyManager) Status() (bool, string, string, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running && p.listener != nil {
		return true, p.listener.Addr().String(), p.host, p.port
	}
	return false, "", "", 0
}

// Start 启动 SOCKS5 代理
func (p *ProxyManager) Start(host string, port int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("代理已在运行")
	}

	if host == "" {
		host = "0.0.0.0"
	}

	if port <= 0 {
		port = config.Cfg.Proxy.Port
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("端口范围无效: %d", port)
	}

	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("创建SOCKS5服务失败: %v", err)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听端口失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.listener = ln
	p.cancel = cancel
	p.running = true
	p.host = host
	p.port = port

	go func() {
		log.Printf("SOCKS5 代理启动: %s", addr)
		if err := server.Serve(ln); err != nil {
			select {
			case <-ctx.Done():
				// 正常关闭
			default:
				log.Printf("SOCKS5 代理错误: %v", err)
			}
		}
		p.mu.Lock()
		p.running = false
		p.host = ""
		p.port = 0
		p.mu.Unlock()
	}()

	return nil
}

// Stop 停止 SOCKS5 代理
func (p *ProxyManager) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return fmt.Errorf("代理未运行")
	}

	if p.cancel != nil {
		p.cancel()
	}
	if p.listener != nil {
		p.listener.Close()
	}
	p.running = false
	p.host = ""
	p.port = 0
	log.Printf("SOCKS5 代理已停止")
	return nil
}

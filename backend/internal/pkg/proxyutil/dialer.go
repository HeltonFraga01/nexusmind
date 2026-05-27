// Package proxyutil 提供统一的代理配置功能
//
// 支持的代理协议：
//   - HTTP/HTTPS: 通过 Transport.Proxy 设置
//   - SOCKS5: 通过 Transport.DialContext 设置（客户端本地解析 DNS）
//   - SOCKS5H: 通过 Transport.DialContext 设置（代理端远程解析 DNS，推荐）
//
// 注意：proxyurl.Parse() 会自动将 socks5:// 升级为 socks5h://，
// 确保 DNS 也由代理端解析，防止 DNS 泄漏。
package proxyutil

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// ConfigureTransportProxy 根据代理 URL 配置 Transport
//
// 支持的协议：
//   - http/https: 设置 transport.Proxy
//   - socks5: 设置 transport.DialContext（客户端本地解析 DNS）
//   - socks5h: 设置 transport.DialContext（代理端远程解析 DNS，推荐）
//
// 参数：
//   - transport: 需要配置 of http.Transport
//   - proxyURL: 代理地址，nil 表示直连
//
// 返回：
//   - error: 代理配置错误（协议不支持或 dialer 创建失败）
//
// 优化说明：
//   对于 SOCKS 代理，我们通过自定义 DialContext 对本地主机、Docker 内部服务（无点域名）以及私有 IP 进行直连旁路，
//   避免本地区域/内部网络流量被不必要且可能失败地路由至外部 SOCKS 代理。
func ConfigureTransportProxy(transport *http.Transport, proxyURL *url.URL) error {
	if proxyURL == nil {
		return nil
	}

	scheme := strings.ToLower(proxyURL.Scheme)
	switch scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(proxyURL)
		return nil

	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return fmt.Errorf("create socks5 dialer: %w", err)
		}

		directDialer := &net.Dialer{
			Timeout: 5 * time.Second,
		}

		// 包装 DialContext，识别本地及内部地址以进行 Proxy Bypass
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}

			if isBypassAddr(host) {
				// 本地/ Docker 内部服务，直接连接
				return directDialer.DialContext(ctx, network, addr)
			}

			// 外部请求，使用 SOCKS 代理
			if contextDialer, ok := dialer.(proxy.ContextDialer); ok {
				return contextDialer.DialContext(ctx, network, addr)
			}
			return dialer.Dial(network, addr)
		}
		return nil

	default:
		return fmt.Errorf("unsupported proxy scheme: %s", scheme)
	}
}

// isBypassAddr 检查地址是否为本地环回、Docker 内部无点域名或私有 IP 网段，以此决定是否绕过代理。
func isBypassAddr(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return true
	}

	// 1. 环回地址及无点域名（Docker 内部服务，如 cortexx-omniroute_cortexx-omniroute）
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}
	if !strings.Contains(host, ".") {
		return true
	}

	// 2. 本地局域网 Multicast DNS 域名
	if strings.HasSuffix(host, ".local") {
		return true
	}

	// 3. 私有 IP 地址网段校验（RFC 1918 & RFC 4193 Unique Local Address）
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return true
		}
		// IPv4 私有网段 (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
		if ip4 := ip.To4(); ip4 != nil {
			return ip4[0] == 10 ||
				(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
				(ip4[0] == 192 && ip4[1] == 168)
		}
		// IPv6 局域网私有网段 Unique Local Address (fc00::/7)
		if len(ip) == 16 {
			return (ip[0]&0xfe) == 0xfc
		}
	}

	return false
}

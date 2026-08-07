package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Rule 定义一条 mock 规则。
type Rule struct {
	Method  string            `json:"method"`  // GET/POST/...
	Path    string            `json:"path"`    // /api/users
	Status  int               `json:"status"`  // 200/404/...
	Headers map[string]string `json:"headers"` // 响应头
	Body    string            `json:"body"`    // 响应体
	Delay   string            `json:"delay"`   // 模拟延迟如 "500ms"
}

type Config struct {
	Name  string `json:"name"`
	Port  int    `json:"port"`
	Rules []Rule `json:"rules"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	return &c, nil
}

// Match 找匹配的规则。
func (c *Config) Match(method, path string) *Rule {
	for i := range c.Rules {
		r := &c.Rules[i]
		if r.Method != "" && !strings.EqualFold(r.Method, method) {
			continue
		}
		if r.Path != "" && r.Path != path {
			continue
		}
		return r
	}
	return nil
}

// Serve 启动 mock 服务器，返回 mux 方便测试。
func (c *Config) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rule := c.Match(r.Method, r.URL.Path)
		if rule == nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, `{"error":"no rule for %s %s"}`, r.Method, r.URL.Path)
			return
		}
		for k, v := range rule.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(rule.Status)
		fmt.Fprint(w, rule.Body)
	})
	return mux
}

var _ = fmt.Println

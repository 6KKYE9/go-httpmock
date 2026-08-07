package mock

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatch(t *testing.T) {
	c := &Config{
		Rules: []Rule{
			{Method: "GET", Path: "/api/health", Status: 200, Body: `{"ok":true}`},
			{Method: "POST", Path: "/api/users", Status: 201, Body: `{"id":1}`},
		},
	}
	r := c.Match("GET", "/api/health")
	if r == nil || r.Status != 200 {
		t.Fatal("应匹配到 health 规则")
	}
	r = c.Match("GET", "/nope")
	if r != nil {
		t.Fatal("不应匹配到规则")
	}
}

func TestHandler(t *testing.T) {
	c := &Config{
		Rules: []Rule{
			{Method: "GET", Path: "/ok", Status: 200, Body: "hello"},
		},
	}
	h := c.Handler()
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ok")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("状态码应为 200")
	}
}

func TestHandler404(t *testing.T) {
	c := &Config{}
	h := c.Handler()
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/nope")
	if resp.StatusCode != 404 {
		t.Fatalf("未匹配应返回 404")
	}
}

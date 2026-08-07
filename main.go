package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-httpmock/internal/mock"
)

func run() error {
	config := flag.String("c", "mock.json", "配置文件路径")
	port := flag.Int("port", 0, "端口（覆盖配置里的）")
	flag.Parse()

	cfg, err := mock.LoadConfig(*config)
	if err != nil {
		return fmt.Errorf("读配置失败: %w", err)
	}
	if *port > 0 {
		cfg.Port = *port
	}

	fmt.Printf("Mock 服务器 %s 监听 :%d\n", cfg.Name, cfg.Port)
	for _, r := range cfg.Rules {
		fmt.Printf("  %s %s → %d\n", r.Method, r.Path, r.Status)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), cfg.Handler()))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}

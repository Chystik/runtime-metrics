package main

import (
	"flag"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.AgentConfig) {
	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(&cfg.HTTPServer)
	_ = flag.Value(&cfg.PollInterval)
	_ = flag.Value(&cfg.ReportInterval)

	flag.Var(&cfg.HTTPServer, "a", "Net address host:port")
	flag.Var(&cfg.PollInterval, "p", "Poll interval in seconds")
	flag.Var(&cfg.ReportInterval, "r", "Report interval in seconds")
	flag.Parse()
}
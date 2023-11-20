package run

import (
	"context"
	"os"
	"runtime/pprof"

	"github.com/Chystik/runtime-metrics/config"
)

type profile struct {
	cfg config.ProfileConfig
	cpu *os.File
	mem *os.File
}

func NewProfile(cfg config.ProfileConfig) (*profile, error) {
	CPUFile, err := os.Create(cfg.CPUFilePath)
	if err != nil {
		return nil, err
	}

	MemFile, err := os.Create(cfg.MemFilePath)
	if err != nil {
		return nil, err
	}

	return &profile{
		cfg: cfg,
		cpu: CPUFile,
		mem: MemFile,
	}, nil
}

func (p *profile) Run(ctx context.Context) error {
	var err error

	if err = pprof.StartCPUProfile(p.cpu); err != nil {
		return err
	}

	// interrupt signal
	<-ctx.Done()

	pprof.StopCPUProfile()
	err = p.cpu.Close()
	if err != nil {
		return err
	}

	err = pprof.WriteHeapProfile(p.mem)
	if err != nil {
		return err
	}
	err = p.mem.Close()
	if err != nil {
		return err
	}

	return nil
}

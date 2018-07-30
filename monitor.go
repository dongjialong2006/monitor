package monitor

import (
	"context"
	"os"
	"time"

	"github.com/dongjialong2006/log"
	"github.com/shirou/gopsutil/process"
)

var def *monitor = nil

type monitor struct {
	cpu  float64
	mem  float32
	ctx  context.Context
	log  *log.Entry
	stop chan struct{}
}

func new(ctx context.Context, cpu float64, mem float32) *monitor {
	m := &monitor{
		ctx:  ctx,
		cpu:  5.0,
		mem:  5.0,
		log:  log.New("util/monitor"),
		stop: make(chan struct{}),
	}

	go m.watch()

	return m
}

func (m *monitor) watch() {
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if nil != err {
		m.log.Error(err)
		return
	}
	if nil == proc {
		m.log.Warnf("process id:%d is not exist.", pid)
		return
	}

	var cpuPercent float64 = 0.0
	var memPercent float32 = 0.0
	tick := time.Tick(time.Second)
	for {
		select {
		case <-m.stop:
			return
		case <-m.ctx.Done():
			return
		case <-tick:
			if cpuPercent, err = proc.CPUPercentWithContext(m.ctx); nil != err {
				m.log.Error(err)
			}
			if cpuPercent > m.cpu {
				m.log.WithField("pid", pid).Infof("cpu percent:%f.", cpuPercent)
			}

			if memPercent, err = proc.MemoryPercentWithContext(m.ctx); nil != err {
				m.log.Error(err)
			}

			if memPercent > m.mem {
				m.log.WithField("pid", pid).Infof("used memory percent:%f.", memPercent)
			}
		}
	}
}

func Watch(ctx context.Context, cpu float64, mem float32) {
	if nil == def {
		def = new(ctx, cpu, mem)
	}
}

func Stop() {
	if nil != def {
		close(def.stop)
	}
}

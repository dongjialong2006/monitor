package monitor

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/dongjialong2006/log"
	"github.com/shirou/gopsutil/process"
)

var rw sync.RWMutex
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
		cpu:  cpu,
		mem:  mem,
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

	var num time.Duration = time.Duration(1)
	var cpuPercent float64 = 0.0
	var memPercent float32 = 0.0

	var tick1 = time.Tick(time.Second * num)
	var tick2 = time.Tick(time.Second * num)

	for {
		select {
		case <-m.stop:
			return
		case <-m.ctx.Done():
			return
		case <-tick1:
			if cpuPercent, err = proc.CPUPercentWithContext(m.ctx); nil != err {
				num++
				tick1 = time.Tick(time.Second * num)
				if num == 120 {
					num = 1
				}
				m.log.Error(err)
			}
			if cpuPercent > m.cpu {
				m.log.WithField("pid", pid).Infof("cpu percent:%f.", cpuPercent)
			}
		case <-tick2:
			if memPercent, err = proc.MemoryPercentWithContext(m.ctx); nil != err {
				num++
				tick2 = time.Tick(time.Second * num)
				if num == 120 {
					num = 1
				}
				m.log.Error(err)
			}

			if memPercent > m.mem {
				m.log.WithField("pid", pid).Infof("used memory percent:%f.", memPercent)
			}
		}
	}
}

func Watch(ctx context.Context, cpu float64, mem float32) {
	rw.Lock()
	defer rw.Unlock()
	if nil == def {
		def = new(ctx, cpu, mem)
	}
}

func Stop() {
	rw.Lock()
	defer rw.Unlock()
	if nil != def {
		close(def.stop)
		def = nil
	}
}

package posixsignal

import (
	"github.com/cuizhaoyue/iams/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name 定义了服务关闭管理器的名称.
const Name = "PosixSignalManager"

// PosixSignalManager 是一种服务关闭管理器，它实现了ShutdownManager接口.
type PosixSignalManager struct {
	signals []os.Signal
}

var _ shutdown.ShutdownManager = &PosixSignalManager{}

func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{signals: sig}
}

// GetName 返回服务关闭管理器的名称.
func (pm *PosixSignalManager) GetName() string {

	return Name
}

// Start 关闭服务启动入口.
func (pm *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, pm.signals...)

		// 接收到信号前一直阻塞
		<-c

		// 调用服务关闭相关的操作函数
		gs.StartShutdown(pm)
	}()

	return nil
}

// PreShutdown does nothing.
func (pm *PosixSignalManager) PreShutdown() error {
	return nil
}

// PostShutdown 退出应用.
func (pm *PosixSignalManager) PostShutdown() error {
	os.Exit(0)

	return nil
}

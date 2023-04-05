package shutdown

import "sync"

// ShutdownManager 是一个由服务关闭管理器实现的接口.
type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	PreShutdown() error
	PostShutdown() error
}

// ShutdownCallback 是一个由执行关闭操作时的回调函数实现的接口.
type ShutdownCallback interface {
	OnShutdown(string) error
}

// ShutdownFunc 是一个helper类型，可以提供匿名函数作为ShutdownCallback实例
type ShutdownFunc func(string) error

// OnShutdown 定义在关闭服务操作触发时需要执行的操作.
func (f ShutdownFunc) OnShutdown(managerName string) error {
	return f(managerName)
}

// ErrorHandler 异步处理错误的接口
type ErrorHandler interface {
	OnError(error)
}

// ErrorFunc 是一个helper类型，可以提供匿名函数作为ErrorHandler实例.
type ErrorFunc func(error)

// OnError 定义发生错误时需要执行的操作.
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GSInterface 是一个优雅服务必须实现的接口.
// 它接收一个ShutdownManager实例，然后调用该实例的ShutdownStart函数完成关闭操作.
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(callback ShutdownCallback)
}

var _ GSInterface = &GracefuleShutdown{}

// GracefuleShutdown 用于实现优雅关闭的主要结构体.由这来处理ShutdownCallback实例和ShutdownManager实例.
type GracefuleShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

func New() *GracefuleShutdown {
	return &GracefuleShutdown{
		callbacks: make([]ShutdownCallback, 10),
		managers:  make([]ShutdownManager, 3),
	}
}

// Start 调用所有ShutdownManager的Start函数，执行服务关闭操作.
func (gs *GracefuleShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

// AddShutdownManager 添加一个ShutdownManager实例.
func (gs *GracefuleShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

// SetErrorHandler 设置ErrorHandler实例
func (gs *GracefuleShutdown) SetErrorHandler(errHandler ErrorHandler) {
	gs.errorHandler = errHandler
}

func (gs *GracefuleShutdown) StartShutdown(sm ShutdownManager) {
	// 执行关闭服务的预操作
	gs.ReportError(sm.PreShutdown())

	// 执行关闭服务相关的操作函数
	var wg sync.WaitGroup
	for _, callback := range gs.callbacks {
		wg.Add(1)
		go func(callback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(callback.OnShutdown(sm.GetName()))
		}(callback)
	}

	// 执行关闭服务完成之后的操作
	gs.ReportError(sm.PostShutdown())
}

// ReportError 用来汇报error
func (gs *GracefuleShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

func (gs *GracefuleShutdown) AddShutdownCallback(callback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, callback)
}

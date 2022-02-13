package shutdown

import "sync"

// Callback 是Shutdown时回调的接口
type Callback interface {
	OnShutdown(string) error
}

// OnShutdownFunc 为Callback提供匿名函数
type OnShutdownFunc func(string) error

// OnShutdown 定义了关机触发时需要运行的动作
func (f OnShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

type Manager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

type ErrorHandler interface {
	OnError(err error)
}

type ErrorFunc func(err error)

func (f ErrorFunc) OnError(err error) {
	f(err)
}

type GSInterface interface {
	StartShutdown(sm Manager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback Callback)
}

type GracefulShutdown struct {
	callbacks    []Callback
	managers     []Manager
	errorHandler ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]Callback, 0, 10),
		managers:  make([]Manager, 0, 3),
	}
}

func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

func (gs *GracefulShutdown) AddShutdownManager(manager Manager) {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback 添加shutdown回调
func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback Callback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// StartShutdown 发起Shutdown操作
// 1. 调用 ShutdownStart
// 2. 依次执行回调
// 3. 最后执行 ShutdownFinish
func (gs *GracefulShutdown) StartShutdown(sm Manager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback Callback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}
	wg.Done()

	gs.ReportError(sm.ShutdownFinish())
}

// ReportError Shutdown上报错误
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

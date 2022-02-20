package posixsignal

import (
	"github.com/shencw/login/pkg/shutdown"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const Name = "PosixSignalManager"

type Manager struct {
	signals []os.Signal
}

func NewPosixSignalManager(sig ...os.Signal) *Manager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &Manager{sig}
}

func (m *Manager) GetName() string {
	return Name
}

func (m *Manager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, m.signals...)
		<-c
		gs.StartShutdown(m)
	}()

	return nil
}

func (m *Manager) ShutdownStart() error {
	return nil
}

func (m *Manager) ShutdownFinish() error {
	log.Println("finish")
	os.Exit(0)

	return nil
}

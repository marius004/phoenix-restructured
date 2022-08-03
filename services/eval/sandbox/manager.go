package sandbox

import (
	"context"

	"github.com/marius004/phoenix-/internal"
	"golang.org/x/sync/semaphore"
)

type Manager struct {
	semaphore  *semaphore.Weighted
	evalConfig *internal.EvalConfig

	maxConcurrentSandboxes int64
	availableSandboxes     chan int
}

func (m *Manager) RunTask(ctx context.Context, task internal.Task) error {
	logger := internal.GetGlobalLoggerInstance()
	sandbox, err := m.getSandbox()

	if err != nil {
		logger.Println("Could not create sandbox", err)
		return err
	}

	defer m.ReleaseSandbox(sandbox)
	return task.Run(ctx, sandbox)
}

func (m *Manager) ReleaseSandbox(sandbox internal.Sandbox) {
	m.semaphore.Release(1)
	m.availableSandboxes <- sandbox.GetID()
}

func (m *Manager) Stop(ctx context.Context) error {
	if err := m.semaphore.Acquire(ctx, m.maxConcurrentSandboxes); err != nil {
		return err
	}

	close(m.availableSandboxes)
	return nil
}

func (m *Manager) getSandbox() (internal.Sandbox, error) {
	if err := m.semaphore.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}

	return m.newSandbox(<-m.availableSandboxes)
}

func (m *Manager) newSandbox(id int) (*Sandbox, error) {
	return newSandbox(id, m.evalConfig)
}

func NewManager(evalConfig *internal.EvalConfig) *Manager {
	manager := &Manager{
		evalConfig: evalConfig,

		semaphore:              semaphore.NewWeighted(int64(evalConfig.MaxSandboxes)),
		availableSandboxes:     make(chan int, evalConfig.MaxSandboxes),
		maxConcurrentSandboxes: int64(evalConfig.MaxSandboxes),
	}

	for i := 1; i <= int(evalConfig.MaxSandboxes); i++ {
		manager.availableSandboxes <- i
	}

	return manager
}

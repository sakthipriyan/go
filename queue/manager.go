package queue

import (
	"time"
    "os"
    "log"
)

type Manager struct {
	basedir string
	created time.Time
	queues  map[string]Queue
}

type SystemStats struct {
}

type QueueStats struct {
}

func prepDir(dir string) {
    _, err := os.Stat(dir)
    if os.IsNotExist(err) {
        err := os.MkdirAll(dir, 0700)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func loadQueues(dir string) map[string]Queue {

    return nil
}

func newManager(dir string) *Manager {
    prepDir(dir)
    queues := loadQueues(dir)
	return &Manager{
        dir,
        time.Now(),
        queues}
}

func (m *Manager) Status() SystemStats {
	return SystemStats{}
}

func (m *Manager) Enqueue(q string, data []byte) error {
	return nil
}

func (m *Manager) Dequeue(q string) ([]byte, error) {
	return []byte{}, nil
}

func (m *Manager) GetQueue(q string) (QueueStats, error) {
	return QueueStats{}, nil
}

func (m *Manager) CreateQueue(q string) (bool, error) {
	return true, nil
}

func (m *Manager) DeleteQueue(q string) (bool, error) {
	return true, nil
}

func (m *Manager) Shutdown() (bool, error) {
	return true, nil
}

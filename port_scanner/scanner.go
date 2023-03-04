package port_scanner

import (
	"sync"
)

// Внутренняя реализация сканера портов
type Scanner struct {
	// IP
	ip IP
	// Диапазон портов для сканирования
	portRange PortRange
	// Стратегии сканирование
	strategies []scanStrategy
	// Максимальное количество работающих горутин
	workersNum int
}

func newScanner(ip IP, portRange PortRange, strategies []scanStrategy, workersNum int) *Scanner {
	return &Scanner{ip, portRange, strategies, workersNum}
}

// Run запускает сканер и возвращает канал успешно выполненных задач сканирования
func (s *Scanner) Run() chan *scanTaskSuccess {
	tasksChan := make(chan *scanTask, s.workersNum)
	go s.taskSpawner(tasksChan)

	successesChan := make(chan *scanTaskSuccess)
	go s.runWorkers(tasksChan, successesChan)

	return successesChan
}

// taskSpawner генерирует задачи по диапазону портов
func (s *Scanner) taskSpawner(tasksChan chan *scanTask) {
	for _, p := range s.portRange.ports {
		port := p
		tasksChan <- &scanTask{ip: s.ip, port: port}
	}
	close(tasksChan)
}

// runWorkers занимается запуском горутин, которые будут выполнять задачи сканирования
func (s *Scanner) runWorkers(tasksChan chan *scanTask, successesChan chan *scanTaskSuccess) {
	wg := sync.WaitGroup{}
	for i := 0; i < s.workersNum; i++ {
		go s.spawnWorker(&wg, tasksChan, successesChan)
	}
	wg.Wait()
	close(successesChan)
}

// Цикл обработки задач воркером
func (s *Scanner) spawnWorker(wg *sync.WaitGroup, tasksChan chan *scanTask, successesChan chan *scanTaskSuccess) {
	wg.Add(1)
	defer wg.Done()
	for {
		task, more := <-tasksChan
		if more {
			s.processWorkerTask(task, successesChan)
		} else {
			break
		}
	}
}

// Обработка конкретной задачи воркером
func (s *Scanner) processWorkerTask(task *scanTask, successesChan chan *scanTaskSuccess) {
	for _, s := range s.strategies {
		// Здесь происходит последовательное выполнение задач сканирования по указанным стратегиям. В случае, если одна
		// из стратегий вернула успех, то считаем задачу выполненной
		if res := s.scan(task); res != nil {
			successesChan <- res
			break
		}
	}
}

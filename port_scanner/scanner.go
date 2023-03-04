package port_scanner

import (
	"context"
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
func (s *Scanner) Run(ctx context.Context) chan *scanTaskSuccess {
	tasksChan := make(chan *scanTask, s.workersNum)
	go s.taskSpawner(ctx, tasksChan)

	successesChan := make(chan *scanTaskSuccess)
	go s.runWorkers(ctx, tasksChan, successesChan)

	return successesChan
}

// taskSpawner генерирует задачи по диапазону портов
func (s *Scanner) taskSpawner(ctx context.Context, tasksChan chan *scanTask) {
	defer close(tasksChan)
	for _, p := range s.portRange.ports {
		select {
		case _ = <-ctx.Done():
			// Контекст отменен, больше создавать задачи не нужно
			return
		default:
			port := p
			tasksChan <- &scanTask{ip: s.ip, port: port}
		}
	}

}

// runWorkers занимается запуском горутин, которые будут выполнять задачи сканирования
func (s *Scanner) runWorkers(ctx context.Context, tasksChan chan *scanTask, successesChan chan *scanTaskSuccess) {
	wg := sync.WaitGroup{}
	for i := 0; i < s.workersNum; i++ {
		select {
		case _ = <-ctx.Done():
			// Контекст отменен, больше запускать воркеров не нужно
			break
		default:
			go s.spawnWorker(ctx, &wg, tasksChan, successesChan)
		}
	}
	wg.Wait()
	close(successesChan)
}

// Цикл обработки задач воркером
func (s *Scanner) spawnWorker(ctx context.Context, wg *sync.WaitGroup, tasksChan chan *scanTask, successesChan chan *scanTaskSuccess) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case _ = <-ctx.Done():
			// Контекст отменен, больше обрабатывать задачи не нужно
			return
		case task, more := <-tasksChan:
			if more {
				s.processWorkerTask(ctx, task, successesChan)
			} else {
				return
			}
		}
	}
}

// Обработка конкретной задачи воркером
func (s *Scanner) processWorkerTask(ctx context.Context, task *scanTask, successesChan chan *scanTaskSuccess) {
	for _, s := range s.strategies {
		select {
		case _ = <-ctx.Done():
			return
		default:
			// Здесь происходит последовательное выполнение задач сканирования по указанным стратегиям. В случае, если одна
			// из стратегий вернула успех, то считаем задачу выполненной
			if res := s.scan(task); res != nil {
				successesChan <- res
				return
			}
		}
	}
}

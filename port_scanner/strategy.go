package port_scanner

// Стратегия сканирования
type scanStrategy interface {
	// Выполнение задачи сканирования и возврат ненулевого результата в случае успеха
	scan(task *scanTask) *scanTaskSuccess
}

package port_scanner

// Задача на сканирование порта
type scanTask struct {
	ip   IP
	port Port
}

// Успешно выполненная задача сканирования
type scanTaskSuccess struct {
	IP   IP
	Port Port
}

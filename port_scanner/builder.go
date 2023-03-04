package port_scanner

import (
	"math"
	"time"
)

// Количество рабочих горутин по умолчанию
const defaultNumWorkers = 200

type IP = string
type Port = uint16

// ScannerBuilder создает сканер портов
type ScannerBuilder struct {
	// Стратегии сканирования
	strategies []scanStrategy
	// IP хоста
	ip IP
	// Диапазон портов для сканирования
	portRange *PortRange
	// Количество рабочих горутин
	numWorkers int
}

func NewScannerBuilder() *ScannerBuilder {
	return &ScannerBuilder{ip: "127.0.0.1", portRange: nil, numWorkers: defaultNumWorkers}
}

// IP устанавливает ip-адрес сканируемого хоста
func (b *ScannerBuilder) IP(ip IP) *ScannerBuilder {
	b.ip = ip
	return b
}

// PortRange устанавливает диапазон портов сканируемого хоста
func (b *ScannerBuilder) PortRange(portRange PortRange) *ScannerBuilder {
	b.portRange = &portRange
	return b
}

// Tcp добавляет режим tcp-сканирования. respTimeout - максимальное время ожидания установления подключения
func (b *ScannerBuilder) Tcp(respTimeout time.Duration) *ScannerBuilder {
	b.strategies = append(b.strategies, newTcpStrategy(respTimeout))
	return b
}

// NumWorkers устанавливает максимальное количество рабочих горутин для сканирования
func (b *ScannerBuilder) NumWorkers(numWorkers int) *ScannerBuilder {
	b.numWorkers = numWorkers
	return b
}

// Build создает сканер портов
func (b *ScannerBuilder) Build() *Scanner {
	if b.portRange == nil {
		pr := OrderedRange(1, math.MaxUint16)
		b.portRange = &pr
	}
	return newScanner(b.ip, *b.portRange, b.strategies, b.numWorkers)
}

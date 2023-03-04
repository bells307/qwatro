package port_scanner

// PortRange предоставляет диапазон сканирования портов
type PortRange struct {
	ports []Port
}

// OrderedRange создает упорядоченный диапазон портов от min до max
func OrderedRange(min, max Port) PortRange {
	ports := make([]Port, max-min+1)
	for i := range ports {
		ports[i] = min + uint16(i)
	}
	return PortRange{ports}
}

// SpecificRange создает пецифический набор портов
func SpecificRange(ports ...Port) PortRange {
	return PortRange{ports}
}

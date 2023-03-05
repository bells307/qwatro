package port_scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// PortRange предоставляет диапазон сканирования портов
type PortRange struct {
	ports []Port
}

// OrderedRange создает упорядоченный диапазон портов от min до max
func OrderedRange(min, max Port) (PortRange, error) {
	if min > max {
		return PortRange{}, fmt.Errorf("invalid port range: min can't be greater than max")
	}

	ports := make([]Port, max-min+1)
	for i := range ports {
		ports[i] = min + uint16(i)
	}
	return PortRange{ports}, nil
}

// SpecificRange создает специфический набор портов
func SpecificRange(ports ...Port) PortRange {
	return PortRange{ports}
}

func RangeFromString(s string) (PortRange, error) {
	splitted := strings.Split(s, "-")
	if len(splitted) > 1 {
		min, err := strconv.ParseUint(splitted[0], 10, 64)
		if err != nil {
			return PortRange{}, fmt.Errorf("can't parse %s as uint16", splitted[0])
		}
		max, err := strconv.ParseUint(splitted[1], 10, 64)
		if err != nil {
			return PortRange{}, fmt.Errorf("can't parse %s as uint16", splitted[0])
		}
		return OrderedRange(uint16(min), uint16(max))
	} else {
		p, err := strconv.ParseUint(splitted[0], 10, 64)
		if err != nil {
			return PortRange{}, fmt.Errorf("can't parse %s as uint16", splitted[0])
		}
		return SpecificRange(uint16(p)), nil
	}
}

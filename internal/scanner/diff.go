package scanner

// Diff holds the result of comparing two snapshots.
type Diff struct {
	New    []PortInfo
	Closed []PortInfo
}

// Compare returns a Diff describing ports that appeared or disappeared
// between prev and next.
func Compare(prev, next Snapshot) Diff {
	prevMap := toMap(prev)
	nextMap := toMap(next)

	var d Diff

	for key, p := range nextMap {
		if _, exists := prevMap[key]; !exists {
			d.New = append(d.New, p)
		}
	}

	for key, p := range prevMap {
		if _, exists := nextMap[key]; !exists {
			d.Closed = append(d.Closed, p)
		}
	}

	return d
}

// toMap indexes a snapshot's ports by "proto:port" for O(1) lookup.
func toMap(s Snapshot) map[string]PortInfo {
	m := make(map[string]PortInfo, len(s.Ports))
	for _, p := range s.Ports {
		key := p.Proto + ":" + itoa(p.Port)
		m[key] = p
	}
	return m
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

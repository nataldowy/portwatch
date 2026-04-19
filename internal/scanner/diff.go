package scanner

// Diff describes changes between two snapshots.
type Diff struct {
	Opened []PortState
	Closed []PortState
}

// Compare returns the difference between a previous and current Snapshot.
func Compare(prev, curr Snapshot) Diff {
	prevMap := toMap(prev)
	currMap := toMap(curr)

	var diff Diff
	for port, state := range currMap {
		if _, existed := prevMap[port]; !existed {
			diff.Opened = append(diff.Opened, state)
		}
	}
	for port, state := range prevMap {
		if _, exists := currMap[port]; !exists {
			diff.Closed = append(diff.Closed, state)
		}
	}
	return diff
}

// HasChanges returns true when any ports opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

func toMap(s Snapshot) map[int]PortState {
	m := make(map[int]PortState, len(s.Ports))
	for _, p := range s.Ports {
		m[p.Port] = p
	}
	return m
}

package watch

import (
	"fmt"
	"sort"
	"strings"

	"portwatch/internal/scanner"
)

// Fingerprint produces a stable string hash representing a set of open ports.
// It can be used to detect whether the overall port landscape has changed
// between two scan cycles without performing a full diff.
type Fingerprint struct {
	value string
}

// NewFingerprint builds a Fingerprint from a scanner.Snapshot.
func NewFingerprint(snap scanner.Snapshot) Fingerprint {
	ports := make([]string, 0, len(snap.Ports))
	for _, p := range snap.Ports {
		ports = append(ports, fmt.Sprintf("%d/%s", p.Number, p.Proto))
	}
	sort.Strings(ports)
	return Fingerprint{value: strings.Join(ports, ",")}
}

// String returns the underlying fingerprint string.
func (f Fingerprint) String() string {
	return f.value
}

// Equal reports whether two fingerprints represent identical port sets.
func (f Fingerprint) Equal(other Fingerprint) bool {
	return f.value == other.value
}

// Empty reports whether the fingerprint represents an empty port set.
func (f Fingerprint) Empty() bool {
	return f.value == ""
}

package config

import (
	"errors"
	"fmt"
	"time"
)

// Validate checks the configuration for logical errors.
func (c *Config) Validate() error {
	if len(c.PortRanges) == 0 {
		return errors.New("config: at least one port range is required")
	}
	for i, r := range c.PortRanges {
		if err := validateRange(r); err != nil {
			return fmt.Errorf("config: port_ranges[%d]: %w", i, err)
		}
	}
	if c.Interval.Duration < time.Second {
		return errors.New("config: interval must be at least 1s")
	}
	return nil
}

func validateRange(r [2]int) error {
	start, end := r[0], r[1]
	if start < 1 || start > 65535 {
		return fmt.Errorf("start port %d out of range [1, 65535]", start)
	}
	if end < 1 || end > 65535 {
		return fmt.Errorf("end port %d out of range [1, 65535]", end)
	}
	if start > end {
		return fmt.Errorf("start port %d is greater than end port %d", start, end)
	}
	return nil
}

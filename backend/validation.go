package main

import "fmt"

var wellStatuses = map[string]bool{"producing": true, "inspection": true, "isolated": true}

func ValidateWellStatus(status string) error {
	if !wellStatuses[status] {
		return fmt.Errorf("status must be producing, inspection, or isolated")
	}
	return nil
}

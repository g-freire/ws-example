package profile

import (
	"fmt"
	"time"
)

// Profiles the function execution time.
func Start(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("\n %s took %v\n", name, time.Since(start))
	}
}
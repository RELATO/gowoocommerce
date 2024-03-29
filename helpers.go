package gowoocommerce

import (
	"fmt"
	"strconv"
)

func progressBar(completed, total int) {
	progress := float64(completed) / float64(total) * 100.0
	fmt.Print("[")
	for pct := 0.0; pct <= 100.0; pct += 4.0 {
		if pct <= progress {
			fmt.Print("#")
		} else {
			fmt.Print("-")
		}
	}
	fmt.Printf("] %s%% completed\n", strconv.FormatFloat(progress, 'f', 2, 64))
}

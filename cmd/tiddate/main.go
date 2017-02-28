package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var years int
	var months int
	var days int

	flag.IntVar(&years, "years", years, "A number to offset the year by")
	flag.IntVar(&months, "months", months, "A number to offset the month by")
	flag.IntVar(&days, "days", days, "A number to offset the day by")
	flag.Parse()

	fmt.Println(time.Now().AddDate(years, months, days).Format("2006-01-02"))
}

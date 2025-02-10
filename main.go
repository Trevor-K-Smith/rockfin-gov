package main

import (
	"fmt"
	"rockfin-gov/internal/aggregator/federal"
)

func main() {
	fmt.Println("Calling SamGovAPI...")
	federal.CallSamGovAPI()
	fmt.Println("Successfully called the SAM.gov API from client")
}

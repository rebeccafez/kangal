package main

import (
	"fmt"
	"github.com/rebeccafez/kangal/internal/config"
)

func main() {
	cfg := config.ConfigFromEnv()

	fmt.Printf("%v", cfg)
}

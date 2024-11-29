package main

import (
	"fmt"
	"github.com/jake-abed/gatorcli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	err = cfg.SetUser("jake")
	if err != nil {
		fmt.Errorf("Error; %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
	fmt.Println(cfg)
}

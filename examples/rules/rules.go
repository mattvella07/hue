package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mattvella07/hue"
)

func main() {
	h := hue.Connection{
		UserID: os.Getenv("hueUserID"),
	}

	rules, err := h.GetRules()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(rules), " rules ", rules)
}

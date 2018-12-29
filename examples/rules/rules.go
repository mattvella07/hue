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

	rule, err := h.GetRule(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Rule 1 ", rule)

	// conditions := []hue.RuleConditions{}
	// actions := []hue.RuleActions{}

	// err = h.CreateRule("NEW RULE", conditions, actions)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Rule created")
}

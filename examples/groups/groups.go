package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mattvella07/hue"
)

func main() {
	//Create connection using Hue User ID
	h := hue.Connection{
		UserID: os.Getenv("hueUserID"),
	}

	groups, err := h.GetGroups()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(groups), " groups ", groups)

	group, err := h.GetGroup(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Group 1 ", group)
}

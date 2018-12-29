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

	config, err := h.GetConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Config ", config)
}

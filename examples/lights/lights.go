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

	lights, err := h.GetLights()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(lights), " lights ", lights)

	light, err := h.GetLight(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Light 1 ", light)
}

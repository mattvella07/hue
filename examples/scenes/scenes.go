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

	scenes, err := h.GetScenes()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(scenes), " scenes ", scenes)

	for _, s := range scenes {
		fmt.Println(s.Name)
	}

	scene, err := h.GetScene("")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Scene 1", scene)

}

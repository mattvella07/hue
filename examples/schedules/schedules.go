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

	schedules, err := h.GetSchedules()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(schedules), " schedules ", schedules)

	schedule, err := h.GetSchedule(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Schedule 1 ", schedule)
}

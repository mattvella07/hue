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

	// cmd := hue.ScheduleCommand{
	// 	Address: "abc",
	// 	Body: hue.ScheduleCommandBody{
	// 		Scene: "123",
	// 	},
	// 	Method: "POST",
	// }
	// err = h.CreateSchedule("New Schedule 5", "Created by API", cmd, "2018-12-29T22:30:40", "enabled", false, false)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Schedule created")
}

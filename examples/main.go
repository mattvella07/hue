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

	//Get all lights
	lights, err := h.GetAllLights()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Found %d lights\n", len(lights))
	fmt.Println(lights)

	//Turn on light 1
	//err = h.TurnOnLight(2)
	/* err = h.TurnOnLightWithColor(4, 0.2, 0.9, 100, 200, 200)
	if err != nil {
		log.Fatalln(err)
	}

	//Sleep 3 seconds
	time.Sleep(time.Second * 3)

	//Turn off light 1
	err = h.TurnOffLight(4)
	if err != nil {
		log.Fatalln(err)
	}

	err = h.RenameLight(2, "Hue color lamp 2")
	if err != nil {
		log.Fatalln(err)
	} */

	groups, err := h.GetAllGroups()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Found %d groups\n", len(groups))
	fmt.Println(groups)
	for _, g := range groups {
		fmt.Println(fmt.Sprintf("%s: %s", g.Name, g.Type))
	}

	// err = h.CreateGroup("Blah New Group", "Room", "", []int{3, 4})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Group added!")

	// err = h.DeleteGroup(8)
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }

	// err = h.RenameGroup(7, "Living Room")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Updated room name")

	// err = h.SetLightsInGroup(7, []int{3, 4})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Updated lights in group")

	// err = h.SetGroupClass(7, "Other")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Updated class name")

	// err = h.TurnOnGroupWithColor(1, 0.2, 0.9, 100, 200, 200)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Turned on all in group 1")

	err = h.TurnOffGroup(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Turned off all lights in group 1")

	schedules, err := h.GetAllSchedules()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Schedules")
	fmt.Println(schedules)

	// cmd := hue.ScheduleCommand{
	// 	Address: "abc",
	// 	Body: hue.ScheduleCommandBody{
	// 		Scene: "123",
	// 	},
	// 	Method: "POST",
	// }
	// err = h.CreateSchedule("New Schedule", "Created by API", cmd, "2018-12-29T22:30:40", "enabled", false, false)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Schedule created")

	schedule, err := h.GetSchedule(1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Schedule 1")
	fmt.Println(schedule)
	fmt.Println(schedule.AutoDelete)

	// err = h.DeleteSchedule(6)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Schedule 6 deleted")

	scenes, err := h.GetAllScenes()
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range scenes {
		fmt.Println(s.ID, ": ", s.Name, " ", s.Lights)
	}

	// err = h.CreateLightScene("Scene 1", []int{1}, true, hue.SceneAppData{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Scene created!")

	// err = h.CreateGroupScene("Scene Group 1", 1, true, hue.SceneAppData{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Group Scene created!")

	sc, err := h.GetScene("MwnFecYhwouu-Hq")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("GetScene")
	fmt.Println(sc)

	// err = h.DeleteScene("3B1VzB2CKSuoAA-")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Scene Deleted")

	// err = h.RenameScene("XPnyhaWsrEF1VR9", "Scene 3")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Scene renamed")

	// err = h.SetLightsInScene("XPnyhaWsrEF1VR9", []int{2})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Scene lights updated")
}

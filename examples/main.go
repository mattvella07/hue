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
}

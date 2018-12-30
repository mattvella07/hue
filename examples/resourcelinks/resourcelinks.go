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

	links, err := h.GetResourceLinks()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(links), " resource links ", links)
	for _, l := range links {
		fmt.Println(l.ID)
	}

	// link, err := h.GetResourceLink(1)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Print("\n\nResource link 1 ", link)

	// err = h.CreateResourceLink("New new", "blah blah blah", true, []string{"\"/schedules/1\""})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Resource link created")

	// err = h.RenameResourceLink(1, "Routine 5")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Resource link renamed")

	// err = h.SetResourceLinkDescription(1, "Routine 5")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Resource link description updated")
}

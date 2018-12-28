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

	sensors, err := h.GetSensors()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All ", len(sensors), " sensors ", sensors)

	for _, s := range sensors {
		fmt.Println(s.Name)
	}

	// state := hue.SensorState{}
	// config := hue.SensorConfig{}
	// err = h.CreateSensor("new sensor", "123", "1.0", "Light", "abcd", "Phillips", state, config, true)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Sensor created")

	sensor, err := h.GetSensor(2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Sensor 1 ", sensor.Name, " ", sensor.Config.On)

	// err = h.TurnOnSensor(2)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println("Sensor turned on")
}

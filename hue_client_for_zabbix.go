package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/amimof/huego"
	"math"
	"os"
	"strconv"
)

type DiscoveryResult struct {
	Id   int    `json:"{#DEVICE_ID}"`
	Name string `json:"{#DEVICE_NAME}"`
}

func main() {
	hueIp := flag.String("h", "127.0.0.1", "Hue Bridge IP address")
	hueLogin := flag.String("u", "this_is_dummy_login", "Hue Bridge Username")

	flag.Parse()

	command := flag.Args()[0]

	bridge := huego.New(*hueIp, *hueLogin)

	switch command {
	case "discovery_temp_sensors":
		listSensorsByStateKey(bridge, "temperature")
	case "discovery_light_sensors":
		listSensorsByStateKey(bridge, "lightlevel")
	case "get_sensor_temp":
		sensorId, err := strconv.Atoi(flag.Args()[1])
		if err != nil {
			fmt.Printf("Parse sensor id Error: %s\n", err.Error())
			os.Exit(1)
		}
		getSensorTemp(bridge, sensorId)
	case "get_sensor_lux":
		sensorId, err := strconv.Atoi(flag.Args()[1])
		if err != nil {
			fmt.Printf("Parse sensor id Error: %s\n", err.Error())
			os.Exit(1)
		}
		getSensorLux(bridge, sensorId)
	default:
		fmt.Println("command not supported.")
		os.Exit(2)
	}
}

func listSensorsByStateKey(bridge *huego.Bridge, stateKey string) {
	sensors, err := bridge.GetSensors()
	if err != nil {
		fmt.Printf("Get sensors Error: %s\n", err.Error())
		os.Exit(1)
	}

	var result []DiscoveryResult
	for _, sensor := range sensors {
		for key, _ := range sensor.State {
			if key == stateKey {
				result = append(result, DiscoveryResult{Id: sensor.ID, Name: sensor.Name})
				break
			}
		}
	}

	outputJson, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Marshal zone state json error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s", outputJson)
}

func getSensorTemp(bridge *huego.Bridge, deviceId int) {
	sensor, err := bridge.GetSensor(deviceId)
	if err != nil {
		fmt.Printf("Get sensors Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%f", sensor.State["temperature"].(float64)/100)
}

func getSensorLux(bridge *huego.Bridge, deviceId int) {
	sensor, err := bridge.GetSensor(deviceId)
	if err != nil {
		fmt.Printf("Get sensors Error: %s\n", err.Error())
		os.Exit(1)
	}

	// https://developers.meethue.com/develop/hue-api/supported-devices/
	sensorRaw := sensor.State["lightlevel"].(float64)
	lux := math.Pow(10, (sensorRaw-1)/10000)

	fmt.Printf("%f", lux)
}

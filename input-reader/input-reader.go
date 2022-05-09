package main

import (
	"flag"
	evdev "github.com/gvalkov/golang-evdev"
	log "gitlab.com/pmoscode/golang-shared-libs/logging"
	mqttclient "gitlab.com/pmoscode/golang-shared-libs/mqtt"
	"strings"
)

var logger *log.Logger

func setupCommandLine() (*string, *string, *string, *string) {
	interfaceParam := flag.String("interface", "", "Input interface (id) to listen on")
	mqttBrokerIp := flag.String("mqtt-broker", "", "Ip of MQTT broker")
	mqttTopic := flag.String("mqtt-topic", "/input/<interfaceParam>/", "Define topic to publish message to")
	mqttClientId := flag.String("mqtt-client-id", "input-reader", "Client id for Mqtt connection")
	flag.Parse()

	return interfaceParam, mqttBrokerIp, mqttTopic, mqttClientId
}

func showInterfaces() {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return
	}

	for _, device := range devices {
		id := strings.Split(device.Fn, "/")[3]
		logger.Println("id=", id, " ## name=", device.Name)
	}
}

func cleanTag(tag string) string {
	return strings.Replace(tag, "SLASH", "-", -1)
}

func main() {
	logger = log.NewLogger("input-reader.log")
	interfaceParam, mqttBrokerIp, mqttTopic, mqttClientId := setupCommandLine()

	if *interfaceParam == "" {
		showInterfaces()
	} else {
		var rfidReader *evdev.InputDevice

		devices, err := evdev.ListInputDevices("/dev/input/" + *interfaceParam)
		if err != nil {
			logger.Fatalln(err)
		}

		if len(devices) == 0 {
			logger.Fatalln("No devices found with name: ", *interfaceParam)
		}

		rfidReader = devices[0]

		mqttClient := mqttclient.CreateClient(*mqttBrokerIp, 1883, *mqttClientId)
		mqttClient.Connect()

		logger.Println("Listening on reader...")
		rfidReader.Grab()

		defer func() {
			rfidReader.Release()
		}()

		container := make([]string, 0)
		for {
			read, err := rfidReader.ReadOne()
			if err != nil {
				return
			}

			if read.Type == evdev.EV_KEY && read.Value == 1 {
				digit := evdev.KEY[int(read.Code)]
				if digit == "KEY_ENTER" {
					tag := cleanTag(strings.Join(container, ""))
					logger.Println("Tag is: ", tag)
					msg := &mqttclient.Message{
						Topic: strings.Replace(*mqttTopic, "<interfaceParam>", *interfaceParam, -1),
						Value: tag,
					}
					mqttClient.SendMessage(msg)
					container = make([]string, 0)
				} else {
					container = append(container, strings.Split(digit, "_")[1])
				}
			}
		}
	}
}

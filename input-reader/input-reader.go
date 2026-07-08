package main

import (
	"flag"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/pmoscode/go-common/logging"
	mqttclient "github.com/pmoscode/go-common/mqtt"
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
		logger.Info("id=%s ## name=%s", id, device.Name)
	}
}

func cleanTag(tag string) string {
	return strings.Replace(tag, "SLASH", "-", -1)
}

func main() {
	logger = log.NewLogger()
	interfaceParam, mqttBrokerIp, mqttTopic, mqttClientId := setupCommandLine()

	if *interfaceParam == "" {
		showInterfaces()
	} else {
		var rfidReader *evdev.InputDevice

		devices, err := evdev.ListInputDevices("/dev/input/" + *interfaceParam)
		if err != nil {
			logger.Fatal(err, "Could not list input devices")
		}

		if len(devices) == 0 {
			logger.Fatal(nil, "No devices found with name: %s", *interfaceParam)
		}

		rfidReader = devices[0]

		mqttClient := mqttclient.NewClient(mqttclient.WithBroker(*mqttBrokerIp, 1883), mqttclient.WithClientId(*mqttClientId))
		if err := mqttClient.Connect(); err != nil {
			logger.Fatal(err, "Could not connect to mqtt broker")
		}

		logger.Info("Listening on reader...")
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
					logger.Info("Tag is: %s", tag)
					msg := &mqttclient.Message{
						Topic: strings.Replace(*mqttTopic, "<interfaceParam>", *interfaceParam, -1),
						Value: tag,
					}
					mqttClient.Publish(msg)
					container = make([]string, 0)
				} else {
					container = append(container, strings.Split(digit, "_")[1])
				}
			}
		}
	}
}

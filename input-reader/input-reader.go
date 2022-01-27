package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	evdev "github.com/gvalkov/golang-evdev"
	"log"
	"strings"
)

func setupCommandLine() (*string, *string, *string) {
	interfaceParam := flag.String("interface", "", "Input interface (id) to liste on")
	mqttBrokerIp := flag.String("mqtt-broker", "", "Ip of MQTT broker")
	mqttTopic := flag.String("mqtt-topic", "/input/<interfaceParam>/", "Define topic to publish message to")
	flag.Parse()

	return interfaceParam, mqttBrokerIp, mqttTopic
}

func showInterfaces() {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return
	}

	for _, device := range devices {
		id := strings.Split(device.Fn, "/")[3]
		fmt.Println("id=", id, " ## name=", device.Name)
	}
}

func main() {
	interfaceParam, mqttBrokerIp, mqttTopic := setupCommandLine()

	if *interfaceParam == "" {
		showInterfaces()
	} else {
		var rfidReader *evdev.InputDevice

		devices, err := evdev.ListInputDevices("/dev/input/" + *interfaceParam)
		if err != nil {
			log.Fatalln(err)
		}

		if len(devices) == 0 {
			log.Fatalln("No devices found with name: ", *interfaceParam)
		}

		rfidReader = devices[0]

		mqttClient := connectToMqtt(mqttBrokerIp)

		fmt.Println("Listening on reader...")
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
					tag := strings.Join(container, "")
					fmt.Println("Tag is: ", tag)
					sendMessage(mqttClient, strings.Replace(*mqttTopic, "<interfaceParam>", *interfaceParam, -1), tag)
					container = make([]string, 0)
				} else {
					container = append(container, strings.Split(digit, "_")[1])
				}
			}
		}
	}
}

func connectToMqtt(mqttBrokerIp *string) *mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", *mqttBrokerIp, 1883))
	opts.SetClientID("go_mqtt_client")
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Could not connect with broker: ", token.Error())
		log.Println("MQTT is disabled now")

		return nil
	}

	return &client
}

func sendMessage(mqttClient *mqtt.Client, mqttTopic string, tag string) {
	if mqttClient != nil {
		client := *mqttClient
		token := client.Publish(mqttTopic, 2, false, tag)
		token.Wait()
	}
}

package main

import (
	"flag"
	evdev "github.com/gvalkov/golang-evdev"
	mqttclient "github.com/pmoscode/golang-mqtt"
	"log"
	writer2 "serial-writer/writer"
	"strings"
)

var debugWriter = false

func setupCommandLine() (*string, *string, *string, *string, *bool) {
	interfaceParam := flag.String("interface", "", "Output interface (id) where data is send")
	mqttBrokerIp := flag.String("mqtt-broker", "", "Ip of MQTT broker")
	mqttTopic := flag.String("mqtt-topic", "/output/<interfaceParam>/", "Define topic to subscribe on")
	mqttClientId := flag.String("mqtt-client-id", "serial-writer", "Client id for Mqtt connection")
	debug := flag.Bool("debug", false, "Check if data should be send to console for debugging")
	flag.Parse()

	return interfaceParam, mqttBrokerIp, mqttTopic, mqttClientId, debug
}

func showInterfaces() {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return
	}

	for _, device := range devices {
		id := strings.Split(device.Fn, "/")[3]
		log.Println("id=", id, " ## name=", device.Name)
	}
}

func onMessageReceived(message mqttclient.Message) {
	var writer writer2.Writer

	if debugWriter {
		writer = &writer2.Debug{}
	} else {
		writer = &writer2.Console{}
	}

	writer.Write(message.ToString())
}

func main() {
	interfaceParam, mqttBrokerIp, mqttTopic, mqttClientId, debug := setupCommandLine()
	debugWriter = *debug

	if *interfaceParam == "" {
		showInterfaces()
	} else {
		mqttClient := mqttclient.CreateClient(*mqttBrokerIp, 1883, *mqttClientId)
		mqttClient.Connect()
		mqttClient.Subscribe(*mqttTopic, onMessageReceived)
	}
}

package main

import (
	"flag"
	evdev "github.com/gvalkov/golang-evdev"
	mqttclient "github.com/pmoscode/golang-mqtt"
	"log"
	writer2 "serial-writer/writer"
	"strings"
)

var writer writer2.Writer

func setupCommandLine() (*string, *int, *string, *string, *string, *bool) {
	interfaceParam := flag.String("interface", "", "Output interface (id) where data is send")
	interfaceBadRateParam := flag.Int("baudRate", 19200, "Baud rate of the serial console")
	mqttBrokerIp := flag.String("mqtt-broker", "", "Ip of MQTT broker")
	mqttTopic := flag.String("mqtt-topic", "/output/<interfaceParam>/", "Define topic to subscribe on")
	mqttClientId := flag.String("mqtt-client-id", "serial-writer", "Client id for Mqtt connection")
	debug := flag.Bool("debug", false, "Check if data should be send to console for debugging")
	flag.Parse()

	return interfaceParam, interfaceBadRateParam, mqttBrokerIp, mqttTopic, mqttClientId, debug
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
	writer.Write(message.ToString())
}

func connectWriter(interfaceName string, interfaceBadRate int, debugWriter bool) {
	if debugWriter {
		writer = &writer2.Debug{}
	} else {
		writer = &writer2.SerialConsole{
			InterfaceName:     interfaceName,
			InterfaceBaudRate: interfaceBadRate,
		}
	}

	writer.Connect()
}

func main() {
	interfaceParam, interfaceBadRateParam, mqttBrokerIp, mqttTopic, mqttClientId, debug := setupCommandLine()

	if *interfaceParam == "" {
		showInterfaces()
	} else {
		connectWriter(*interfaceParam, *interfaceBadRateParam, *debug)

		mqttClient := mqttclient.CreateClient(*mqttBrokerIp, 1883, *mqttClientId)
		mqttClient.Connect()
		mqttClient.Subscribe(*mqttTopic, onMessageReceived)
	}
}

package main

import (
	"flag"
	mqttclient "github.com/pmoscode/golang-mqtt/mqtt"
	"log"
	writer2 "serial-writer/writer"
)

var writer writer2.Writer

func setupCommandLine() (*string, *int, *string, *string, *string, *string, *bool) {
	interfaceParam := flag.String("interface", "", "Output interface (id) where data is send")
	interfaceBadRateParam := flag.Int("baudRate", 19200, "Baud rate of the serial console")
	writerParam := flag.String("writer", "shell", "writer implementation to use: shell or port (default: shell)")
	mqttBrokerIp := flag.String("mqtt-broker", "", "Ip of MQTT broker")
	mqttTopic := flag.String("mqtt-topic", "/output/<interfaceParam>/", "Define topic to subscribe on")
	mqttClientId := flag.String("mqtt-client-id", "serial-writer", "Client id for Mqtt connection")
	debug := flag.Bool("debug", false, "Check if data should be send to console for debugging")
	flag.Parse()

	return interfaceParam, interfaceBadRateParam, writerParam, mqttBrokerIp, mqttTopic, mqttClientId, debug
}

func onMessageReceived(message mqttclient.Message) {
	writer.Write(message.ToString())
}

func connectWriter(interfaceName string, interfaceBadRate int, writerImpl string, debugWriter bool) {
	if debugWriter {
		writer = &writer2.Debug{}
	} else {
		switch writerImpl {
		case "shell":
			writer = &writer2.Shell{
				InterfaceName: interfaceName,
			}
		case "port":
			writer = &writer2.Port{
				InterfaceName:     interfaceName,
				InterfaceBaudRate: interfaceBadRate,
			}
		default:
			log.Fatalln("Unknown writer implementation: ", writerImpl)
		}
	}

	writer.Connect()
}

func main() {
	interfaceParam, interfaceBadRateParam, writerParam, mqttBrokerIp, mqttTopic, mqttClientId, debug := setupCommandLine()

	if *interfaceParam == "" {
		log.Fatalln("'interface' parameter is required!")
	} else {
		connectWriter(*interfaceParam, *interfaceBadRateParam, *writerParam, *debug)

		mqttClient := mqttclient.CreateClient(*mqttBrokerIp, 1883, *mqttClientId)
		mqttClient.Connect()
		mqttClient.Subscribe(*mqttTopic, onMessageReceived)
	}
}

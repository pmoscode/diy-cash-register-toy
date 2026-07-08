package main

import (
	"flag"

	log "github.com/pmoscode/go-common/logging"
	mqttclient "github.com/pmoscode/go-common/mqtt"
	"serial-writer/hardware"
	writer2 "serial-writer/writer"
)

var writer writer2.Writer
var logger *log.Logger

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
	writer.Write(message.ToRawString())
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
			logger.Fatal(nil, "Unknown writer implementation: %s", writerImpl)
		}
	}

	writer.Connect()
}

func main() {
	logger = log.NewLogger()
	interfaceParam, interfaceBadRateParam, writerParam, mqttBrokerIp, mqttTopic, mqttClientId, debug := setupCommandLine()

	if *interfaceParam == "" {
		logger.Fatal(nil, "'interface' parameter is required!")
	} else {
		isPi := hardware.IsRaspberryPiHardware()
		if !isPi {
			*debug = true
		}
		connectWriter(*interfaceParam, *interfaceBadRateParam, *writerParam, *debug)

		mqttClient := mqttclient.NewClient(mqttclient.WithBroker(*mqttBrokerIp, 1883), mqttclient.WithClientId(*mqttClientId))
		if err := mqttClient.Connect(); err != nil {
			logger.Fatal(err, "Could not connect to mqtt broker")
		}
		mqttClient.Subscribe(*mqttTopic, onMessageReceived)
		mqttClient.LoopForever()
	}
}

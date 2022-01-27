package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tarm/serial"
	"log"
)

func main() {
	// Connect to broker

	// setup serial

	// run Sub in a loop
	// send data to printer when topic has got data

	c := &serial.Config{Name: "COM5", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
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

func sub(client *mqtt.Client) {
	clientO := *client
	topic := "topic/test"
	token := clientO.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

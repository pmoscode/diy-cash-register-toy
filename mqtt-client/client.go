package mqtt_client

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	Topic string
	Value interface{}
}

func (m Message) ToJson() string {
	marshal, err := json.Marshal(m.Value)
	if err != nil {
		return fmt.Sprintf("%v", m.Value)
	}
	return string(marshal)
}

func (m Message) ToString() string {
	return string(m.Value.([]uint8))
}

type client struct {
	brokerIp string
	port     int
	topic    string
	clientId string
	client   *mqtt.Client
}

func (c *client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.brokerIp, c.port))
	opts.SetClientID(c.clientId)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Could not connect to broker: ", token.Error())

		return token.Error()
	}

	c.client = &client
	log.Println("Mqtt connected to", c.brokerIp)

	return nil
}

func (c *client) Disconnect() {
	client := *c.client
	client.Disconnect(100)
}

func (c *client) SendMessage(message *Message) {
	if c.client == nil {
		log.Println("Mqtt client not connected! Call 'connect' method...")
	} else {
		client := *c.client
		token := client.Publish(message.Topic, 2, false, message.ToJson())
		token.Wait()
	}
}

func (c *client) Subscribe(topic string, fn func(message Message)) {
	if c.client == nil {
		log.Println("Mqtt client not connected! Call 'connect' method...")
	} else {
		client := *c.client
		client.Subscribe(topic, 2, func(client mqtt.Client, msg mqtt.Message) {
			message := Message{
				Topic: msg.Topic(),
				Value: msg.Payload(),
			}
			fn(message)
		})

		channel := make(chan os.Signal)
		signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-channel
			c.Disconnect()
			log.Println("Exiting...")
			os.Exit(1)
		}()

		for {
			time.Sleep(1 * time.Second)
		}
	}
}

func CreateClient(brokerIp string, port int, clientId string) *client {
	return &client{
		brokerIp: brokerIp,
		port:     port,
		clientId: clientId,
	}
}

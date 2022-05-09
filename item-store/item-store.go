package main

import (
	"flag"
	log "gitlab.com/pmoscode/golang-shared-libs/logging"
	mqttclient "gitlab.com/pmoscode/golang-shared-libs/mqtt"
	"item-store/product"
)

var productsFilename = "products.yaml"
var mqttTopicPublish string
var productList *product.List
var mqttClient *mqttclient.Client

var logger *log.Logger

func setupCommandLine() (*string, *string, *string, *string) {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "item-store", "Client id for Mqtt connection")
	mqttTopicSub := flag.String("mqtt-topic-sub", "/input/item/", "Define topic to subscribe resolver requests to")
	mqttTopicPub := flag.String("mqtt-topic-pub", "/output/item/", "Define topic to publish resolved requests to")
	flag.Parse()

	return mqttBrokerIp, mqttClientId, mqttTopicSub, mqttTopicPub
}

func onMessage(message mqttclient.Message) {
	code := message.ToString()

	productItem := productList.FromCode(code)
	if productItem != nil {
		mqttMessage := mqttclient.Message{
			Topic: mqttTopicPublish,
			Value: productItem,
		}
		mqttClient.SendMessage(&mqttMessage)
	} else {
		mqttMessage := mqttclient.Message{
			Topic: mqttTopicPublish,
			Value: "not found",
		}
		mqttClient.SendMessage(&mqttMessage)
	}
}

func main() {
	logger = log.NewLogger("item-store.log")
	mqttBrokerIp, mqttClientId, mqttTopicSub, mqttTopicPub := setupCommandLine()
	mqttTopicPublish = *mqttTopicPub

	list, err := product.NewProductList(productsFilename)
	if err != nil {
		logger.Println(err)
	}

	productList = list

	mqttClient = mqttclient.CreateClient(*mqttBrokerIp, 1883, *mqttClientId)
	mqttClient.Connect()

	mqttClient.Subscribe(*mqttTopicSub, onMessage)
}

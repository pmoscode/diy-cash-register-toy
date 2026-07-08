package main

import (
	"flag"

	log "github.com/pmoscode/go-common/logging"
	mqttclient "github.com/pmoscode/go-common/mqtt"
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
	code := message.ToRawString()

	productItem := productList.FromCode(code)
	if productItem != nil {
		mqttMessage := mqttclient.Message{
			Topic: mqttTopicPublish,
			Value: productItem,
		}
		mqttClient.Publish(&mqttMessage)
	} else {
		mqttMessage := mqttclient.Message{
			Topic: mqttTopicPublish,
			Value: "not found",
		}
		mqttClient.Publish(&mqttMessage)
	}
}

func main() {
	logger = log.NewLogger()
	mqttBrokerIp, mqttClientId, mqttTopicSub, mqttTopicPub := setupCommandLine()
	mqttTopicPublish = *mqttTopicPub

	list, err := product.NewProductList(productsFilename)
	if err != nil {
		logger.Error(err, "Could not load product list")
	}

	productList = list

	mqttClient = mqttclient.NewClient(mqttclient.WithBroker(*mqttBrokerIp, 1883), mqttclient.WithClientId(*mqttClientId))
	if err := mqttClient.Connect(); err != nil {
		logger.Fatal(err, "Could not connect to mqtt broker")
	}

	mqttClient.Subscribe(*mqttTopicSub, onMessage)
	mqttClient.LoopForever()
}

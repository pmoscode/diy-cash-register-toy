package main

import (
	"fmt"
	"github.com/256dpi/lungo"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Item struct {
	Code  string  `yaml:"code"`
	Name  string  `yaml:"name"`
	Value float32 `yaml:"value"`
}

type Products struct {
	Items []Item `yaml:"items"`
}

var productsFilename = "products.yaml"

func main() {
	var products *Products

	_, error := os.Stat(productsFilename)

	if !os.IsNotExist(error) {
		file, err := os.Open(productsFilename)
		if err != nil {
			log.Println(err.Error())
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(products); err != nil {
			log.Println(err.Error())
		}

		fmt.Println(products)
		fmt.Println("First item: ", products.Items[0].Name)
		fmt.Println("Second item: ", products.Items[1].Name)
	}

	// prepare options
	opts := lungo.Options{
		Store: lungo.NewFileStore("data/products.db", os.ModePerm),
	}

	// open database
	client, engine, err := lungo.Open(nil, opts)
	if err != nil {
		panic(err)
	}

	defer engine.Close()

	foo := client.Database("products")

	bar := foo.Collection("items")

	if products != nil {
		for _, item := range products.Items {
			_, err = bar.InsertOne(nil, &item)
			if err != nil {
				panic(err)
			}
		}
	}

	csr, err := bar.Find(nil, bson.D{})
	if err != nil {
		panic(err)
	}

	var items []Item
	err = csr.All(nil, &items)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", items)
	fmt.Println()

	err = os.Rename(productsFilename, productsFilename+".inserted")
	if err != nil {
		return
	}
}

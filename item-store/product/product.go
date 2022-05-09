package product

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Product struct {
	Code  string  `yaml:"code"`
	Name  string  `yaml:"name"`
	Value float32 `yaml:"value"`
}

type List struct {
	Products []Product `yaml:"products"`
}

func (p *List) FromCode(code string) *Product {
	for _, item := range p.Products {
		if item.Code == code {
			return &item
		}
	}

	return nil
}

func (p Product) String() string {
	return fmt.Sprintf("{ Code: %s ## Name: %s ## Value: %f }", p.Code, p.Name, p.Value)
}

func NewProductList(path string) (*List, error) {
	var productList = &List{}
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		getwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		absPath := filepath.Join(getwd, path)

		fmt.Println("Loading products from: ", absPath)
		file, err := os.Open(absPath)
		if err != nil {
			log.Println(err.Error())
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(productList); err != nil {
			log.Println(err.Error())
		}

		return productList, nil
	} else {
		return nil, errors.New("no file found with path: " + path)
	}
}

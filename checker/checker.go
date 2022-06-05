package checker

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/andrysds/dropship-checker/csv"
	"github.com/andrysds/dropship-checker/product"
)

const (
	stockLevelKeyEnvKey  = "STOCK_LEVEL_KEY"
	priceKeyEnvKey       = "PRICE_KEY"
	productSlugKeyEnvKey = "PRODUCT_SLUG_KEY"
)

type Partner interface {
	Login() error
	GetProduct(slug string) (*product.Product, error)
}

type Checker struct {
	records        []csv.Record
	partner        Partner
	stockLevelKey  string
	priceKey       string
	productSlugKey string
}

func NewChecker(records []csv.Record, partner Partner) *Checker {
	return &Checker{
		records:        records,
		partner:        partner,
		stockLevelKey:  os.Getenv(stockLevelKeyEnvKey),
		priceKey:       os.Getenv(priceKeyEnvKey),
		productSlugKey: os.Getenv(productSlugKeyEnvKey),
	}
}

func (c *Checker) Check() error {
	if err := c.partner.Login(); err != nil {
		return err
	}

	for i, record := range c.records {
		data := record.Data

		product, err := c.partner.GetProduct(data[c.productSlugKey])
		if err != nil {
			log.Println("[ERROR] [GetProduct]", err)
		}

		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Println("[ERROR] [parsing product json]", err)
		}

		variants := product.Variants
		for _, variant := range variants {
			oldPrice, err := strconv.ParseInt(data[c.priceKey], 10, 32)
			if err != nil {
				log.Println("[ERROR] [parsing old price]", err)
			}

			if variant.IsPriceChanged(int(oldPrice)) {
				log.Printf("[WARN] price change detected at row no. %d; product: %v\n", i+1, string(productJSON))
			}

			oldStockLevel, err := strconv.ParseInt(data[c.stockLevelKey], 10, 32)
			if err != nil {
				log.Println("[ERROR] [parsing old stock level]", err)
			}

			if variant.IsStockLevelChange(int(oldStockLevel)) {
				log.Printf("[WARN] stock level change detected at row no. %d; product: %v\n", i+1, string(productJSON))
			}
		}
	}

	return nil
}

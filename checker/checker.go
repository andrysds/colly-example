package checker

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andrysds/dropship-checker/csv"
	"github.com/andrysds/dropship-checker/product"
)

const (
	stockLevelKeyEnvKey  = "STOCK_LEVEL_KEY"
	priceKeyEnvKey       = "PRICE_KEY"
	productSlugKeyEnvKey = "PRODUCT_SLUG_KEY"
	variantNameKeyEnvKey = "VARIANT_NAME_KEY"
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
	variantKey     string
}

func NewChecker(records []csv.Record, partner Partner) *Checker {
	return &Checker{
		records:        records,
		partner:        partner,
		stockLevelKey:  os.Getenv(stockLevelKeyEnvKey),
		priceKey:       os.Getenv(priceKeyEnvKey),
		productSlugKey: os.Getenv(productSlugKeyEnvKey),
		variantKey:     os.Getenv(variantNameKeyEnvKey),
	}
}

func (c *Checker) Check() error {
	if err := c.partner.Login(); err != nil {
		return err
	}

	for i, record := range c.records {
		data := record.Data
		slug := data[c.productSlugKey]

		if slug == "" {
			log.Println("found empty slug at row no.", i+1)
			break
		}

		product, err := c.partner.GetProduct(slug)
		if err != nil {
			log.Println("[ERROR] [GetProduct]", err)
		}

		found := false
		for _, variant := range product.Variants {
			if variant.Name == record.Data[c.variantKey] {
				found = true

				oldPriceStr := data[c.priceKey]
				oldPriceStr = strings.ReplaceAll(oldPriceStr, "Rp", "")
				oldPriceStr = strings.ReplaceAll(oldPriceStr, ",", "")
				oldPrice, err := strconv.ParseInt(oldPriceStr, 10, 32)
				if err != nil {
					log.Println("[ERROR] [parsing old price]", err)
				}

				if variant.IsPriceChanged(int(oldPrice)) {
					log.Printf("[WARN] price change detected at row no. %d; new price: %d; product: %v\n", i+1, variant.Price, slug)
				}

				oldStockLevel, err := strconv.ParseInt(data[c.stockLevelKey], 10, 32)
				if err != nil {
					log.Println("[ERROR] [parsing old stock level]", err)
				}

				if variant.IsStockLevelChange(int(oldStockLevel)) {
					log.Printf("[WARN] stock level change detected at row no. %d; new stock level: %d; slug: %v\n", i+1, variant.StockLevel(), slug)
				}

				break
			}
		}

		if !found {
			log.Printf("[ERROR] product: %v,  variant: %v, not found", product.Name, data[c.variantKey])
		}
	}

	return nil
}

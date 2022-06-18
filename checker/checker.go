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
	skuKeyEnvKey         = "SKU_KEY"
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
	skuKey         string
}

func NewChecker(records []csv.Record, partner Partner) *Checker {
	return &Checker{
		records:        records,
		partner:        partner,
		stockLevelKey:  os.Getenv(stockLevelKeyEnvKey),
		priceKey:       os.Getenv(priceKeyEnvKey),
		productSlugKey: os.Getenv(productSlugKeyEnvKey),
		variantKey:     os.Getenv(variantNameKeyEnvKey),
		skuKey:         os.Getenv(skuKeyEnvKey),
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
				sku := record.Data[c.skuKey]
				oldPriceStr := data[c.priceKey]
				oldPriceStr = strings.ReplaceAll(oldPriceStr, "Rp", "")
				oldPriceStr = strings.ReplaceAll(oldPriceStr, ",", "")
				oldPrice, err := strconv.ParseInt(oldPriceStr, 10, 32)
				if err != nil {
					log.Println("[ERROR] [parsing old price]", err)
				}

				if variant.IsPriceChanged(int(oldPrice)) {
					log.Printf("[WARN] price change detected; row: %d; new price: %d; sku: %s\n", i+1, variant.Price, sku)
				}

				oldStockLevel, err := strconv.ParseInt(data[c.stockLevelKey], 10, 32)
				if err != nil {
					log.Println("[ERROR] [parsing old stock level]", err)
				}

				if variant.IsStockLevelChange(int(oldStockLevel)) {
					log.Printf("[WARN] stock level change detected; row: %d; new stock level: %d; sku: %s\n", i+1, variant.StockLevel(), sku)
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

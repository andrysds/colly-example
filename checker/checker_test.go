package checker

import (
	"os"
	"reflect"
	"testing"

	"github.com/andrysds/dropship-checker/csv"
	"github.com/andrysds/dropship-checker/product"
)

func TestNewChecker(t *testing.T) {
	mockStockLevelKey := "header1"
	mockPriceKey := "header2"
	mockProductSlugKey := "header3"
	mockVariantKey := "header4"
	mockPartner := &MockPartner{}

	mockRecords := []csv.Record{
		{
			Data: map[string]string{
				"header1": "data1",
				"header2": "data2",
				"header3": "data3",
			},
		},
	}

	os.Setenv(stockLevelKeyEnvKey, mockStockLevelKey)
	os.Setenv(priceKeyEnvKey, mockPriceKey)
	os.Setenv(productSlugKeyEnvKey, mockProductSlugKey)
	os.Setenv(variantNameKeyEnvKey, mockVariantKey)
	defer os.Unsetenv(stockLevelKeyEnvKey)
	defer os.Unsetenv(mockPriceKey)
	defer os.Unsetenv(productSlugKeyEnvKey)
	defer os.Unsetenv(variantNameKeyEnvKey)

	want := &Checker{
		records:        mockRecords,
		partner:        mockPartner,
		stockLevelKey:  mockStockLevelKey,
		priceKey:       mockPriceKey,
		productSlugKey: mockProductSlugKey,
		variantKey:     mockVariantKey,
	}

	if got := NewChecker(mockRecords, mockPartner); !reflect.DeepEqual(got, want) {
		t.Errorf("NewChecker() = %v, want %v", got, want)
	}
}

func TestChecker_Check(t *testing.T) {
	mockStockLevelKey := "header1"
	mockPriceKey := "header2"
	mockProductSlugKey := "header3"
	mockSlug := "sample-slug"

	mockRecords := []csv.Record{
		{
			Data: map[string]string{
				"header1": "2",
				"header2": "1000",
				"header3": mockSlug,
			},
		},
	}

	mockProduct := &product.Product{
		Name:        "sample name",
		Description: "description",
		Variants: []product.Variant{
			{
				Name:  "sample name",
				Price: 2000,
				Stock: 0,
			},
		},
	}

	mockPartner := &MockPartner{}
	mockPartner.On("Login").Return(nil)
	mockPartner.On("GetProduct", mockSlug).Return(mockProduct, nil)

	c := &Checker{
		records:        mockRecords,
		partner:        mockPartner,
		stockLevelKey:  mockStockLevelKey,
		priceKey:       mockPriceKey,
		productSlugKey: mockProductSlugKey,
	}

	wantErr := false
	if err := c.Check(); (err != nil) != wantErr {
		t.Errorf("Checker.Check() error = %v, wantErr %v", err, wantErr)
	}
}

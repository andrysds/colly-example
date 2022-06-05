package product

import (
	"reflect"
	"testing"
)

func TestProduct_VariantMap(t *testing.T) {
	v := Variant{
		Name:  "sample name",
		Price: 1000,
		Stock: 0,
	}

	p := &Product{
		Name:        "sample name",
		Description: "sample description",
		Variants:    []Variant{v},
	}

	want := map[string]Variant{"sample name": v}

	if got := p.VariantMap(); !reflect.DeepEqual(got, want) {
		t.Errorf("product.variantMap() = %v, want %v", got, want)
	}
}

func TestVariant_IsPriceChanged(t *testing.T) {
	tests := []struct {
		name         string
		currentPrice int
		oldPrice     int
		want         bool
	}{
		{
			name:         "price is not changed",
			currentPrice: 1000,
			oldPrice:     1000,
			want:         false,
		},
		{
			name:         "price is changed",
			currentPrice: 2000,
			oldPrice:     1000,
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Variant{Price: tt.currentPrice}
			if got := v.IsPriceChanged(tt.oldPrice); got != tt.want {
				t.Errorf("variant.isPriceChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariant_IsStockLevelChange(t *testing.T) {
	tests := []struct {
		name          string
		currentStock  int
		oldStockLevel int
		want          bool
	}{
		{
			name:          "stock level is not changed",
			currentStock:  0,
			oldStockLevel: 0,
			want:          false,
		},
		{
			name:          "stock level is changed",
			currentStock:  1,
			oldStockLevel: 0,
			want:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Variant{Stock: tt.currentStock}
			if got := v.IsStockLevelChange(tt.oldStockLevel); got != tt.want {
				t.Errorf("variant.isStockLevelChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariant_stockLevel(t *testing.T) {
	tests := []struct {
		name  string
		stock int
		want  int
	}{
		{
			name:  "out of stock",
			stock: 0,
			want:  OutOfStock,
		},
		{
			name:  "negative stock",
			stock: -100,
			want:  OutOfStock,
		},
		{
			name:  "low stock",
			stock: 5,
			want:  LowStock,
		},
		{
			name:  "High stock",
			stock: 100,
			want:  HighStock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Variant{Stock: tt.stock}
			if got := v.stockLevel(); got != tt.want {
				t.Errorf("variant.stockLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

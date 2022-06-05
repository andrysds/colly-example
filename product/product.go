package product

type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
}

func (p *Product) VariantMap() map[string]Variant {
	res := map[string]Variant{}
	for _, v := range p.Variants {
		res[v.Name] = v
	}
	return res
}

type Variant struct {
	Name  string `json:"variants_name"`
	Price int    `json:"price"`
	Stock int    `json:"stok"`
}

func (v *Variant) IsPriceChanged(oldPrice int) bool {
	return oldPrice != v.Price
}

func (v Variant) IsStockLevelChange(oldStockLevel int) bool {
	return oldStockLevel != v.stockLevel()
}

const (
	OutOfStock = iota
	LowStock   = iota
	HighStock  = iota
)

func (v *Variant) stockLevel() int {
	if v.Stock < 1 {
		return OutOfStock
	}

	if v.Stock < 10 {
		return LowStock
	}

	return HighStock
}

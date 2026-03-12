package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"` //dont forget to add unique val. tag
	CreatedOn   string  `json:"created_on"`
	UpdatedOn   string  `json:"updated_on"`
	DeletedOn   string  `json:"deleted_on"`
}

type Products []*Product 

func GetAllProduct() Products{
	return productList
}

func (p *Products) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// add
func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// update
func UpdateProduct(id int, p *Product) error{
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	// pointer basically not really need smth, so just return error / nil
	p.ID = id
	productList[pos] = p

	return nil
}

// validate
func (p *Product) Validate() error {
	validate := validator.New()

	// can also register new one :
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// format : abs-abdf-absdif
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}
	return true
}

// reusable 
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, fmt.Errorf("Product %d not found", id)
}


var productList = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
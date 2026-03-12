package data

import "testing"

func TestCheckValidation(t *testing.T){
	p := &Product{
		Name: "test",
		Price: 1,
		SKU: "test-test-test",
	}

	err := p.Validate()
	if err != nil{
		t.Fatal(err)
	}
}
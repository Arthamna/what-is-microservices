package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log hclog.Logger
	rates map[string]float64 
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: make(map[string]float64)}
	return er, nil
}

// if succeed, return itself (sometimes i forgot)
func (e *ExchangeRates) getRates() error{

	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64) // str to float
		if err != nil{
			return err
		}
		e.rates[c.Currency] = r
	}

	return nil
}

type Cubes struct{
	CubeData []Cube `xml:"Cube>Cube>Cube"` // this one depends on your xml format file
}

type Cube struct{
	Currency string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}


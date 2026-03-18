package data

import (
	"encoding/xml"
	"fmt"
	"log"
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
	err := er.getRates()
	return er, err
}

func (e *ExchangeRates) GetRate(base string, dest string) (float64, error){
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("currency %s not found", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("currency %s not found", dest)
	}

	// conv := destination / base
	return dr / br, nil
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
	// case when converting to eur, given there isn't any in the xml
	e.rates["EUR"] = 1	
	
	//
	log.Printf("loaded %d exchange rates", len(e.rates))
	return nil
}

type Cubes struct{
	CubeData []Cube `xml:"Cube>Cube>Cube"` // this one depends on your xml format file
}

type Cube struct{
	Currency string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}


package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/nicholasjackson/building-microservices-youtube/currency/protos/currency"

	"github.com/nicholasjackson/building-microservices-youtube/currency/data"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	protos.UnimplementedCurrencyServer
	rates *data.ExchangeRates
	log hclog.Logger
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{
		rates : r,
		log : l,
	}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	// unicode conv. vs string conv.
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error{
	
	// two blocking func, need to run both
	go func() {
		for {
			rr, err := src.Recv()
			if err == io.EOF{
				c.log.Info("Client closed connection")
				break
			}
			
			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				break
			}

			c.log.Info("Handle request", "request", rr )
		}
	}()

	for {
		err := src.Send(&protos.RateResponse{Rate: 10.0})
		if err != nil {
			c.log.Error("Unable to send to client", "error", err)
			break
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

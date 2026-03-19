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
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	
	c := &Currency{
		rates : r,
		log : l,
		subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest),
	}

	go c.handleUpdates()
	return c
}

// debug global func ex
// {
// 	"Base": "USD",
// 	"Destination": "EUR",
// }


func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Rates updated")

		// loop over subs clients
		for k,v := range c.subscriptions {
			
			// loop over subs rates
			for _, rr := range v {
				rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

			if err != nil {
				c.log.Error("Unable to get rate", "base", rr.GetBase(), "dest", rr.GetDestination(), "error", err)
			}
			
			err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate})

			if err != nil {
				c.log.Error("Unable to get rate", "base", rr.GetBase(), "dest", rr.GetDestination(), "error", err) }
			}

		}
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

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error{

	// two blocking func, need to run both
	// handle client
	for {
		rr, err := src.Recv()
		if err == io.EOF{
			c.log.Info("Client closed connection")
			break
		}
		
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			return err
		}

		c.log.Info("Handle request", "request", rr)

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = make([]*protos.RateRequest, 0)
		}
		
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs	

		// c.subscriptions[src] = append(c.subscriptions[src], rr)
	}
	
	// // handle server responses
	// for {
	// 	err := src.Send(&protos.RateResponse{Rate: 10.0})
	// 	if err != nil {
	// 		c.log.Error("Unable to send to client", "error", err)
	// 		break
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }
	return nil
}

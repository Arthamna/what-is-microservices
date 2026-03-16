package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	protos "github.com/nicholasjackson/building-microservices-youtube/currency/protos/currency"
)

func unaryLogInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	log.Printf("GRPC CALL -> method: %s, req: %+v", method, req)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func main() {
	conn, err := grpc.Dial(
		"localhost:9092",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryLogInterceptor),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("gRPC client connected to localhost:9092")

	// disini biasanya dibuat client dari proto
	// example:
	client := protos.NewCurrencyClient(conn)
	rr := &protos.RateRequest{
		Base: protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["GBP"]),
	}

	resp, err := client.GetRate(context.Background(), rr)
	if err != nil {
		log.Fatalf("failed to call GetRate: %v", err)
	}

	log.Printf("response: %+v", resp)
}
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func brasilApiRequest(ctx context.Context, brasilApiChan chan<- io.ReadCloser, cep string) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	brasilApiChan <- response.Body
}

func viaCepApi(ctx context.Context, viaCepChan chan<- io.ReadCloser, cep string) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	viaCepChan <- response.Body
}

func main() {
	brasilApiChan := make(chan io.ReadCloser)
	viaCepApiChan := make(chan io.ReadCloser)

	cep := os.Args[1]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go brasilApiRequest(ctx, brasilApiChan, cep)
	go viaCepApi(ctx, viaCepApiChan, cep)

	select {
	case brasilApi := <-brasilApiChan:
		defer brasilApi.Close()
		log.Println("Brasil API")
		io.Copy(os.Stdout, brasilApi)
	case viaCep := <-viaCepApiChan:
		defer viaCep.Close()
		log.Println("Via Cep API")
		io.Copy(os.Stdout, viaCep)
	case <-time.After(time.Second):
		log.Println("Timeout")
	}
}

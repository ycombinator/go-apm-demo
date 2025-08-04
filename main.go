package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"golang.org/x/net/context/ctxhttp"
)

func main() {
	printAPMEnv()
	var tracingClient = apmhttp.WrapClient(http.DefaultClient)
	tracer := apm.DefaultTracer()
	defer tracer.Flush(nil)

	tx := tracer.StartTransaction("GET https://google.com", "request")
	defer tx.End()

	ctx := context.Background()
	resp, err := ctxhttp.Get(ctx, tracingClient, "https://google.com")
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return
	}

	tx.Result = "HTTP 2xx"
	tx.Context.SetLabel("region", "us-east-1")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return
	}

	fmt.Printf("Received %d bytes\n", len(body))
	tx.Context.SetLabel("response_size", len(body))
}

func printAPMEnv() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "ELASTIC_APM_") {
			fmt.Println(env)
		}
	}
}

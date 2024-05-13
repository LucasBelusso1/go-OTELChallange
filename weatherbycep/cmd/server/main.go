package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	config "github.com/LucasBelusso1/go-OTELChallange/weatherbycep/configs"
	"github.com/LucasBelusso1/go-OTELChallange/weatherbycep/internal/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const collectorUrl = "otel-collector:4317"

func init() {
	config.LoadConfig()
}

func main() {
	r := chi.NewRouter()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := shutdown(ctx)
		if err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get(`/{cep}`, handlers.GetTemperatureByZipCode)

	http.ListenAndServe(":8081", r)
	log.Printf("Listening on port %s", ":8081")
}

func initTracer() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("weatherbycep"),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}

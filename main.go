package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/exp/slices"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/shirou/gopsutil/cpu"

	apiv1 "connect-go-embed-vue/gen/api/v1"
	"connect-go-embed-vue/gen/api/v1/apiv1connect"
)

type SayHelloServer struct{}

func (s *SayHelloServer) SayHello(
	ctx context.Context,
	req *connect.Request[apiv1.SayHelloRequest],
) (*connect.Response[apiv1.SayHelloResponse], error) {
	res := connect.NewResponse(&apiv1.SayHelloResponse{
		Reply: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	return res, nil
}

type MetricsServer struct{}

func (s *MetricsServer) CpuUsageStream(
	ctx context.Context,
	req *connect.Request[apiv1.CpuUsageStreamRequest],
	stream *connect.ServerStream[apiv1.CpuUsageStreamResponse],
) error {
	for {
		percent, errPercent := cpu.Percent(1*time.Second, false)
		if errPercent != nil {
			log.Println(errPercent)
			return connect.NewError(connect.CodeUnknown, errPercent)
		}

		err := stream.Send(&apiv1.CpuUsageStreamResponse{
			Percent: float32(percent[0]),
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}
}

func getEnv(key string, fallback string, acceptable []string) string {
	if value, ok := os.LookupEnv(key); ok {
		if len(acceptable) == 0 || slices.Contains(acceptable, value) {
			return value
		}
	}
	return fallback
}

//go:embed frontend/dist/*
var static embed.FS

func main() {
	isDev := getEnv("APP_ENV", "production", []string{"development", "production"}) == "development"

	host := getEnv("APP_HOST", "0.0.0.0", nil)
	port := getEnv("APP_PORT", "8080", nil)

	mux := http.NewServeMux()

	if isDev {
		log.Println("[Dev] proxy vite dev page.")
		frontProxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = ":5173"
		}}
		mux.Handle("/", frontProxy)
	} else {
		log.Println("[Prod] serve embed source.")
		documentRoot, err := fs.Sub(static, "frontend/dist")
		if err != nil {
			log.Fatal(err)
		}
		mux.Handle("/", http.FileServer(http.FS(documentRoot)))
	}

	api := http.NewServeMux()
	api.Handle(apiv1connect.NewSayHelloServiceHandler(&SayHelloServer{}))
	api.Handle(apiv1connect.NewMetricsServiceHandler(&MetricsServer{}))
	mux.Handle("/api/", http.StripPrefix("/api", api))

	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Bind address %s", addr)

	http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c so we can serve HTTP/2 without TLS.
	)
}

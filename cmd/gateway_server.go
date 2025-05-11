package cmd

import (
	"context"
	suggest "darius/pkg/proto/suggest"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func startGateway() {
	grpcPort := viper.GetString("grpc.port")

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := suggest.RegisterSuggestServiceHandlerFromEndpoint(context.Background(), grpcMux, "localhost:"+grpcPort, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	mainMux := http.NewServeMux()

	mainMux.Handle("/metrics", promhttp.Handler())

	mainMux.Handle("/", grpcMux)

	gatewayPort := viper.GetString("gateway.port")
	log.Println("HTTP Gateway running on port " + gatewayPort)
	http.ListenAndServe(":"+gatewayPort, corsMiddleware(mainMux))
}

func test() {

}

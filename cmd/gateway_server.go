package cmd

import (
	"context"
	suggest "darius/pkg/proto/suggest"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	ctxdata "darius/ctx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	proto "google.golang.org/protobuf/proto"
)

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userIds, ok := r.Header["X-User-Id"]; ok {
			log.Printf("[Gateway] x-user-id: %v", userIds)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-user-id")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func customHeaderMatcher(key string) (string, bool) {
	if strings.ToLower(key) == "x-user-id" {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}

func forwardResponseFunc(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	smd, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}
	if vals := smd.HeaderMD.Get(ctxdata.HttpCodeHeader); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
	}
	return nil
}

// customErrorHandler handles custom error responses with proper HTTP status codes
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// Check if we have a custom HTTP status code in the context
	smd, ok := runtime.ServerMetadataFromContext(ctx)
	if ok {
		if vals := smd.HeaderMD.Get(ctxdata.HttpCodeHeader); len(vals) > 0 {
			code, parseErr := strconv.Atoi(vals[0])
			if parseErr == nil {
				w.WriteHeader(code)
				// Write error response
				errorResponse := map[string]interface{}{
					"error": err.Error(),
				}
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	// Default error handling
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func startGateway() {
	grpcPort := viper.GetString("grpc.port")

	grpcMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithForwardResponseOption(forwardResponseFunc),
		runtime.WithErrorHandler(customErrorHandler),
	)
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

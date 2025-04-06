package cmd

import (
	"darius/internal/handler"
	llm_grpc "darius/internal/llm-grpc"

	hello "darius/pkg/proto/hello"
	suggest "darius/pkg/proto/suggest"
	"fmt"
	"log"
	"net"
	"strings"

	"flag"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGRPC() {
	//server gateway
	port := viper.GetString("grpc.port")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//get llm host and model
	// llmHost := viper.GetString("llm.host")
	// if llmHost == "" || !strings.HasPrefix(llmHost, "http") {
	// 	llmHost = "http://104.199.250.71:2525/api/chat/completions"
	// }
	// llmModel := viper.GetString("llm.model")
	// if llmModel == "" || strings.HasPrefix(llmModel, "$") {
	// 	llmModel = "gpt-4o-mini"
	// }

	llmGRPCAddress := viper.GetString("LLM_GRPC_HOST")
	fmt.Print("llmGRPCAddress before hardcode: ", llmGRPCAddress)
	if llmGRPCAddress == "" || strings.HasPrefix(llmGRPCAddress, "$") {
		llmGRPCAddress = "arceus"
	}
	llmGRPCPort := viper.GetString("LLM_GRPC_PORT")
	fmt.Print("llmGRPCPort before hardcode: ", llmGRPCPort)

	if llmGRPCPort == "" || strings.HasPrefix(llmGRPCPort, "$") {
		llmGRPCPort = "8080"
	}
	llmGRPCModel := viper.GetString("llm.model")
	if llmGRPCModel == "" || strings.HasPrefix(llmGRPCModel, "$") {
		llmGRPCModel = "gpt-4o-mini"
	}
	addr := flag.String("addr", llmGRPCAddress+":"+llmGRPCPort, "the address to connect to")
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	arceusClient := suggest.NewArceusClient(conn)
	llmGRPCService := llm_grpc.NewService(arceusClient, llmGRPCModel)

	// r, err := c.GenerateText(context.Background(),
	// 	&suggest.GenerateTextRequest{
	// 		Model:   llmModel,
	// 		Content: "Hello, how are you?"})

	// LlmService := llm.NewLLM(&llm.Config{
	// 	Host:  llmHost,
	// 	Model: llmModel,
	// })
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.Content)

	handler := handler.NewHandlerWithDeps(handler.Dependency{
		// LlmService: LlmService,
		LLMGRPC: llmGRPCService,
	})

	grpcServer := grpc.NewServer()
	hello.RegisterHelloServiceServer(grpcServer, handler)
	suggest.RegisterSuggestServiceServer(grpcServer, handler)

	fmt.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

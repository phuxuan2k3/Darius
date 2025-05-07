package cmd

import (
	"context"
	"darius/cmd/db"
	"darius/internal/handler"
	f2_score "darius/internal/handler/f2-score"
	databaseService "darius/internal/services/database"
	llm_grpc "darius/internal/services/llm-grpc"
	llmManager "darius/managers/llm"
	hello "darius/pkg/proto/hello"
	suggest "darius/pkg/proto/suggest"
	"fmt"
	"log"
	"net"
	"strings"

	arceus "darius/pkg/proto/deps/arceus"
	"flag"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGRPC() {
	//server gateway
	port := viper.GetString("grpc.port")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}

	llmGRPCAddress := viper.GetString("LLM_GRPC_HOST")
	log.Print("llmGRPCAddress before hardcode: ", llmGRPCAddress)
	if llmGRPCAddress == "" || strings.HasPrefix(llmGRPCAddress, "$") {
		llmGRPCAddress = "arceus"
	}
	llmGRPCPort := viper.GetString("LLM_GRPC_PORT")
	log.Print("llmGRPCPort before hardcode: ", llmGRPCPort)

	if llmGRPCPort == "" || strings.HasPrefix(llmGRPCPort, "$") {
		llmGRPCPort = "8080"
	}
	llmGRPCModel := viper.GetString("LLM_GRPC_MODEL")
	log.Print("llmGRPCModel before hardcode: ", llmGRPCModel)

	if llmGRPCModel == "" || strings.HasPrefix(llmGRPCModel, "$") {
		llmGRPCModel = "gpt-4o-mini"
	}
	addr := flag.String("addr", llmGRPCAddress+":"+llmGRPCPort, "the address to connect to")
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	arceusClient := arceus.NewArceusClient(conn)
	log.Printf("Connected to LLM gRPC server at %s", *addr)
	llmGRPCService := llm_grpc.NewService(arceusClient, llmGRPCModel)

	// message queue
	f2scoreReqQueueAddr := viper.GetString("F2_SCORE_REQ_QUEUE_ADDRESS")
	f2scoreReqQueueName := viper.GetString("F2_SCORE_REQ_QUEUE_NAME")

	f2reqConn, f2reqCh, f2reqQ := conectQueue(f2scoreReqQueueAddr, f2scoreReqQueueName)
	if f2reqCh == nil || f2reqQ == nil {
		log.Printf("Failed to connect to RabbitMQ")
	}
	defer f2reqConn.Close()

	f2scoreRespQueueAddr := viper.GetString("F2_SCORE_RESP_QUEUE_ADDRESS")
	f2scoreRespQueueName := viper.GetString("F2_SCORE_RESP_QUEUE_NAME")
	f2respConn, f2respCh, f2respQ := conectQueue(f2scoreRespQueueAddr, f2scoreRespQueueName)
	if f2respCh == nil || f2respQ == nil {
		log.Printf("Failed to connect to RabbitMQ")
	}
	defer f2respConn.Close()

	db, err := db.NewDatabase()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	dbService := databaseService.NewService(db)
	llmManager := llmManager.NewManager(llmGRPCService, dbService)

	if f2reqCh != nil && f2respCh != nil {
		f2scoringHandler := f2_score.NewScoringHandler(llmManager, f2respCh, f2respQ)

		msgs, err := f2reqCh.Consume(f2reqQ.Name, "", false, false, false, false, nil)
		if err != nil {
			log.Print(err)
		}

		const maxWorker = 2
		for i := 0; i < maxWorker; i++ {
			go func() {
				for msg := range msgs {
					f2scoringHandler.Score(context.Background(), &f2_score.ScoreRequest{
						Msg: msg,
					})
					msg.Ack(false)
				}
			}()
		}
	}

	// r, err := c.GenerateText(context.Background(),
	// 	&suggest.GenerateTextRequest{
	// 		Model:   llmModel,
	// 		Content: "Hello, how are you?"})

	// LlmService := llm.NewLLM(&llm.Config{
	// 	Host:  llmHost,
	// 	Model: llmModel,
	// })
	// if err != nil {
	// 	log.Printff("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.Content)

	handler := handler.NewHandlerWithDeps(handler.Dependency{
		// LlmService: LlmService,
		LLMManager: llmManager,
	})

	grpcServer := grpc.NewServer()
	hello.RegisterHelloServiceServer(grpcServer, handler)
	suggest.RegisterSuggestServiceServer(grpcServer, handler)

	fmt.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func conectQueue(addr, queueName string) (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	if addr == "" || queueName == "" {
		log.Printf("RabbitMQ address or queue is not set")
		return nil, nil, nil
	}

	log.Printf("Connecting to RabbitMQ at %v, queue %v", addr, queueName)

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Print(err)
		return nil, nil, nil
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Print(err)
		return nil, nil, nil
	}

	return conn, ch, &q
}

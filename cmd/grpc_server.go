package cmd

import (
	"context"
	"darius/cmd/db"
	"darius/internal/handler"
	f2_score "darius/internal/handler/f2-score"
	bulbasaurService "darius/internal/services/bulbasaur"
	llm_grpc "darius/internal/services/llm-grpc"
	missfortune "darius/internal/services/missfortune"
	databaseService "darius/internal/services/repo"
	llmManager "darius/managers/llm"
	arceus "darius/pkg/proto/deps/arceus"
	"darius/pkg/proto/deps/bulbasaur"
	suggest "darius/pkg/proto/suggest"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	ctxdata "darius/ctx"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	userId, err := ctxdata.GetUserIdFromContext(ctx)
	if err != nil {
		log.Printf("Error getting user ID from context: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user ID from context")
	}

	if userId == "" {
		return nil, status.Error(codes.Unauthenticated, "user ID is required")
	}

	return handler(ctx, req)
}

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

	bulbasaurHost := viper.GetString("BULBASAUR_HOST")
	log.Print("bulbasaurHost before hardcode: ", bulbasaurHost)
	if bulbasaurHost == "" || strings.HasPrefix(bulbasaurHost, "$") {
		bulbasaurHost = "bulbasaur"
	}
	bulbasaurPort := viper.GetString("BULBASAUR_PORT")
	log.Print("bulbasaurPort before hardcode: ", bulbasaurPort)

	if bulbasaurPort == "" || strings.HasPrefix(bulbasaurPort, "$") {
		bulbasaurPort = "8080"
	}

	bulbasaurAddr := bulbasaurHost + ":" + bulbasaurPort
	bulbasaurConn, err := grpc.NewClient(bulbasaurAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer bulbasaurConn.Close()

	arceusClient := arceus.NewArceusClient(conn)
	log.Printf("Connected to LLM gRPC server at %s", *addr)
	llmGRPCService := llm_grpc.NewService(arceusClient, llmGRPCModel)

	bulbasaurClient := bulbasaur.NewVenusaurClient(bulbasaurConn)
	bulbasaurService := bulbasaurService.NewService(bulbasaurClient)

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

	missfortuneAddr := viper.GetString("MISSFORTUNE_ADDRESS")
	if missfortuneAddr == "" || strings.HasPrefix(missfortuneAddr, "$") {
		missfortuneAddr = "http://missfortune:8080"
	}

	missfortuneService := missfortune.NewService(missfortuneAddr, initMissfortuneClient())

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
					f2scoringHandler.ScoreV2(context.Background(), &f2_score.ScoreRequest{
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
		LLMManager:  llmManager,
		Missfortune: missfortuneService,
		Bulbasaur:   bulbasaurService,
	})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)
	// hello.RegisterHelloServiceServer(grpcServer, handler)
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

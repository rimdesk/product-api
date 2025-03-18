package main

import (
	"fmt"
	"net/http"

	"log"
	"os"

	"connectrpc.com/grpcreflect"
	"github.com/rimdesk/product-api/gen/protos/rimdesk/product/v1/productv1connect"

	"github.com/rimdesk/product-api/pkg/config"

	"github.com/rimdesk/product-api/pkg/data/repository"
	"github.com/rimdesk/product-api/pkg/database"
	
	"github.com/rimdesk/product-api/pkg/server"
	"github.com/rimdesk/product-api/pkg/service"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/gorm"
)

var (
	cfg = config.New()
	db  = database.NewGormDatabase()
)

func init() {
	cfg.LoadEnv()
	db.SetConfig(cfg.DatabaseConfig())
	db.ConnectDB()
}

func main() {
	serverAddr := fmt.Sprintf(":%s", os.Getenv("APP.PORT"))
	dbEngine := db.GetEngine().(*gorm.DB)

	inventoryRepository := repository.NewProductRepository(dbEngine)
	inventoryService := service.NewProductService(inventoryRepository)
	walletServer := server.NewProductServer(inventoryService)

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(productv1connect.ProductServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(productv1connect.NewProductServiceHandler(walletServer))

	// Start the server
	log.Printf("Starting server on %s...", serverAddr)
	err := http.ListenAndServe(
		serverAddr,
		h2c.NewHandler(mux, &http2.Server{}), // Use h2c for HTTP/2 without TLS
	)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/rimdesk/product-api/gen/rimdesk/product/v1/productv1connect"
	"github.com/rimdesk/product-api/pkg/auth"
	"github.com/rimdesk/product-api/pkg/config"
	"github.com/rimdesk/product-api/pkg/data/repository"
	"github.com/rimdesk/product-api/pkg/database"
	"github.com/rimdesk/product-api/pkg/helpers"
	middlewares "github.com/rimdesk/product-api/pkg/middlewares"
	"github.com/rimdesk/product-api/pkg/server"
	"github.com/rimdesk/product-api/pkg/service"
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
	serverAddr := cfg.GetServerAddress()
	dbEngine := db.GetEngine().(*gorm.DB)

	authenticator, err := auth.New()
	if err != nil {
		log.Println("Failed to create authenticator:", err)
		os.Exit(1)
	}

	contextHelper := helpers.NewContextHelper(authenticator)

	middleware := middlewares.New()
	interceptors := connect.WithInterceptors(
		middleware.CheckTenantIdPresenceInHeader(),
		middleware.UnaryTokenInterceptor(authenticator),
		middleware.LoggingUnaryInterceptor(),
	)

	productRepository := repository.NewProductRepository(dbEngine)
	productService := service.NewProductService(productRepository, contextHelper)
	productServer := server.NewProductServer(productService)

	mux := http.NewServeMux()
	path, handler := productv1connect.NewProductServiceHandler(productServer, interceptors)
	corsMiddleware := middleware.CorsMiddleware(handler)
	mux.Handle(path, corsMiddleware)

	reflector := grpcreflect.NewStaticReflector(productv1connect.ProductServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Start the server
	log.Printf("gRPC server started on port: %s...", serverAddr)
	if err := http.ListenAndServe(serverAddr, h2c.NewHandler(mux, cfg.Http2())); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

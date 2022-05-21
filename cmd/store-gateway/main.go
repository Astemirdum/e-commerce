package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Astemirdum/e-commerce/cmd/store-gateway/config"
	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/gen/openapiv2"
	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const textArt = "E-commerce GATEWAY"

func main() {
	fmt.Println(textArt)

	cfg := config.InitConfigs()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	s := newServer(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	log.Printf("server has started listen on %s", addr)
	go func() {
		if err := s.Run(addr); err != nil {
			log.Printf("server stop %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("graceful shutdown")
}

func newServer(cfg *config.Config) *gin.Engine {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			token := request.Header.Get("Authorization")
			log.Println("token", token)
			return metadata.Pairs("auth", token)
		}),
	)

	{
		if err := authv1.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux,
			cfg.Auth, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
			log.Fatal(err)
		}
		if err := orderv1.RegisterOrderServiceHandlerFromEndpoint(context.Background(), mux,
			cfg.Order, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
			log.Fatal(err)
		}
		if err := productv1.RegisterProductServiceHandlerFromEndpoint(context.Background(), mux,
			cfg.Product, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
			log.Fatal(err)
		}
	}

	return initRoutes(mux)
}

func initRoutes(mux *runtime.ServeMux) *gin.Engine {
	r := gin.Default()

	{
		r.GET("/swagger/:name", openapiv2.SwagHandler()) // auth | product | order
		r.GET("/swagger-ui/*any", gin.WrapH(openapiv2.SwagUIHandler("/swagger-ui")))
	}

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", gin.WrapH(mux))
		auth.POST("/login", gin.WrapH(mux))
	}

	api.POST("/create-order", gin.WrapH(mux))

	api.POST("/create-product", gin.WrapH(mux))
	api.GET("/get-product/:id", gin.WrapH(mux))

	return r
}

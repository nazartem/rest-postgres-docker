package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"restapi-lesson/internal/buyer"
	buyerDB "restapi-lesson/internal/buyer/db"
	"restapi-lesson/internal/config"
	"restapi-lesson/internal/logging"
	"restapi-lesson/internal/note"
	noteDB "restapi-lesson/internal/note/db"
	"restapi-lesson/internal/prdlist"
	productListDB "restapi-lesson/internal/prdlist/db"
	"restapi-lesson/internal/product"
	productDB "restapi-lesson/internal/product/db"
	"restapi-lesson/pkg/client/postgresql"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	logger := &logging.Logger{
		Info: infoLog,
		Err:  errorLog,
	}

	infoLog.Println("create router")
	router := httprouter.New()

	cfg := config.GetConfig(logger)

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		errorLog.Fatalf("%v", err)
	}

	productRepository := productDB.NewRepository(postgreSQLClient, logger)
	logger.Info.Println("register product handler")
	productHandler := product.NewHandler(productRepository, logger)
	productHandler.Register(router)

	buyerRepository := buyerDB.NewRepository(postgreSQLClient, logger)
	logger.Info.Println("register buyer handler")
	buyerHandler := buyer.NewHandler(buyerRepository, logger)
	buyerHandler.Register(router)

	noteRepository := noteDB.NewRepository(postgreSQLClient, logger)
	logger.Info.Println("register note handler")
	noteHandler := note.NewHandler(noteRepository, logger)
	noteHandler.Register(router)

	productListRepository := productListDB.NewRepository(postgreSQLClient, logger)
	logger.Info.Println("register productList handler")
	productListHandler := prdlist.NewHandler(productListRepository, logger)
	productListHandler.Register(router)

	start(router, cfg, logger)
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {
	logger.Info.Println("start application")

	var listener net.Listener
	var listenErr error

	logger.Info.Println("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Info.Printf("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	if listenErr != nil {
		logger.Err.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Err.Fatal(server.Serve(listener))
}

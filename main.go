package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fakhripraya/emailing-service/data"
	"github.com/fakhripraya/emailing-service/entities"
	protos "github.com/fakhripraya/emailing-service/protos/email"
	"github.com/fakhripraya/emailing-service/server"

	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var err error

func main() {
	logger := hclog.Default()

	// load configuration from env file
	logger.Info("Loading env")
	err = godotenv.Load(".env")

	if err != nil {
		// log the fatal error if load env failed
		log.Fatal(err)
	}

	// Initialize app configuration
	logger.Info("Initialize application configuration")
	var appConfig entities.Configuration
	err = data.ConfigInit(&appConfig)

	if err != nil {
		// log the fatal error if initializing app configuration failed
		log.Fatal(err)
	}

	logger.Info("Creating a new email flow interface")
	email := data.NewEmail(logger)

	// create a new gRPC server, use WithInsecure to allow http connections
	logger.Info("Creating a new gRPC server")
	gs := grpc.NewServer()

	// create an instance of the mailer server
	logger.Info("Creating a new mailer instance")
	mailer := server.NewMailer(logger, email, &appConfig.EmailCredential)

	// register the mailer server
	logger.Info("Registering mailer into the gRPC server")
	protos.RegisterEmailServer(gs, mailer)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	logger.Info("Registering reflection service")
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	logger.Info("Creating TCP socket on " + appConfig.EmailConfig.Host + ":" + appConfig.EmailConfig.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":"+appConfig.EmailConfig.Port))
	if err != nil {
		logger.Error("Unable to create listener", "error", err.Error())
		os.Exit(1)
	}

	// listen for requests
	logger.Info("Successfully creating TCP socket")
	gs.Serve(listener)
}

package main

import (
	"log"
	"net/http"
	"os"

	"app/internal/endpoint"
	"app/internal/entity"
	"app/internal/service"
	"app/internal/transport"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env loaded")
	}

	infServ := service.InfoServices{
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		TokenHost: os.Getenv("TOKEN_HOST"),
		TokenPort: os.Getenv("TOKEN_PORT"),
		Secret:    os.Getenv("SECRET"),
	}

	runServer(
		os.Getenv("PORT"),
		&infServ,
	)
}

func runServer(port string, infServ *service.InfoServices) {
	svc := service.NewService(
		&http.Client{},
		infServ,
	)

	getSignUpHandler := httptransport.NewServer(
		endpoint.MakeSignUpEndpoint(svc),
		transport.DecodeRequestWithBody(entity.UsernamePasswordEmailRequest{}),
		transport.EncodeResponse,
	)

	getSignInHandler := httptransport.NewServer(
		endpoint.MakeSignInEndpoint(svc),
		transport.DecodeRequestWithBody(entity.UsernamePasswordRequest{}),
		transport.EncodeResponse,
	)

	getLogOutHandler := httptransport.NewServer(
		endpoint.MakeLogOutEndpoint(svc),
		transport.DecodeRequestWithHeader(entity.Token{}),
		transport.EncodeResponse,
	)

	getAllUsersHandler := httptransport.NewServer(
		endpoint.MakeGetAllUsersEndpoint(svc),
		transport.DecodeRequestWithoutBody(),
		transport.EncodeResponse,
	)

	getProfileHandler := httptransport.NewServer(
		endpoint.MakeProfileEndpoint(svc),
		transport.DecodeRequestWithHeader(entity.Token{}),
		transport.EncodeResponse,
	)

	getDeleteAccountHandler := httptransport.NewServer(
		endpoint.MakeDeleteAccountEndpoint(svc),
		transport.DecodeRequestWithHeader(entity.Token{}),
		transport.EncodeResponse,
	)

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/signup").Handler(getSignUpHandler)
	router.Methods(http.MethodPost).Path("/signin").Handler(getSignInHandler)
	router.Methods(http.MethodPost).Path("/logout").Handler(getLogOutHandler)
	router.Methods(http.MethodGet).Path("/users").Handler(getAllUsersHandler)
	router.Methods(http.MethodPost).Path("/profile").Handler(getProfileHandler)
	router.Methods(http.MethodDelete).Path("/profile").Handler(getDeleteAccountHandler)

	log.Println("ListenAndServe on localhost:" + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+port, router))
}

package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mehmetalisavas/message-sender/internal/api"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Routers initializes the routes and Swagger documentation
// @title Message Sender API
// @version 1.0
// @description This is the API documentation for the message sender service.
// @host localhost:8080
// @BasePath /

// Routers function sets up the routes and handlers
// @Summary Set up routes for message processing
// @Description Define all the routes for message processing and listing
// @Tags routes
// @Accept json
// @Produce json
func Routers(api *api.Api) http.Handler {
	r := mux.NewRouter()

	// Process message command (start/stop)
	// @Summary Update message processing
	// @Description Start or stop the message processing based on the command
	// @Accept json
	// @Produce json
	// @Param command query string true "Command: start or stop"
	// @Success 200 {string} string "Message processing started/stopped"
	// @Failure 400 {string} string "Command is required or invalid command"
	// @Router /process_message [get]
	r.HandleFunc("/process_message", api.UpdateMessageProcessing).Methods("GET")

	// List sent messages with pagination
	// @Summary List sent messages
	// @Description Get a list of sent messages with optional pagination parameters
	// @Accept json
	// @Produce json
	// @Param limit query int false "Limit of messages to return"
	// @Param offset query int false "Offset for pagination"
	// @Param page query int false "Page number"
	// @Success 200 {array} models.Message "List of sent messages"
	// @Failure 500 {string} string "Internal server error"
	// @Router /messages [get]
	r.HandleFunc("/messages", api.ListSentMessages).Methods("GET")

	// Serve the Swagger UI at /swagger route
	// r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
	// 	httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	// 	httpSwagger.DeepLinking(true),
	// 	httpSwagger.DocExpansion("none"),
	// 	httpSwagger.DomID("swagger-ui"),
	// )).Methods(http.MethodGet)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	return r
}

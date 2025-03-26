package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mehmetalisavas/message-sender/internal/models"
)

// ListSentMessages handles listing the sent messages with optional pagination
// @Summary List sent messages
// @Description Get a list of sent messages with optional pagination parameters (limit, offset, page)
// @Param limit query int false "Limit of messages to return"
// @Param offset query int false "Offset for pagination"
// @Param page query int false "Page number"
// @Success 200 {array} models.Message "List of sent messages"
// @Failure 500 {string} string "Internal server error"
// @Router /messages [get]
func (a *Api) ListSentMessages(w http.ResponseWriter, r *http.Request) {
	opts := models.ListOptions{}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		opts.Limit, _ = strconv.Atoi(limit)
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		opts.Offset, _ = strconv.Atoi(offset)
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		opts.Page, _ = strconv.Atoi(page)
	}

	messages, err := a.storageService.ListSentMessages(r.Context(), opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// UpdateMessageProcessing handles the command to start or stop message processing
// @Summary Update message processing state
// @Description Start or stop the message processing based on the command (start/stop)
// @Param command query string true "Command: start or stop"
// @Success 200 {string} string "Message processing started or stopped"
// @Failure 400 {string} string "Command is required or invalid command"
// @Router /process_message [get]
func (a *Api) UpdateMessageProcessing(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("command")

	if command == "" {
		http.Error(w, "command is required", http.StatusBadRequest)
		return
	}

	switch command {
	case "start":
		a.config.SetMessageProcessing(true)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message processing started"))
	case "stop":
		a.config.SetMessageProcessing(false)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message processing stopped"))
	default:
		http.Error(w, "invalid command", http.StatusBadRequest)
	}
}

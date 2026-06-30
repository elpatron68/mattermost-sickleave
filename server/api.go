package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/medisoftware/mattermost-sickleave/server/command"
)

type sickLeaveContextResponse struct {
	Active          any    `json:"active,omitempty"`
	MaxBackdateDays int    `json:"max_backdate_days"`
	CommandTrigger  string `json:"command_trigger"`
}

func (p *Plugin) initRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(p.MattermostAuthorizationRequired)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/context", p.handleSickLeaveContext).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dialog/submit", p.handleDialogSubmit).Methods(http.MethodPost)
	apiRouter.HandleFunc("/end", p.handleEnd).Methods(http.MethodPost)

	return router
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) handleSickLeaveContext(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	settings := p.settingsFromConfig()

	active, err := p.kvstore.GetActive(userID)
	if err != nil {
		p.API.LogError("Failed to load active sick leave", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	maxBackdate := settings.MaxBackdateDays
	if maxBackdate <= 0 {
		maxBackdate = 3
	}

	response := sickLeaveContextResponse{
		Active:          active,
		MaxBackdateDays: maxBackdate,
		CommandTrigger:  command.NormalizeCommandTrigger(settings.CommandTrigger),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.API.LogError("Failed to encode sick leave context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Plugin) handleDialogSubmit(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.API.LogError("Failed to read dialog submission", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request model.SubmitDialogRequest
	if err = json.Unmarshal(body, &request); err != nil {
		p.API.LogError("Failed to decode dialog submission", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.UserId == "" {
		request.UserId = r.Header.Get("Mattermost-User-ID")
	}

	var response *model.SubmitDialogResponse
	response, err = p.command.SubmitDialog(&request)
	if err != nil {
		p.API.LogError("Failed to process dialog submission", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.API.LogError("Failed to encode dialog response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type endRequest struct {
	ChannelID string `json:"channel_id"`
}

type endResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (p *Plugin) handleEnd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")

	var request endRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil && err != io.EOF {
			p.API.LogError("Failed to decode end request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	response, err := p.command.End(userID, request.ChannelID)
	if err != nil {
		p.API.LogError("Failed to end sick leave", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	payload := endResponse{Message: response.Text}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		p.API.LogError("Failed to encode end response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

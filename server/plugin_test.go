package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elpatron68/mattermost-sickleave/server/dialog"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCommand struct{}

func (stubCommand) Handle(_ *model.CommandArgs) (*model.CommandResponse, error) {
	return &model.CommandResponse{}, nil
}

func (stubCommand) SubmitDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error) {
	return &model.SubmitDialogResponse{Errors: map[string]string{}}, nil
}

func (stubCommand) End(_ string, _ string) (*model.CommandResponse, error) {
	return &model.CommandResponse{Text: "closed"}, nil
}

func (stubCommand) EnsureSlashCommandRegistered() error {
	return nil
}

func TestDialogSubmitRoute(t *testing.T) {
	plugin := Plugin{
		command: stubCommand{},
	}
	plugin.router = plugin.initRouter()

	payload, err := json.Marshal(model.SubmitDialogRequest{
		CallbackId: dialog.CallbackStart,
		UserId:     "user-1",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/dialog/submit", bytes.NewReader(payload))
	r.Header.Set("Mattermost-User-ID", "user-1")

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	_ = result.Body.Close()
}

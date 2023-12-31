package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vynnydev/chatgpt-microservices-k8s/packages/microservices/chatservice/usecases/chatcompletion"
)

type WebChatGPTHandler struct {
	Completionusecases chatcompletion.ChatCompletionusecases
	Config             chatcompletion.ChatCompletionConfigInputDTO
	AuthToken          string
}

func NewWebChatGPTHandler(usecases chatcompletion.ChatCompletionusecases, config chatcompletion.ChatCompletionConfigInputDTO, authToken string) *WebChatGPTHandler {
	return &WebChatGPTHandler{
		Completionusecases: usecases,
		Config:             config,
		AuthToken:          authToken,
	}
}

func (h *WebChatGPTHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Authorization") != h.AuthToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !json.Valid(body) {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var dto chatcompletion.ChatCompletionInputDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dto.Config = h.Config

	result, err := h.Completionusecases.Execute(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

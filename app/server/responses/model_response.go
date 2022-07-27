package responses

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	UnMarshalRequestError     = "E001"
	InvalidBodyError          = "E002"
	DataBaseQueryFailureError = "E003"
	MarshalError              = "E004"
	ResourceNotFound          = "E005"
	ResourceFinished          = "E006"
)

type ErrorResponse struct {
	Code string `json:"errorCode,omitempty"`
	// Unique object id
	Message string `json:"message,omitempty"`
}

func GenerateErrorResponseBody(ctx context.Context, errorCode string, message string) ErrorResponse {
	return ErrorResponse{Code: errorCode, Message: message}
}

func WriteError(ctx context.Context, w http.ResponseWriter, statusCode int, errorBody ErrorResponse) {
	response, err := json.Marshal(errorBody)
	if err != nil {
		http.Error(w, "couldn't marshal error body", http.StatusInternalServerError)
		return
	}

	log.Ctx(ctx).Warn().Int("statusCode", statusCode).Msg(errorBody.Message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}

func WriteOkResponse(ctx context.Context, w http.ResponseWriter, responseBody interface{}) {
	writeResponse(ctx, w, responseBody, http.StatusOK)
}

func WriteCreatedResponse(ctx context.Context, w http.ResponseWriter, responseBody interface{}) {
	writeResponse(ctx, w, responseBody, http.StatusCreated)
}

func WriteNoContentResponse(ctx context.Context, w http.ResponseWriter) {
	writeResponse(ctx, w, nil, http.StatusNoContent)
}

func writeResponse(ctx context.Context, w http.ResponseWriter, responseBody interface{}, statusCode int) {
	response, err := json.Marshal(responseBody)
	if err != nil {
		WriteError(
			ctx,
			w,
			http.StatusInternalServerError,
			GenerateErrorResponseBody(
				ctx,
				MarshalError,
				"Encoding error",
			),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		log.Ctx(ctx).Warn().AnErr("error", err).Msgf("error writing the response body")
	}
}

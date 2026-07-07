package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
)

func TestRespondServiceErrorIncludesValidationDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	respondServiceError(ctx, &services.ValidationError{Field: "email", Reason: "already_registered", Message: "Email is already registered"})

	var body struct {
		Error   string `json:"error"`
		Details struct {
			Field  string `json:"field"`
			Reason string `json:"reason"`
		} `json:"details"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error != "validation_failed" || body.Details.Field != "email" || body.Details.Reason != "already_registered" {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}

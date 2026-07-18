package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testJSONRequest struct {
	Email string `json:"email" validate:"required,email"`
}

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

func TestBindAndValidateJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns a typed request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"email":"alice@example.com"}`))
		ctx.Request.Header.Set("Content-Type", "application/json")

		req, ok := bindAndValidateJSON[testJSONRequest](ctx, validator.New())
		if !ok || req.Email != "alice@example.com" {
			t.Fatalf("unexpected bind result: %#v, ok=%v", req, ok)
		}
	})

	t.Run("keeps validation error details", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"email":"invalid"}`))
		ctx.Request.Header.Set("Content-Type", "application/json")

		_, ok := bindAndValidateJSON[testJSONRequest](ctx, validator.New())
		if ok || recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected validation failure, code=%d, ok=%v", recorder.Code, ok)
		}

		var body struct {
			Details struct {
				Field  string `json:"field"`
				Reason string `json:"reason"`
			} `json:"details"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Details.Field != "email" || body.Details.Reason != "invalid" {
			t.Fatalf("unexpected response: %s", recorder.Body.String())
		}
	})
}

func TestRequireUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns authenticated id", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		expected := primitive.NewObjectID()
		ctx.Set("user_id", expected.Hex())

		actual, ok := requireUserID(ctx)
		if !ok || actual != expected {
			t.Fatalf("unexpected id: %s, ok=%v", actual.Hex(), ok)
		}
	})

	t.Run("writes the shared unauthorized response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Set("user_id", "invalid")

		_, ok := requireUserID(ctx)
		if ok || recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected unauthorized response, code=%d, ok=%v", recorder.Code, ok)
		}
	})
}

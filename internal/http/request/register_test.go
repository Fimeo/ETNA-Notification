package request

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/maxatome/go-testdeep/td"
)

func TestRegisterRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	var err error
	ctx.Request, err = http.NewRequest(
		http.MethodPost,
		"/register",
		io.NopCloser(strings.NewReader(`{"discordAccountName": "Test#1234", "password": "password", "login": "login"}`)))
	td.CmpNoError(t, err)

	inputBody, err := GetValidatedRegisterInputPayload(ctx)
	td.CmpNoError(t, err)
	td.Cmp(t, inputBody, &RegisterInput{
		DiscordAccountName: "Test#1234",
		Password:           "password",
		Login:              "login",
	})
}

package httpcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"etna-notification/internal/controller"
	"etna-notification/internal/http/request"
)

type HTTPRegisterController struct {
	register controller.IRegisterController
}

// NewRegisterController boostrap quiz endpoints with their controllers.
// These routes are located under /quiz prefix.
func NewRegisterController(e *gin.Engine, controller controller.IRegisterController) *HTTPRegisterController {
	httpController := &HTTPRegisterController{
		register: controller,
	}

	group := e.Group("/register")
	group.POST("", httpController.Register)

	return httpController
}

func (q *HTTPRegisterController) Register(c *gin.Context) {
	req, err := request.GetValidatedRegisterInputPayload(c)
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	register, err := q.register.Register(req.Login, req.Password, req.DiscordAccountName)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invitation": register})
}

func ValidationError(err error) (int, any) {
	return http.StatusUnprocessableEntity, gin.H{"message": err.Error()}
}

func ProcessError(err error) (int, any) {
	return http.StatusBadRequest, gin.H{"message": err.Error()}
}

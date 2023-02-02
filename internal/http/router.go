package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"go.uber.org/fx"

	"etna-notification/internal/controller"
	"etna-notification/internal/http/httpcontroller"
	"etna-notification/internal/http/middleware"
)

// SetupRouter boostrap the application by loading and init dependencies.
func SetupRouter(lc fx.Lifecycle, controllers controller.Controllers) *gin.Engine {
	e := gin.Default()

	// Enable CORS allow origin
	e.Use(middleware.CORSMiddleware())

	// Active kin validation based on swagger interface
	e.Use(middleware.Kin(middleware.NewKinValidator()))

	httpcontroller.NewRegisterController(e, controllers.IRegisterController)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go ListenAndServe(e)
			return nil
		},
	})

	return e
}

func ListenAndServe(e *gin.Engine) {
	if err := e.Run(); err != nil {
		panic(fmt.Errorf("fatal error Gin server: %w", err))
	}
}

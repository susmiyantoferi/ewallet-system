package routes

import (
	"ewallet/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(users controller.UserController) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	v1 := r.Group("/api/v1")
	{
		usr := v1.Group("/users")
		{
			usr.POST("", users.Create)
			usr.PATCH("/:id", users.Update)
			usr.GET("/:id", users.GetByID)
			usr.DELETE("/:id", users.Delete)

		}
	}

	return r
}

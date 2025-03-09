// @Title       	Task API
// @Version     	1.0
// @Description 	Esta é uma API de gerenciamento de tarefas.
// @Host        	localhost:8080
// @BasePath    	/
package delivery

import (
	"go-clean-arch/internal/delivery/dependencies"
	"go-clean-arch/internal/interfaces/handlers"
	"log"
	"time"

	_ "go-clean-arch/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	container := dependencies.Setup()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := container.Invoke(func(taskHandler *handlers.TaskHandler) {
		router.POST("/tasks", taskHandler.CreateTask)
		router.GET("/tasks", taskHandler.GetTasks)
		router.PUT("/tasks/:id", taskHandler.UpdateTask)
		router.DELETE("/tasks/:id", taskHandler.DeleteTask)
	})
	if err != nil {
		log.Fatalf("Erro ao resolver dependências: %v", err)
	}

	log.Println("Server started at: 8080")
	router.Run(":8080")
}

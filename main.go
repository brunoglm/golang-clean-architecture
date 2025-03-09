// @Title       	Task API
// @Version     	1.0
// @Description 	Esta é uma API de gerenciamento de tarefas.
// @Host        	localhost:8080
// @BasePath    	/
package main

import (
	"context"
	"go-clean-arch/models"
	"log"
	"net/http"
	"time"

	_ "go-clean-arch/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func connectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("not connected")
		log.Fatal(err)
	}

	collection = client.Database("taskdb").Collection("tasks")
	log.Println("Connected to MongoDB")
}

// @Summary 		Cria uma nova tarefa
// @Description Cria uma nova tarefa
// @Tags 				tasks
// @Accept 			json
// @Produce 		json
// @Param 			task	body 			models.Task true "Dados da Tarefa"
// @Success 		201 	{object}	map[string]interface{}
// @Failure 		400 	{object} 	map[string]string
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks [post]
func createTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// @Summary 		Lista todas as tarefas
// @Description Retorna todas as tarefas cadastradas
// @Tags 				tasks
// @Produce 		json
// @Success 		200 	{array}		models.Task
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks [get]
func listTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary 		Atualiza uma tarefa
// @Description Atualiza os dados de uma tarefa pelo ID
// @Tags 				tasks
// @Accept 			json
// @Produce 		json
// @Param 			id		path 			string 			true "ID da Tarefa"
// @Param 			task	body 			models.Task true "Dados atualizados da Tarefa"
// @Success 		200 	{object}	map[string]string
// @Failure 		400 	{object} 	map[string]string
// @Failure 		404 	{object} 	map[string]string
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks/{id} [put]
func updateTask(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": task})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// @Summary 		Deleta uma tarefa
// @Description Remove uma tarefa pelo ID
// @Tags 				tasks
// @Produce 		json
// @Param 			id		path 			string 			true "ID da Tarefa"
// @Success 		200 	{object}	map[string]string
// @Failure 		400 	{object} 	map[string]string
// @Failure 		404 	{object} 	map[string]string
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks/{id} [delete]
func deleteTask(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func main() {
	connectDB()

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

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.GET("/tasks", listTasks)

	router.POST("/tasks", createTask)

	router.PUT("/tasks/:id", updateTask)

	router.DELETE("/tasks/:id", deleteTask)

	router.Run(":8080")
}

package handlers

import (
	"go-clean-arch/internal/entities"
	"go-clean-arch/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	usecase usecases.TaskUseCase
}

func NewTaskHandler(usecase usecases.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		usecase: usecase,
	}
}

// @Summary 		Cria uma nova tarefa
// @Description Cria uma nova tarefa
// @Tags 				tasks
// @Accept 			json
// @Produce 		json
// @Param 			task	body 			entities.Task true "Dados da Tarefa"
// @Success 		201 	{object}	map[string]interface{}
// @Failure 		400 	{object} 	map[string]string
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task entities.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.usecase.CreateTask(c.Request.Context(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary 		Lista todas as tarefas
// @Description Retorna todas as tarefas cadastradas
// @Tags 				tasks
// @Produce 		json
// @Success 		200 	{array}		entities.Task
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.usecase.GetTasks(c.Request.Context())
	if err != nil {
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
// @Param 			task	body 			entities.Task true "Dados atualizados da Tarefa"
// @Success 		200 	{object}	map[string]string
// @Failure 		400 	{object} 	map[string]string
// @Failure 		404 	{object} 	map[string]string
// @Failure 		500 	{object} 	map[string]string
// @Router 			/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not sented"})
		return
	}

	var task entities.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.usecase.UpdateTask(c.Request.Context(), id, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not sented"})
		return
	}

	err := h.usecase.DeleteTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

package rest

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"planner/entity"
	"strconv"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}

}
func (s *Service) CreateTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var task entity.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var ID int
	err = s.db.QueryRowContext(ctx, "INSERT INTO planner (title, description, status, priority) VALUES (?, ?, ?, ?) RETURNING ID",
		task.Title, task.Description, task.Status, task.Priority).Scan(&ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ID": ID})
}

func (s *Service) GetTasks(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctx.Header("Cache-Control", "public, max-age=3600")
	getStatus := ctx.Query("status")
	status, err := strconv.ParseBool(getStatus)
	if err != nil {
		rows, err := s.db.QueryContext(ctxt, "SELECT * FROM planner")
		if err != nil {
			return
		}
		var sl []entity.Task
		for rows.Next() {
			var t entity.Task
			rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority)
			sl = append(sl, t)
		}
		ctx.JSONP(http.StatusOK, sl)
		return
	}
	rows, err := s.db.Query("SELECT * FROM planner WHERE status = ?", status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var sl []entity.Task
	var flag bool

	for rows.Next() {
		flag = true
		var t entity.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority)
		if err != nil {
			return
		}
		sl = append(sl, t)
	}

	if flag == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Нет задач с таким статусом"})
		return
	}
	ctx.JSONP(http.StatusOK, sl)
}

func (s *Service) GetTaskByID(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	taskID := s.getTaskID(ctx)
	if taskID == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}
	rows, err := s.db.QueryContext(ctxt, "SELECT * FROM planner WHERE ID = ?", taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	var t entity.Task
	err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, t)
}

func (s *Service) getTaskID(ctx *gin.Context) int {
	taskIDStr := ctx.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return -1
	}

	return taskID
}

func (s *Service) UpdateTasks(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	key := ctx.Param("id")
	ID, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный индекс"})
		return
	}
	var task entity.Task
	err = ctx.BindJSON(&task)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, err)
		return
	}
	_, err = s.db.ExecContext(ctxt, "UPDATE planner SET title = ?, description = ?, status = ?, priority = ? WHERE ID = ?", task.Title, task.Description, task.Status, task.Priority, ID)
	ctx.JSON(http.StatusOK, gin.H{"Статус": "изменения сохранены"})
}

func (s *Service) DelTask(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	key := ctx.Param("id")
	num, err := strconv.Atoi(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный индекс"})
		return
	}
	result, err := s.db.ExecContext(ctxt, "DELETE FROM planner WHERE ID = ?", num)
	i, err := result.RowsAffected()
	if i == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Задача не найдена"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Статус": "Задача удалена"})
}

package controllers

import (
	"course-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Tasks struct {
	DB *gorm.DB
}

type taskResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	CourseID uint   `json:"courseId"`
	Course   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"course"`
	Question []struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Answer string `json:"answer"`
		Hint   string `json:"hint"`
	} `json:"question"`
}

type taskCreateResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type allTaskResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type createTaskForm struct {
	Name     string `form:"name" binding:"required"`
	Desc     string `form:"desc" binding:"required"`
	CourseID uint   `form:"courseId" binding:"required"`
}

type updateTaskForm struct {
	Name string `form:"name"`
	Desc string `form:"desc"`
}

func (c *Tasks) FindAll(ctx *gin.Context) {
	var tasks []models.Task
	c.DB.Order("id desc").Find(&tasks)

	var serializedTask []allTaskResponse
	copier.Copy(&serializedTask, &tasks)
	ctx.JSON(http.StatusOK, gin.H{"tasks": serializedTask})
}

func (c *Tasks) FindOne(ctx *gin.Context) {
	task, err := c.findTaskByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedTask taskResponse
	copier.Copy(&serializedTask, &task)
	ctx.JSON(http.StatusOK, gin.H{"task": serializedTask})
}

func (c *Tasks) Create(ctx *gin.Context) {
	var form createTaskForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	copier.Copy(&task, &form)
	if err := c.DB.Create(&task).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedTask taskCreateResponse
	copier.Copy(&serializedTask, &task)
	ctx.JSON(http.StatusCreated, gin.H{"task": serializedTask})
}

func (c *Tasks) Update(ctx *gin.Context) {
	var form updateTaskForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	task, err := c.findTaskByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&task).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedTask taskResponse
	copier.Copy(&serializedTask, &task)
	ctx.JSON(http.StatusOK, gin.H{"task": serializedTask})
}

func (c *Tasks) Delete(ctx *gin.Context) {
	task, err := c.findTaskByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Unscoped().Delete(&task)
	ctx.Status(http.StatusNoContent)
}

func (c *Tasks) findTaskByID(ctx *gin.Context) (*models.Task, error) {
	var task models.Task
	id := ctx.Param("id")

	if err := c.DB.Preload("Course").Preload("Question").First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

package http

import (
	"net/http"
	"strconv"
	"to-do-list/internal/model"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	taskUsecase model.ITaskUsecase
}

func NewTaskHandler(e *echo.Echo, taskUsecase model.ITaskUsecase) {
	handler := &TaskHandler{
		taskUsecase: taskUsecase,
	}

	routeTask := e.Group("v1/tasks")
	routeTask.GET("", handler.FindAll, AuthMiddleware)
	routeTask.GET("/:id", handler.FindById, AuthMiddleware)
	routeTask.POST("/create", handler.Create, AuthMiddleware)
	routeTask.PUT("/update/:id", handler.Update, AuthMiddleware)
	routeTask.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (handler *TaskHandler) FindAll(c echo.Context) error {
	var filter model.FindAllParam

	// Bind query params for pagination
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters")
	}

	tasks, err := handler.taskUsecase.FindAll(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   tasks,
	})
}

func (handler *TaskHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	task, err := handler.taskUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Task Not Found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   task,
	})
}

func (handler *TaskHandler) Create(c echo.Context) error {
	var body model.CreateTaskInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := handler.taskUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Task created successfully",
	})
}

func (handler *TaskHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var body model.UpdateTaskInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = handler.taskUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		if err.Error() == "task not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Task updated successfully",
		Data:    body,
	})
}

func (handler *TaskHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = handler.taskUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Task deleted successfully",
	})
}

package controller

import (
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
}

type userControllerImpl struct {
	UserService service.UserService
}

func NewUserControllerImpl(userService service.UserService) UserController {
	return &userControllerImpl{
		UserService: userService,
	}
}

func (u *userControllerImpl) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	resp, err := u.UserService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("failed create user"))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(resp))

}

func (u *userControllerImpl) Update(c *gin.Context) {
	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	userId := c.Param("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	req.ID = id

	resp, err := u.UserService.Update(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resp))

}

func (u *userControllerImpl) Delete(c *gin.Context) {
	userId := c.Param("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	if err := u.UserService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil))

}

func (u *userControllerImpl) GetByID(c *gin.Context) {
	userId := c.Param("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	resp, err := u.UserService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resp))
}

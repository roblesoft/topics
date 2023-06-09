package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	entity "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/pkg/token"
)

type UserParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) Register(ctx *gin.Context) {
	var (
		user   entity.User
		params UserParams
	)

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Username = params.Username
	user.Password = params.Password

	if err := server.Service.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func (server *Server) Login(ctx *gin.Context) {
	var params UserParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := entity.User{}
	u.Username = params.Username
	u.Password = params.Password

	token, err := server.Service.LoginCheck(u.Username, u.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (server *Server) CurrentUser(ctx *gin.Context) *entity.User {
	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	user, err := server.Service.UserRepo.GetUserById(user_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	return user
}

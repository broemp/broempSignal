package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/broemp/broempSignal/db/sqlc"
)

type createUserRequest struct {
	Username  string `json:"username" binding:"required"`
	DiscordId int64  `json:"discordid" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:  req.Username,
		Discordid: req.DiscordId,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) Create_user_if_not_exists(ctx *gin.Context, userid int64, username string) (db.User, error) {
	arg := db.CreateUserParams{
		Username:  username,
		Discordid: userid,
	}

	user, err := server.store.GetUser(ctx, userid)
	if err == sql.ErrNoRows {
		user, err = server.store.CreateUser(ctx, arg)
		if err != nil {
			log.Println("user creation error")
			return db.User{}, err
		}

	} else if err != nil {
		log.Println("other error")
		return db.User{}, err
	}

	log.Println("user created")
	return user, nil
}

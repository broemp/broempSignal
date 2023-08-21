package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request_afk_create struct {
	Username string `json:"username" binding:"required"`
	Userid   int64  `json:"userid" binding:"required"`
}

func (server *Server) createAFK(ctx *gin.Context) {
	var req request_afk_create
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Create_user_if_not_exists(ctx, req.Userid, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	afk, err := server.store.CreateAFK(ctx, sql.NullInt64{Int64: user.Discordid, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, afk)
}

type listAFKRequest struct {
	Userid int64 `uri:"id" binding:"required"`
}

func (server *Server) listAFK(ctx *gin.Context) {
	var req listAFKRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	afk, err := server.store.ListAFK(ctx, sql.NullInt64{Int64: req.Userid, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, afk)
}

func (server *Server) countAFK(ctx *gin.Context) {
	var req listAFKRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	afk, err := server.store.GetAFKCount(ctx, sql.NullInt64{Int64: req.Userid, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, afk)
}

func (server *Server) listAFKByCount(ctx *gin.Context) {
	afk, err := server.store.ListUsersByAFKCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, afk)
}

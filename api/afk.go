package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type afkRequest struct {
	User int64 `json:"discordid" binding:"required"`
}

func (server *Server) createAFK(ctx *gin.Context) {
	var req afkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	afk, err := server.store.CreateAFK(ctx, sql.NullInt64{Int64: req.User, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, afk)
}

type listAFKRequest struct {
	User int64 `uri:"id" binding:"required"`
}

func (server *Server) listAFK(ctx *gin.Context) {
	var req listAFKRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	afk, err := server.store.ListAFK(ctx, sql.NullInt64{Int64: req.User, Valid: true})
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
	afk, err := server.store.GetAFKCount(ctx, sql.NullInt64{Int64: req.User, Valid: true})
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

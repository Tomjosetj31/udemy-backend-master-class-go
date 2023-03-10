package api

import (
	"database/sql"
	"net/http"

	simplebank "github.com/Tomjosetj31/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err :=ctx.ShouldBindJSON(&req); err !=nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := simplebank.CreateAccountParams{
		Owner: req.Owner,
		Balance: 0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err :=ctx.ShouldBindUri(&req); err !=nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`

}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err :=ctx.ShouldBindQuery(&req); err !=nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg :=simplebank.ListAccountParams{
		Limit: req.PageSize,
		Offset: (req.PageID-1)*req.PageSize,
	}
	accounts, err := server.store.ListAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err :=ctx.ShouldBindJSON(&req); err !=nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := simplebank.UpdateAccountParams{
		ID: req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err :=ctx.ShouldBindUri(&req); err !=nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusNoContent, "")
}


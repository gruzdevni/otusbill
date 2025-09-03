package restapi

import (
	"errors"

	"otusbill/apperr"
	"otusbill/internal/models"
	"otusbill/internal/restapi/operations/balance"
	"otusbill/internal/restapi/operations/other"
	"otusbill/internal/restapi/operations/user_c_r_u_d"
	"otusbill/internal/service/api/bill"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	billSrv bill.Service
}

func NewHandler(billSrv bill.Service) *Handler {
	return &Handler{
		billSrv: billSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) ReduceBalance(params balance.PostUserBalanceReduceParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	reqParams := params.Request

	err := h.billSrv.ReduceBalance(ctx, *reqParams)
	if err != nil {
		if errors.Is(err, apperr.NotEnoughMoney) {
			return balance.NewPostUserBalanceReduceForbidden().WithPayload(&models.DefaultStatusResponse{
				Message: apperr.NotEnoughMoney.Error(), Code: "03",
			})
		}

		return balance.NewPostUserBalanceReduceInternalServerError()
	}

	return balance.NewPostUserBalanceReduceOK()
}

func (h *Handler) IncreaseBalance(params balance.PostUserBalanceIncreaseParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	reqParams := params.Request

	err := h.billSrv.IncreaseBalance(ctx, *reqParams)
	if err != nil {
		return balance.NewPostUserBalanceIncreaseInternalServerError()
	}

	return balance.NewPostUserBalanceIncreaseOK()
}

func (h *Handler) GetUserBalance(params balance.GetUserBalanceGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		errText = err.Error()
		return balance.NewGetUserBalanceGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.billSrv.GetUserBalance(ctx, userGUID)
	if err != nil {
		errText = err.Error()
		return balance.NewGetUserBalanceGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return balance.NewGetUserBalanceGUIDOK().WithPayload(&res)
}

func (h *Handler) CreateUser(params user_c_r_u_d.PostUserParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.Request.GUID.String())
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPostUserBadRequest().WithPayload(&models.DefaultStatusResponse{Code: "3", Message: errText})
	}

	err = h.billSrv.InsertUser(ctx, userGUID)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPostUserInternalServerError().WithPayload(&models.DefaultStatusResponse{Code: "3", Message: errText})
	}

	return user_c_r_u_d.NewPostUserOK()
}

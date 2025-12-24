package httpserver

import (
	"gameapp/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userRegisterHandler(ctx echo.Context) error {
	var uReq userservice.RegisterRequest
	err := ctx.Bind(&uReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, err := s.userSvc.Register(uReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, resp)
}

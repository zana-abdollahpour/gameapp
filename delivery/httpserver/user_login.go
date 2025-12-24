package httpserver

import (
	"gameapp/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userLoginHandler(ctx echo.Context) error {
	var loginRequest userservice.LoginRequest
	if err := ctx.Bind(&loginRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, loginErr := s.userSvc.Login(loginRequest)
	if loginErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, loginErr.Error())
	}

	return ctx.JSON(http.StatusOK, resp)
}

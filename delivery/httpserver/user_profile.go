package httpserver

import (
	"gameapp/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userProfileHandler(ctx echo.Context) error {
	authToken := ctx.Request().Header.Get("Authorization")
	claims, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, profileErr := s.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if profileErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, profileErr.Error())
	}

	return ctx.JSON(http.StatusOK, resp)

}

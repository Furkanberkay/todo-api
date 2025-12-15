package httpx

import (
	"strconv"
	"todoApp3/internal/domain"

	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) *domain.MyCustomClaim {
	return c.Get("user").(*domain.MyCustomClaim)
}

func GetClaimsUserID(c echo.Context) uint {
	return GetClaims(c).UserID
}

func GetIDParam(c echo.Context, key string) (uint, error) {
	idStr := c.Param(key)
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		return 0, domain.ErrInvalidID
	}

	return uint(id), nil

}

package httpx

import (
	"todoApp3/internal/domain"

	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) *domain.MyCustomClaim {
	claims := c.Get("user").(domain.MyCustomClaim)

	return &claims
}

func GetClaimsUserID(c echo.Context, myClaim *domain.MyCustomClaim) uint {
	return GetClaims(c).UserID
}

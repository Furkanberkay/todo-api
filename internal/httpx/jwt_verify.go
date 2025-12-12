package httpx

import (
	"context"
	"log/slog"
	"strings"
	"todoApp3/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JwtVerify struct {
	secretKey string
	logger    *slog.Logger
}

func NewJwtVerify(secretKey string, logger *slog.Logger) *JwtVerify {
	return &JwtVerify{
		secretKey: secretKey,
		logger:    logger,
	}
}

func (j *JwtVerify) Verify(ctx context.Context, tokenString string) (*domain.MyCustomClaim, error) {
	if tokenString == "" {
		j.logger.Warn("authorization header missing",
			"component", "jwt_verify",
		)
		return nil, domain.ErrUnAuthorized
	}

	if len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer ") {
		j.logger.Warn("authorization header malformed",
			"component", "jwt_verify",
		)
		return nil, domain.ErrUnAuthorized
	}

	tokenstr := tokenString[7:]
	myClaim := domain.MyCustomClaim{}

	token, err := jwt.ParseWithClaims(tokenstr, &myClaim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.logger.Error("unexpected signing method",
				"method", token.Header["alg"],
				"component", "jwt_verify",
			)
			return nil, domain.ErrUnAuthorized
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		j.logger.Warn("jwt parse/validate failed",
			"error", err.Error(),
			"component", "jwt_verify",
		)
		return nil, domain.ErrUnAuthorized
	}

	if !token.Valid {
		j.logger.Warn("jwt token invalid",
			"component", "jwt_verify",
		)
		return nil, domain.ErrUnAuthorized
	}

	j.logger.Debug("jwt token verified successfully",
		"user_id", myClaim.UserID,
		"component", "jwt_verify",
	)

	return &myClaim, nil
}

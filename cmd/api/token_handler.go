package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		fmt.Println("error", err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	newAccessToken, _, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newRefreshToken, _, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// vazi za https i localhost

	// ctx.SetCookie("access_token", newAccessToken, int(server.config.AccessTokenDuration.Seconds()), "/", "localhost", true, true)
	// ctx.SetCookie("refresh_token", newRefreshToken, int(server.config.RefreshTokenDuration.Seconds()), "/", "localhost", true, true)

	ctx.SetCookie(
		"access_token", // 1) Name
		newAccessToken, // 2) Value
		int(server.config.AccessTokenDuration.Seconds()), // 3) MaxAge (u sekundama)
		"/",                          // 4) Path
		server.config.FrontendDomain, // 5) Domain
		server.config.SecureCookie,   // 6) Secure
		true,                         // 7) HttpOnly
	)

	ctx.SetCookie(
		"refresh_token", // 1) Name
		newRefreshToken, // 2) Value
		int(server.config.RefreshTokenDuration.Seconds()), // 3) MaxAge (u sekundama)
		"/",                          // 4) Path
		server.config.FrontendDomain, // 5) Domain
		server.config.SecureCookie,   // 6) Secure
		true,                         // 7) HttpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tijanadmi/ddn_rdc/models"
	"github.com/tijanadmi/ddn_rdc/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username string   `json:"username"`
	FullName string   `json:"full_name"`
	Role     []string `json:"user_role"`
}

func newUserResponse(user *models.User) userResponse {
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {

		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	/*session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}*/

	// **Set HttpOnly cookie za refresh token**
	/*ctx.SetCookie(
		"refresh_token",
		refreshToken,
		int(server.config.RefreshTokenDuration.Seconds()),
		"/auth/refresh", // putanja za refresh
		"",              // domen (prazno = trenutni)
		true,            // Secure (HTTPS)
		true,            // HttpOnly
	)*/

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		Path:  "/",
		// Domain:   "localhost",
		MaxAge:   int(server.config.AccessTokenDuration.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/",
		// Domain:   "localhost",
		MaxAge:   int(server.config.RefreshTokenDuration.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	rsp := loginUserResponse{
		//SessionID:             session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		//RefreshToken:          refreshToken,
		//RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

// // Custom JWT regex: 3 dela razdvojena tačkama
// var jwtRegex = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`)

// // Custom validator funkcija
// func validateJWT(fl validator.FieldLevel) bool {
// 	token := fl.Field().String()
// 	return jwtRegex.MatchString(token)
// }

// type GetUserByTokenRequest struct {
// 	AccessToken string `json:"access_token" binding:"required,jwt"`
// }

type GetUserByTokenRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (server *Server) GetUserByToken(ctx *gin.Context) {
	// _, err := server.authorizeUser(ctx)
	// if err != nil {
	// 	return nil, unauthenticatedError(err)
	// }

	// validate := validator.New()

	// // Registruj custom validator
	// validate.RegisterValidation("jwt", validateJWT)

	var req GetUserByTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(req)
	accessToken := req.AccessToken
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(payload)

	user, err := server.store.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return
	}

	fmt.Println(user)

	rsp := userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}
	ctx.JSON(http.StatusOK, rsp)
}

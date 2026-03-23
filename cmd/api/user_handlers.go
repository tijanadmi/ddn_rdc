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
	Username string `json:"username"`
	FullName string `json:"full_name"`
	// Role     []string `json:"user_role"`
	Role []string `json:"user_role"`

	DDN *models.DDNData `json:"ddn,omitempty"`
	TDN *models.TDNData `json:"tdn,omitempty"`
	PGI *models.PGIData `json:"pgi,omitempty"`
}

// func newUserResponse(user *models.User) userResponse {
// 	return userResponse{
// 		Username: user.Username,
// 		FullName: user.FullName,
// 		Role:     user.Roles,
// 	}
// }

// func newUserResponse(user *models.User) userResponse {
// 	// Napravi slice stringova iz Role.Code
// 	var roleCodes []string
// 	for _, r := range user.Roles {
// 		roleCodes = append(roleCodes, r.Code)
// 	}

// 	return userResponse{
// 		Username: user.Username,
// 		FullName: user.FullName,
// 		Role:     roleCodes, // samo Code svake role
// 	}
// }

func newUserResponse(user *models.User) userResponse {
	var roleCodes []string
	resp := userResponse{
		Username: user.Username,
		FullName: user.FullName,
	}

	// mapa: code -> funkcija koja popunjava response
	roleHandlers := map[string]func(r models.Role){
		"DDN": func(r models.Role) {
			if r.DDN.TipPrivPrip != "" {
				resp.DDN = &r.DDN
			}
		},
		"TDN": func(r models.Role) {
			if r.TDN.TipPrivPrip != "" {
				resp.TDN = &r.TDN
			}
		},
		"PGI": func(r models.Role) {
			if r.PGI.TipPrivPrip != "" {
				resp.PGI = &r.PGI
			}
		},
	}

	for _, r := range user.Roles {
		roleCodes = append(roleCodes, r.Code)

		if handler, ok := roleHandlers[r.Code]; ok {
			handler(r)
		}
	}

	resp.Role = roleCodes
	return resp
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

func getRoleCodes(roles []models.Role) []string {
	var codes []string
	for _, r := range roles {
		codes = append(codes, r.Code)
	}
	return codes
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	// fmt.Println(user)
	if err != nil {

		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	roleCodes := getRoleCodes(user.Roles)

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		roleCodes,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		roleCodes,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// kod koji radi za produkciju

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   server.config.FrontendDomain,
		MaxAge:   int(server.config.AccessTokenDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,       // ⬅ radi bez HTTPS inace ako je HTTPS http.SameSiteNoneMode
		Secure:   server.config.SecureCookie, // ⬅ važno za HTTP inace je true ako je HTTPS
	})
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   server.config.FrontendDomain,
		MaxAge:   int(server.config.RefreshTokenDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,       // ⬅ radi bez HTTPS inace ako je HTTPS http.SameSiteNoneMode
		Secure:   server.config.SecureCookie, // ⬅ važno za HTTP inace je true ako je HTTPS
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

func (server *Server) logoutUser(ctx *gin.Context) {
	// Clear the access_token cookie

	/*var UserLogout struct {
		UserId string `json:"user_id"`
	}*/

	/*err := ctx.ShouldBindJSON(&UserLogout)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}*/

	// key := ctx.Get("authorizationPayloadKey")

	// fmt.Println("User ID from Logout request:", key.ID)

	// kod za localhost

	/*http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  "access_token",
		Value: "",
		Path:  "/",
		// Domain:   "localhost",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})*/

	// kod za produkciju

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   server.config.FrontendDomain,
		MaxAge:   -1,
		Secure:   server.config.SecureCookie, // ⬅ važno za HTTP inace je true ako je HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // ⬅ radi bez HTTPS inace ako je HTTPS http.SameSiteNoneMode
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   server.config.FrontendDomain,
		MaxAge:   -1,
		Secure:   server.config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

type GetUserByTokenRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (server *Server) GetUserByToken(ctx *gin.Context) {

	token, err := util.GetAccessToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("nije pronađen access token u cookie-ju: %w", err)))
		return
	}

	// fmt.Println("access token je ", token)
	payload, err := server.tokenMaker.VerifyToken(token)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// fmt.Println(payload)

	user, err := server.store.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return
	}

	// fmt.Println(user)

	// Mapiramo Role -> []string (samo Code)
	var roleCodes []string
	for _, r := range user.Roles {
		roleCodes = append(roleCodes, r.Code)
	}

	// rsp := userResponse{
	// 	Username: user.Username,
	// 	FullName: user.FullName,
	// 	Role:     roleCodes,
	// }

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

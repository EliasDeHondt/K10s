/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var envUsername, usernameOk = os.LookupEnv("USERNAME")
var envPassword, passwordOk = os.LookupEnv("PASSWORD")

func Init() {
	if !usernameOk {
		envUsername = "admin"
	}
	if !passwordOk {
		envPassword = "password"
	}
}

func HandleLogin(ctx *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if credentials.Username != envUsername || credentials.Password != envPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateToken(credentials.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.SetCookie("jwt", token, 3600*24, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func HandleLogout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

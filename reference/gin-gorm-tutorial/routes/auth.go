package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"../config"
	"../models"
	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)



func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "success"})
}

// Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
	}

	providerScopes := map[string][]string{
		"github":   []string{"public_repo"},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, user)
	
	var jwtToken = createToken(&newUser)

	c.JSON(200, gin.H{
		"data": newUser,
		"token": jwtToken,
		"message": "Successfully login",
	});
	
}

func getOrRegisterUser(provider string, user *structs.User) models.User{
	var userData models.User

	config.DB.Where("provider = ? and social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User {
			FullName: user.FullName,
			Email: user.Email,
			SocialID: user.ID,
			Provider: provider,
			Avatar: user.Avatar,
		}
		config.DB.Create(&newUser)
		return newUser
	} else {
		return userData
	}
}


func createToken(user *models.User) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"user_role": user.Role,
		"exp": json.Number(strconv.FormatInt(time.Now().Add(time.Hour*time.Duration(1)).Unix(), 10)),
		"iat": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SUPER_SECRET")))

	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}


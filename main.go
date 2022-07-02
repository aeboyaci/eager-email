package main

import (
	"eager-email/api/db"
	"eager-email/api/routes"
	"eager-email/api/security"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	db.Connect()
	defer db.Disconnect()

	ConfigureOauth2()

	app := gin.New()

	app.POST("/api/account/sign-in", routes.SignIn)
	app.POST("/api/account/sign-up", routes.SignUp)
	app.GET("/api/account/oauth2/sign-in", routes.Oauth2SignIn)
	app.GET("/api/account/oauth2/google/callback/", routes.Oauth2Callback)

	app.GET("/api/account/test", security.Authorize, func(ctx *gin.Context) {
		email, _ := ctx.Get("email")

		ctx.JSON(200, gin.H{
			"email": email,
		})
	})

	app.Run("0.0.0.0:8080")
}

func ConfigureOauth2() {
	googleClientKey := os.Getenv("CLIENT_KEY")
	googleClientSecret := os.Getenv("CLIENT_SECRET")

	// "https://www.googleapis.com/auth/gmail.send"
	googleScopes := []string{"profile", "email"}

	// // cookieSecret := os.Getenv("COOKIE_SECRET")
	cookieSecret := "COOKIE_SECRET"

	cookieStore := sessions.NewCookieStore([]byte(cookieSecret))
	cookieStore.MaxAge(1 * 60 * 60) // 1 hour
	cookieStore.Options.Path = "/"
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.Secure = false

	gothic.Store = cookieStore

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleClientSecret,
			"http://localhost:8080/api/account/oauth2/google/callback",
			googleScopes...,
		),
	)
}

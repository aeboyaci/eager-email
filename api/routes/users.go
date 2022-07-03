package routes

import (
	"eager-email/api/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	FullName string `json:"fullName" bson:"fullName"`
}

func SignIn(ctx *gin.Context) {
	user := new(User)

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})

		return
	}

	dbUser := new(User)
	err := db.FindOne("users", bson.M{"email": user.Email}).Decode(dbUser)

	if err != nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Email or password is wrong",
		})

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Email or password is wrong",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})
	tokenString, _ := token.SignedString([]byte("JWT_SECRET"))

	ctx.SetCookie("token", tokenString, 1*60*60, "/", "eager-email.ahmeterenboyaci.com", false, true)

	ctx.JSON(200, gin.H{
		"success":  true,
		"email":    user.Email,
		"fullName": user.FullName,
	})
}

func SignUp(ctx *gin.Context) {
	user := new(User)

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})

		return
	}

	dbUser := new(User)
	err := db.FindOne("users", bson.M{"email": user.Email}).Decode(dbUser)

	if err == nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "User already exist",
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   "Hashing error occured",
		})

		return
	}

	user.Password = string(hashedPassword)

	_, err = db.InsertOne("users", user)

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   "Database error occured",
		})

		return
	}

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "User created successfully",
	})
}

func Oauth2SignIn(ctx *gin.Context) {
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func Oauth2Callback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	dbUser := new(User)
	err = db.FindOne("users", bson.M{"email": user.Email}).Decode(dbUser)

	if err != nil {
		_, err = db.InsertOne("users", user)

		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   "Database error occured",
			})

			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})
	tokenString, _ := token.SignedString([]byte("JWT_SECRET"))

	ctx.SetCookie("token", tokenString, 1*60*60, "/", "eager-email.ahmeterenboyaci.com", false, true)

	ctx.Redirect(http.StatusTemporaryRedirect, "/api/account/test")
}

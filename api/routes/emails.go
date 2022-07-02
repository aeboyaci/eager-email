package routes

import (
	"crypto/rand"
	"eager-email/api/db"
	"fmt"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Email struct {
	From      string    `bson:"from"`
	To        string    `json:"to" bson:"to"`
	Subject   string    `json:"subject" bson:"subject"`
	TraceCode string    `bson:"traceCode"`
	IsRead    bool      `bson:"isRead"`
	Count     int64     `bson:"count"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func GenerateRandomString() string {
	alphabet := "qwertyuopiasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	randomString := ""

	for i := 0; i < 128; i++ {
		bInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		index := bInt.Int64()

		randomString += string(alphabet[index])
	}

	return randomString
}

func GetEmails(ctx *gin.Context) {
	userEmailAddress, _ := ctx.Get("email")

	dbUser := new(User)
	err := db.FindOne("users", bson.M{"email": userEmailAddress}).Decode(dbUser)

	if err != nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "User not found",
		})

		return
	}

	result, err := db.FindMany("emails", bson.M{"from": userEmailAddress})

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   "Databse error occurred",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"emails":  result,
	})
}

func CreateEmailTracer(ctx *gin.Context) {
	userEmailAddress, _ := ctx.Get("email")

	dbUser := new(User)
	err := db.FindOne("users", bson.M{"email": userEmailAddress}).Decode(dbUser)

	if err != nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "User not found",
		})

		return
	}

	email := new(Email)

	if err := ctx.BindJSON(email); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})

		return
	}

	email.From = userEmailAddress.(string)
	email.IsRead = false

	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	email.TraceCode = GenerateRandomString()
	email.Count = 0

	_, err = db.InsertOne("emails", email)

	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   "Database error occured",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"success":        true,
		"tracerImageUrl": fmt.Sprintf("http://localhost:8080/cleardot.gif?code=%s", email.TraceCode),
	})
}

func TrackEmail(ctx *gin.Context) {
	traceCode := ctx.Query("code")

	dbEmail := new(Email)
	err := db.FindOne("emails", bson.M{"traceCode": traceCode}).Decode(dbEmail)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid trace code given",
		})

		return
	}

	db.UpdateOne("emails", bson.M{"traceCode": traceCode}, bson.M{"$set": bson.M{"isRead": true, "updatedAt": time.Now(), "count": dbEmail.Count + 1}})

	ctx.Next()
}

// http://localhost:8080/images?code=fLCgroApvDwMRBxjKwTsteJFXvYAgQSUHxWVehKlhTeAspMVkkLEozZADdhUgxTRMLMgPVYobWBrOGWKrDjuYpQgnwSdbZLgClAbTLkXhhVDFQdyLZpqvNsPvvRbcJGU

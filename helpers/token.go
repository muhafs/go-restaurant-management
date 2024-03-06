package helpers

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateTokens(user models.User) (accesstoken string, refreshToken string, err error) {
	accesstoken, err = GenerateAccessToken(user, SECRET_KEY, 24)
	refreshToken, err = GenerateRefreshToken(SECRET_KEY, 168)

	if err != nil {
		return
	}

	return
}

func GenerateAccessToken(user models.User, secret string, expiry int) (token string, err error) {
	claims := &SignedDetails{
		Email:      *user.Email,
		First_name: *user.FirstName,
		Last_name:  *user.LastName,
		Uid:        user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return
	}

	return
}

func GenerateRefreshToken(secret string, expiry int) (token string, err error) {
	claims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return
	}

	return
}

func UpdateTokens(c context.Context, accessToken string, refreshToken string, userID string) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var updateObj bson.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: accessToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)

	return
}

func ValidateToken(accessToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err.Error()
	}

	//the token is invalid
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, "the token is invalid"
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, "token is expired"
	}

	return
}

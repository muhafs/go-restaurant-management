package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/helpers"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func FindUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * limit
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "user_items", Value: bson.D{
			{Key: "$slice", Value: []interface{}{
				"$data", startIndex, limit,
			}},
		}},
	}}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []bson.M
	if err := result.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users[0])
}

func FindUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	userID := c.Param("user_id")

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	// extract user's data request into struct
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate user's data request
	if vErr := validate.Struct(user); vErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
		return
	}

	// check for the existence email in DB (if exists then return email already taken)
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// check for the existence phone in DB (if exists then return phone already taken)
	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already taken"})
		return
	}

	// hash the password
	pass, _ := helpers.HashPassword(*user.Password)
	user.Password = &pass

	// make current timestamps for request (created, update)
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// generate new objectID
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()

	// generate token and refresh token
	token, refreshToken, _ := helpers.GenerateTokens(user)
	user.Token = &token
	user.RefreshToken = &refreshToken

	// insert all data into DB
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// return response
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	// extract user's data request into struct
	var credentials models.User
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check for the existence email in DB
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this account is not registered"})
		return
	}

	// verify the password
	if isValid := helpers.VerifyPassword(*credentials.Password, *user.Password); isValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	// generate new tokens
	accessToken, refreshToken, _ := helpers.GenerateTokens(user)

	// re-assign tokens with the new generated tokens
	helpers.UpdateTokens(c, accessToken, refreshToken, user.UserID)

	// return response
	c.JSON(http.StatusOK, user)
}

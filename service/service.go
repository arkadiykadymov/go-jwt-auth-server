package service

import (
	"context"
	"dev/jwt-auth-server/auth"
	"dev/jwt-auth-server/config/db"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	var user auth.User

	Collection := db.GetRemoteDBConnection()
	error := c.ShouldBindJSON(&user)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Field Name(s)/ Value(s)"})
		return
	}

	error = user.Validate()

	if error != nil {
		message := "User " + error.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	// Inserts ID for the user object
	user.GUID = uuid.New().String()

	user.Password = auth.HasPassword(user.Password)

	user.AccessToken = "0"

	user.RefreshToken = "0"

	_, error = Collection.InsertOne(context.TODO(), user)

	if error != nil {
		message := "User " + error.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully!", "resourceId": user.GUID})

}

// LoginUser Method
func LoginUser(c *gin.Context) {
	var user auth.User
	userGUID := c.Param("uuid")

	Collection := db.GetRemoteDBConnection()

	err := Collection.FindOne(context.TODO(), bson.D{{"Guid", userGUID}}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid user ID"})
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(user.Username, user.GUID)
	if err != nil || accessToken == "" || refreshToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not generate token"})
		return
	}

	hashedRefreshToken := auth.HasPassword(refreshToken)

	if user.AccessToken == "0" && user.RefreshToken == "0" {

		user.AccessToken = accessToken
		user.RefreshToken = hashedRefreshToken

		update := bson.M{
			"$set": bson.M{
				"access_token":  accessToken,
				"refresh_token": hashedRefreshToken,
				"isused":        false,
			},
		}

		_, err = Collection.UpdateOne(context.TODO(), bson.D{{"Guid", user.GUID}}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not insert user"})
			return
		}

	} else {

		user.AccessToken = accessToken
		user.RefreshToken = hashedRefreshToken
		user.IsUsed = false
		_, error := Collection.InsertOne(context.TODO(), user)

		if error != nil {
			message := "User " + error.Error()
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "access_token": accessToken, "refresh_token": refreshToken})
}

// RefreshAccessToken method creates a new access_token, when the user provides an unexpired and validrefresh_token
func RefreshAccessToken(c *gin.Context) {

	var tokenRequest auth.RefreshTokenRequestBody

	Collection := db.GetRemoteDBConnection()
	err := c.ShouldBindJSON(&tokenRequest)
	if err != nil {
		message := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	valid, id, _, err := auth.ValidateToken(tokenRequest.RefreshToken)
	if valid == false || err != nil {
		message := err.Error()
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": message})
		c.Abort()
		return
	}

	if valid == true && id != "" {

		cursor, err := Collection.Find(context.TODO(), bson.D{{"Guid", id}})
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid id"})
			return
		}

		for cursor.Next(context.TODO()) {
			//Create a value into which the single document can be decoded
			var user auth.User
			err := cursor.Decode(&user)
			if err != nil {
				log.Fatal(err)
			}

			println(user.RefreshToken)

			if auth.ComparePasswords(user.RefreshToken, []byte(tokenRequest.RefreshToken)) && !user.IsUsed {
				user.IsUsed = true
				update := bson.M{
					"$set": bson.M{
						"isused": true,
					},
				}
				_, err = Collection.UpdateOne(context.TODO(), bson.D{{"Guid", id}}, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not insert user"})
					return
				}

				newToken, err := auth.GenerateAccessToken(user.Username, id)
				if err != nil {
					//logger.Logger.Errorf(err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Cannot Generate New Access Token"})
					c.Abort()
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "access_token": newToken, "refresh_token": tokenRequest.RefreshToken})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Token already used"})
				c.Abort()
				return
			}

		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error Found "})
}

func DeleteAllRefreshTokens(c *gin.Context) {

	userUUID := c.Param("uuid")

	Collection := db.GetRemoteDBConnection()

	result, err := Collection.DeleteMany(context.TODO(), bson.D{{"Guid", userUUID}})
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid id"})
		return
	}
	fmt.Print(result)

}

func DeleteRefreshToken(c *gin.Context) {
	var tokenRequest auth.RefreshTokenRequestBody

	Collection := db.GetRemoteDBConnection()

	err := c.ShouldBindJSON(&tokenRequest)
	if err != nil {
		message := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": message})
		return
	}

	valid, id, _, err := auth.ValidateToken(tokenRequest.RefreshToken)
	if valid == false || err != nil {
		message := err.Error()
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": message})
		c.Abort()
		return
	}

	if valid == true && id != "" {

		cursor, err := Collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid id"})
			return
		}

		for cursor.Next(context.TODO()) {
			var user auth.User
			err := cursor.Decode(&user)
			if err != nil {
				log.Fatal(err)
			}

			match := auth.ComparePasswords(user.RefreshToken, []byte(tokenRequest.RefreshToken))
			if match {
				_, err := Collection.DeleteOne(context.TODO(), bson.D{{"refresh_token", user.RefreshToken}})
				if err != nil {
					log.Fatal(err)
				}
				c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Token deleted"})
				return

			}

		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error Found "})

}

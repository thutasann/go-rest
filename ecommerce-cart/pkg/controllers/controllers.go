package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thutasann/ecommerce-cart/pkg/database"
	"github.com/thutasann/ecommerce-cart/pkg/models"
	tokengen "github.com/thutasann/ecommerce-cart/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

// Hash Password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// Verify Password
func verifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is Incorrect"
		valid = false
	}
	return valid, msg
}

// SignUp Controller
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, channel = context.WithTimeout(context.Background(), 100*time.Second)
		defer channel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// validate
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		// check existing email
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already existed"})
		}

		// check existing phone
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer channel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		// HashPassword
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := tokengen.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not"})
			return
		}
		defer channel()
		c.JSON(http.StatusCreated, "Successfully signed up!")
	}
}

// Login Controller
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, channel = context.WithTimeout(context.Background(), 100*time.Second)
		defer channel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer channel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found, Invalid credentials", "details": err})
			return
		}
		PasswordIsValid, msg := verifyPassword(*user.Password, *founduser.Password)
		defer channel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println("Invalid Password --> ", msg)
		}
		token, refreshToken, _ := tokengen.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer channel()
		tokengen.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		c.JSON(http.StatusFound, founduser)
	}
}

// Product Viewer Admin Controller
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Search Product Controller
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Search Product By Query Controller
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/helpers"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func HashPassword(password string) string{
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err !=nil{
		log.Panic(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string)(bool,string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword),[]byte(userPassword))
	check := true
	msg := ""

	if err !=nil{
		msg = fmt.Sprintf("Incorrect Password or Email")
		check = false
	}
	return check,msg 
}

func SignUp()gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user model.User

		if err := c.BindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)
		if validationErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			defer cancel()
			return
		}

		count,err := userCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking for email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count,err = userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err !=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking for phone number"})
		}

		if count>0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this phone number or email already exists"})
		}

		user.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken,_ :=helpers.GenerateAllTokens(*user.Email,*user.First_name,*user.Last_name,*&user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertionErr := userCollection.InsertOne(ctx,user)
		if insertionErr !=nil{
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error" : msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,resultInsertionNumber)
	}
}

func Login()gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		var user model.User
		var foundUser model.User
		
		if err := c.BindJSON(&user);err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			defer cancel()
			return
		}
		
		err := userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"email or password is incorrect"})
			return
		}

		passwordIsValid,msg := VerifyPassword(*user.Password,*foundUser.Password)
		defer cancel()
		if !passwordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not found"})
		}

		token,refreshToken,_ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name,*foundUser.Last_name,*&foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,foundUser)
	}
}

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		// if err := helpers.CheckUserType(c,"ADMIN");err!=nil{
		// 	c.JSON(http.StatusBadGateway,gin.H{"error":err.Error()})
		// 	return
		// }
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
			
		recordPerPage,err := strconv.Atoi(c.Query("recordPerPage"))	
		if err !=nil || recordPerPage <1{
			recordPerPage = 10
		}
		page ,err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}
		startIndex := (page-1)*recordPerPage
		startIndex,err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match",bson.D{{}}}}
		groupStage := bson.D{{"$group",bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count",bson.D{{"$sum",1}}},
			{"data",bson.D{{"$push","$$ROOT"}}},
		}}}
		projectStage := bson.D{
			{"$project",bson.D{
				{"_id",0},
				{"total_count",1},
				{"user_items",bson.D{{"$slice", []interface{}{"$data",startIndex,recordPerPage}}}},
			}},
		}
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}
		result,err := userCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error": "error occured while listing user items"})
		}
		var allUsers[]bson.M
		if err := result.All(ctx, &allUsers);err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allUsers[0])
	}
}

func GetUser()gin.HandlerFunc{
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		// if err := helpers.MatchUserTypeToUid(c,userId);err !=nil{
		// 	c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		// 	return
		// }

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user model.User
		err := userCollection.FindOne(ctx,bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,user)
	}
}
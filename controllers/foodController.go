package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client,"food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		recordPerPage,err :=strconv.Atoi(c.Query("recordPerPage"))
		if err !=nil || recordPerPage<1{
			recordPerPage = 10
		}
		page,err := strconv.Atoi(c.Query("page"))
		if err!=nil || page<1{
			page=1
		}
		startIndex := (page-1)*recordPerPage
		startIndex,err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match",bson.D{{}}}}
		groupStage	:=bson.D{{"$group",bson.D{{"_id",bson.D{{"_id","null"}}},{"total_count",bson.D{{"$sum",1}}},{"data",bson.D{{"$push","$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project",bson.D{
					{"_id",0},
					{"total_count",1},
					{"food_items",bson.D{{"$slice",[]interface{}{"$data",startIndex,recordPerPage}}}},
				}}}

		result,err :=foodCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing items"})
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods);err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allFoods[0])
  }
}

func GetFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		foodId :=c.Param("food_id")
		var food model.Food
		
		err:=foodCollection.FindOne(ctx,bson.M{"food_id":foodId}).Decode(&food)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var menu model.Menu
		var food model.Food

		if err:= c.BindJSON(&food); err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		
		validationErr :=	validate.Struct(food)
		if validationErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		err := menuCollection.FindOne(ctx,bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()
		if err !=nil{
			msg:=fmt.Sprintf("menu not found")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		food.Created_at, _ =	time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		food.ID =primitive.NewObjectID()
		food.Food_id=food.ID.Hex()
		var num = toFixed(*food.Price,2)
		food.Price = &num
		
		result,insertErr := foodCollection.InsertOne(ctx,food)
		if insertErr !=nil{
			msg:=fmt.Sprintf("food item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func round(num float64) int {
	return int(num +math.Copysign(0.5,num))
}

func toFixed(num float64,precision int)float64{
	output := math.Pow(10,float64(precision))
	return float64(round(num*output))/output
}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		var menu model.Menu
		var food model.Food
		foodId := c.Param("food_id")

		if err:=c.BindJSON(&food);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		var updateObj primitive.D
		if food.Name!=nil{
			updateObj = append(updateObj,bson.E{"name",food.Name})
		}
		if food.Price!=nil{
			updateObj = append(updateObj,bson.E{"price",food.Price})
		}
		if food.Food_image!=nil{
			updateObj = append(updateObj,bson.E{"food_image", food.Food_image})
		}
		if food.Menu_id!=nil{
			err := menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
			defer cancel()
			if err !=nil{
				msg := fmt.Sprintf("message:menu not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
			updateObj = append(updateObj,bson.E{"menu", food.Price})
		}

		food.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj,bson.E{"created_at",food.Updated_at})
		
		upsert := true
		filter := bson.M{"food_id":foodId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		defer cancel()

		if err!=nil{
			msg:= fmt.Sprintf("Food item Update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

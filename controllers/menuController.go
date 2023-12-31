package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/kushagra-gupta01/Restaurant-Management/database"
	"github.com/kushagra-gupta01/Restaurant-Management/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		result,err := menuCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing the menu items"})
			var allMenus []bson.M
			if err = result.All(ctx, &allMenus);err!=nil{
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, allMenus)
		}
	}
}


func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		menuId:=c.Param("menu_id")
		var menu model.Menu

		err:= menuCollection.FindOne(ctx,bson.M{"menu_id":menuId}).Decode(&menu)
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while fetching the menu"})
		}
		c.JSON(http.StatusOK,menu)
	}
}

func CreateMenu()  gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var menu model.Menu

		if err:= c.BindJSON(&menu);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}
		defer cancel()
		menu.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		menu.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result,insertErr := menuCollection.InsertOne(ctx,menu)
		if insertErr !=nil{
			msg:=fmt.Sprintf("menu was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func inTimeSpan(start,end, check time.Time) bool{
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		var menu model.Menu
		if err:= c.BindJSON(&menu);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter :=bson.M{"menu_id":menuId}

		var updateObj primitive.D
		if menu.Start_date !=nil && menu.End_date !=nil{
			if !inTimeSpan(*menu.Start_date,*menu.End_date,time.Now()){
				msg := "kindly re-type the time"
				c.JSON(http.StatusInternalServerError,gin.H{"error": msg})
				defer cancel()
				return 
			}
			updateObj = append(updateObj, bson.E{"start_date", menu.Start_date})
			updateObj = append(updateObj, bson.E{"start_date", menu.End_date})

		}
		if menu.Name != ""{
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}	
		if menu.Category !=""{
			updateObj = append(updateObj, bson.E{"category", menu.Category})
		}
		menu.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"created_at", menu.Updated_at})
		
		upsert := true
		opt:= options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		if err !=nil{
			msg := "Menu update failed"
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}
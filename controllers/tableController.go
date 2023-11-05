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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result,err :=tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listing tables"})
		}
		var allTables[]bson.M
		if err = result.All(ctx, &allTables);err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allTables)
	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		tableId := c.Param("table_id")
		var table model.Table

		err:= tableCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&table)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while fetching the table"})
		}
		c.JSON(http.StatusOK,table)
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		var table model.Table

		if err := c.BindJSON(&table); err!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error":err.Error()})
			return
		}

		validationErr := validate.Struct(table)

		if validationErr!= nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error":validationErr.Error()})
			return
		}

		table.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result,InsertErr := tableCollection.InsertOne(ctx,table)
		if InsertErr != nil{
			msg := fmt.Sprintf("Table was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return 
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		var table model.Table

		tableId := c.Param("table_id")

		if err := c.BindJSON(&table);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		var updateObj primitive.D

		if table.Number_of_guests !=nil{
			updateObj = append(updateObj, bson.E{"number_of_guests",table.Number_of_guests})
		}

		if table.Table_number !=nil{
			updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
		}

		table.Updated_at ,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		upsert := true
		opt:= options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		result,err :=tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)

		if err!=nil{
			msg := fmt.Sprintf("Failed to update the table")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/sardorgo/myapp/proto"
)

var userConnection pb.UserProfilesClient
var lisConnection pb.ListUsersRequest

func PostUser(c *gin.Context) {
	var newUser pb.UserProfile
	if err := c.ShouldBindJSON(&newUser); err != nil {
		return
	}

	req := &pb.CreateUserProfileRequest{
		UserProfile: &pb.UserProfile{
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Id:        newUser.Id,
		},
	}
	res, err := userConnection.CreateUser(context.Background(), req)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, res)
}

func GetUser(c *gin.Context) {
	id := c.Param("userId")
	req := &pb.GetUserProfileRequest{
		UserId: id,
	}
	res, err := userConnection.GetUser(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	c.IndentedJSON(http.StatusOK, res)
}

func PutUser(c *gin.Context) {
	var updateWanted pb.UserProfile

	if err := c.ShouldBindJSON(&updateWanted); err != nil {
		return
	}
	req := &pb.UpdateUserProfileRequest{
		UserProfile: &pb.UserProfile{
			FirstName: updateWanted.FirstName,
			LastName:  updateWanted.LastName,
			Id:        updateWanted.Id,
		},
	}
	res, err := userConnection.UpdateUser(context.Background(), req)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("userId")
	req := &pb.DeleteUserProfileRequest{
		UserId: id,
	}
	res, err := userConnection.DeleteUser(context.Background(), req)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func ListAll(c *gin.Context) {
	req := &pb.ListUsersRequest{}
	res, err := userConnection.ListUsers(context.Background(), req)

	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, res)
}

func main() {
	fmt.Println("Welcome Client")
	conn, err := grpc.Dial("localhost:9500", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect %v", err)
	}

	defer conn.Close()

	userConnection = pb.NewUserProfilesClient(conn)

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	router := gin.Default()

	router.POST("/users", PostUser)
	router.DELETE("/users/:userId", DeleteUser)
	router.PUT("/users", PutUser)
	router.GET("/users/:userId", GetUser)
	router.GET("/users", ListAll)

	router.Run("localhost:5000")
}

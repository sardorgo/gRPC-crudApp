package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/satori/uuid"

	"google.golang.org/grpc"

	pb "github.com/sardorgo/myapp/proto"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "sardor"
	PASSWORD = "sardor"
	DBNAME   = "grpcapp"
)

type server struct {
	conn *sql.DB
	pb.UnimplementedUserProfilesServer
}

func (connection *server) CreateUser(ctx context.Context, req *pb.CreateUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	id := uuid.NewV4()
	req.UserProfile.Id = id.String()

	firstName := req.GetUserProfile().GetFirstName()
	lastName := req.GetUserProfile().GetLastName()

	sqlInsert := `insert into users (user_id, first_name, last_name) values ($1, $2, $3);`

	if _, err := db.Exec(sqlInsert, id, firstName, lastName); err != nil {
		return nil, errors.Wrapf(err, "User couldn't be inserted")
	}

	return req.UserProfile, nil

}

func (connection *server) GetUser(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	id := req.GetUserId()
	sqlStatement := `select * from users where user_id = $1`
	var first_name, last_name, user_id string
	// var books []*pb.Book
	err := db.QueryRow(sqlStatement, id).Scan(&user_id, &first_name, &last_name)

	if err != nil {
		errors.Wrapf(err, "User profile couldn't be returned")
	}

	res := &pb.UserProfile{
		Id:        id,
		FirstName: first_name,
		LastName:  last_name,
	}
	return res, nil
}

func (connection *server) UpdateUser(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	sqlStatement := `update users set  first_name=$2, last_name=$3 where user_id=$1`
	if _, err := db.Exec(sqlStatement, req.UserProfile.Id, req.UserProfile.FirstName, req.UserProfile.LastName); err != nil {
		return nil, err
	}

	return req.UserProfile, nil
}

func (connection *server) DeleteUser(ctx context.Context, req *pb.DeleteUserProfileRequest) (*pb.Empty, error) {
	db := connection.conn
	sqlStatement := `delete from users where user_id = $1`

	id := req.GetUserId()

	if _, err := db.Exec(sqlStatement, id); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (connection *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	db := connection.conn
	sqlStatement := `select user_id, first_name, last_name from users`
	result, err := db.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
	}

	defer result.Close()

	res := []*pb.UserProfile{}
	for result.Next() {
		var first, last, id string
		if err = result.Scan(&id, &first, &last); err != nil {
			errors.Wrap(err, "Users couln't be listed")
		}
		u := pb.UserProfile{
			Id:        id,
			FirstName: first,
			LastName:  last,
		}
		res = append(res, &u)
	}
	ans := pb.ListUsersResponse{Profiles: res}
	return &ans, nil
}

func main() {
	fmt.Println("Welcome to the server")
	lis, err := net.Listen("tcp", ":9500")

	if err != nil {
		errors.Wrap(err, "UserProfile couldn't be returned")
	}

	s := grpc.NewServer()

	conn, err := grpc.Dial("localhost: 9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errors.Wrap(err, "UserProfile couldn't be returned")
	}
	defer db.Close()
	err = db.Ping()

	if err != nil {
		errors.Wrap(err, "User couldn't be listed")
	}

	pb.RegisterUserProfilesServer(s, &server{conn: db})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "UserProfile couldn't be returned")
	}
}

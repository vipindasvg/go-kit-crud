package user

import (
	"fmt"
	"context"
	"errors"

	"github.com/couchbase/gocb"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
)

var RepoErr = errors.New("Unable to handle Repo Request")

type repo struct {
	logger log.Logger
}

func NewRepo(logger log.Logger) Repository {
	return &repo{
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	cluster, err := gocb.Connect("couchbase://localhost")
	fmt.Println(err)
    cluster.Authenticate(gocb.PasswordAuthenticator{
        Username: "Administrator",
        Password: "123456",
    })

    bucket, err := cluster.OpenBucket("user", "")
	fmt.Println(err)
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	users := User{Id:user.Id, Name :user.Name, EmailId:user.EmailId, Password : string(bytes)}

	cas, _ := bucket.Insert(user.Id, &users, 0)
	fmt.Println(cas)
	return nil
}

func (repo *repo) UserLogin(ctx context.Context, EmailId string, Password string) (*User, error) {
	cluster, err := gocb.Connect("couchbase://localhost")
	fmt.Println(err)
    cluster.Authenticate(gocb.PasswordAuthenticator{
        Username: "Administrator",
        Password: "123456",
    })

    bucket, err := cluster.OpenBucket("user", "")
	fmt.Println(err)
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `user` WHERE email_id=$1 ;")
	var myParams []interface{}
	myParams = append(myParams, EmailId)
	//fmt.Println(myParams)
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	var row map[string]interface{}
	user := new(User)
	//var retValues []interface{}
	// Stream the values returned from the query into an untyped and unstructred
	// array of interfaces
	fmt.Println("fnte",row)
	err = rows.One(&row)
	fmt.Println("fnte2",row["user"])
	fmt.Println(err)
	jsonOut , err := json.Marshal(row["user"])
	err1 := json.Unmarshal(jsonOut, &user)
	fmt.Println(err1)
	err2 := bcrypt.CompareHashAndPassword([]byte(Password), []byte(user.Password))
	fmt.Println(err2)
	if err2 != nil {
		return nil, nil
	} else {
		return user, nil	
	}
	return nil, nil	
}

func (repo *repo) ListUsers(ctx context.Context) ([]*User, error) {
	cluster, err := gocb.Connect("couchbase://localhost")
	fmt.Println(err)
    cluster.Authenticate(gocb.PasswordAuthenticator{
        Username: "Administrator",
        Password: "123456",
    })

    bucket, err := cluster.OpenBucket("user", "")
	fmt.Println(err)
	query := "SELECT * FROM `user`;"
	rows, err := cluster.Query(query, &gocb.QueryOptions{})
	// check query was successful
	if err != nil {
		panic(err)
	}
	var users []User
	// iterate over rows
	for rows.Next() {
		var u User// this could also just be an interface{} type
		err := rows.Row(&u)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}	
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return users, nil
}
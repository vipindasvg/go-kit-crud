package user

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/dgrijalva/jwt-go"
)


//Create custom jwt token
type UserClaims struct {
	jwt.StandardClaims
	UserName string `json:"user_name"`
	UserId   string `json:"user_id"`
}

// CreateRequest holds the request parameters for the Create method.
type CreateUserRequest struct {
	user *User	`json:"user"`
}

// CreateResponse holds the response values for the Create method.
type CreateUserResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// CreateRequest holds the request parameters for the Create method.
type UserLoginResponse struct {
	user *User
	Err error  `json:"error,omitempty"`
	Status string `json:"status"`
	Token  string `json:"token"`
}

// CreateResponse holds the response values for the Create method.
type UserLoginRequest struct {
	EmailId		string	`json:"email_id"`
	Password	string	`json:"password"`
}

// CreateRequest holds the request parameters for the Create method.
type ListUsersResponse struct {
	users []*User `json:"users"`  
	Err error  `json:"error,omitempty"`
}

// CreateResponse holds the response values for the Create method.
type ListUsersRequest struct {
}

// Endpoints holds all Go kit endpoints for the user service.
type Endpoints struct {
	CreateUser      endpoint.Endpoint
	UserLogin		endpoint.Endpoint
	ListUsers		endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:       makeCreateUserEndpoint(s),
		UserLogin:        makeUserLoginEndpoint(s),
		ListUsers:		  makeListUsersEndpoint(s),
	}
}

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest) // type assertion
		fmt.Println(req)
		id, err := s.CreateUser(ctx, req.user)
		return CreateUserResponse{ID: id, Err: err}, nil
	}
}

func makeUserLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserLoginRequest) // type assertion
		fmt.Println(req)
		users, err := s.UserLogin(ctx, req.EmailId, req.Password)
		fmt.Println(err)
		if users != nil {
			secretKey := "secret"
			tokenDuration := 15 * time.Minute
			claims := UserClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(tokenDuration).Unix(),
				},
				UserName: users.Name,
				UserId: users.Id,
			}
		
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			finalToken, _ := token.SignedString([]byte(secretKey))
			fmt.Println(finalToken)
			return UserLoginResponse{user: users, Token: finalToken, Status: "Login Success"}, nil	
		} else {
			return UserLoginResponse{user: nil, Token: " ", Status: "Login Failed"}, nil
		}
		return UserLoginResponse{user: users}, nil
	}
}

func makeListUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		lusers, err := s.ListUsers(ctx)
		return ListUsersResponse{users: lusers, Err: err}, nil
	}
}

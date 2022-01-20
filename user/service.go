package user

import (
	"context"
	"errors"
	"fmt"
 
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

type Service interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	UserLogin(ctx context.Context, EmailId string, Password string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}

// service implements the Order Service
type service struct {
	repository Repository
	logger     log.Logger
 }
 
 // NewService creates and returns a new Order service instance
 func NewService(rep Repository, logger log.Logger) Service {
	return &service{
	   repository: rep,
	   logger:     logger,
	}
 }
 
 // Create makes an order
 func (s *service) CreateUser(ctx context.Context, user *User) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	user.Id = "user_" + id
 
	if err := s.repository.CreateUser(ctx, *user); err != nil {
	   level.Error(logger).Log("err", err)
	   return "", ErrCmdRepository
	}
	return id, nil
 }

 // User Login
 func (s *service) UserLogin(ctx context.Context, EmailId string, Password string) (*User, error) {
	logger := log.With(s.logger, "method", "Login")
	user, err := s.repository.UserLogin(ctx, EmailId, Password);
	fmt.Println(user)
	if err != nil {
	   level.Error(logger).Log("err", err)
	   return nil, ErrCmdRepository
	}
	return user, nil
 }


 func (s *service) ListUsers(ctx context.Context) ([]*User, error) {
	logger := log.With(s.logger, "method", "list")
	if users, err := s.repository.ListUsers(ctx); err != nil {
	   level.Error(logger).Log("err", err)
	   return nil, ErrCmdRepository
	}
	return users, nil
 }
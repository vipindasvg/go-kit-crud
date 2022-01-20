package user

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//"github.com/shijuvar/gokit-examples/services/account/transport"
	//"github.com/shijuvar/gokit-examples/services/account/transport/pb"
	"github.com/vipindasvg/go-kit-crud/user/pb"
)

// grpc transport service for Account service.
type grpcServer struct {
	createUser kitgrpc.Handler
	userLogin  kitgrpc.Handler
	logger         log.Logger
	pb.UnimplementedUserServiceServer
}

// NewGRPCServer returns a new gRPC service for the provided Go kit endpoints
func NewGRPCServer(
	endpoints Endpoints, options []kitgrpc.ServerOption,
	logger log.Logger,
) pb.UserServiceServer {
	errorLogger := kitgrpc.ServerErrorLogger(logger)
	options = append(options, errorLogger)

	return &grpcServer{
		createUser: kitgrpc.NewServer(
			endpoints.CreateUser, decodeCreateUserRequestGrpc, encodeCreateUserResponseGrpc, options...,
		),
		userLogin: kitgrpc.NewServer(
			endpoints.UserLogin, decodeUserLoginRequestGrpc, encodeUserLoginResponseGrpc, options...,
		),
		listUsers: kitgrpc.NewServer(
			endpoints.ListUsers, decodeListUsersRequestGrpc, encodeListUsersResponseGrpc, options...,
		),
		logger: logger,
	}
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) CreateUser(ctx oldcontext.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, rep, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserResponse), nil
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) UserLogin(ctx oldcontext.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	_, rep, err := s.userLogin.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserLoginResponse), nil
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) ListUser(ctx oldcontext.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, rep, err := s.listUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserResponse), nil
}

// decodeCreateCustomerRequest decodes the incoming grpc payload to our go kit payload
func decodeCreateUserRequestGrpc(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateUserRequest)
	users:= User{ Id:req.User.Id, Name: req.User.Name, Password: req.User.Password, EmailId:req.User.EmailId}
	return CreateUserRequest{
		user : &users,
	}, nil
}

// encodeCreateCustomerResponse encodes the outgoing go kit payload to the grpc payload
func encodeCreateUserResponseGrpc(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(CreateUserResponse)
	err := getError(res.Err)
	if err == nil {
		return &pb.CreateUserResponse{}, nil
	}
	return nil, err
}

// decodeCreateCustomerRequest decodes the incoming grpc payload to our go kit payload
func decodeListUsersRequestGrpc(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListUsersRequest)
	fmt.Println(req)
	return ListUsersRequest{
	}, nil
}

// encodeCreateCustomerResponse encodes the outgoing go kit payload to the grpc payload
func encodeUserLoginResponseGrpc(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(UserLoginResponse)
	err := getError(res.Err)
	pbuser := new(pb.User)
	pbuser.Id = res.user.Id
	pbuser.Name = res.user.Name
	pbuser.EmailId = res.user.EmailId
	if err == nil {
		fmt.Println("enc",res)
		return &pb.UserLoginResponse{User: pbuser, Token: res.Token, Status: res.Status}, nil
	}
	return res, err
}

// decodeCreateCustomerRequest decodes the incoming grpc payload to our go kit payload
func decodeListUsersRequestGrpc(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListUsersRequest)
	users:= User{ Id:req.User.Id, Name: req.User.Name, Password: req.User.Password, EmailId:req.User.EmailId}
	return CreateUserRequest{
		user : &users,
	}, nil
}

// encodeCreateCustomerResponse encodes the outgoing go kit payload to the grpc payload
func encodeListUsersResponseGrpc(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(CreateUserResponse)
	err := getError(res.Err)
	if err == nil {
		return &pb.CreateUserResponse{}, nil
	}
	return nil, err
}

func getError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

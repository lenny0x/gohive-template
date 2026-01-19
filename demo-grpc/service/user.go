package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/gohive/demo-grpc/proto/userpb"
	"github.com/gohive/models/entity"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(server *grpc.Server) {
	userpb.RegisterUserServiceServer(server, s)
}

func (s *UserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	var user entity.User
	if err := s.db.First(&user, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &userpb.GetUserResponse{
		User: toProtoUser(&user),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []entity.User
	var total int64

	if err := s.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count users: %v", err)
	}

	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	protoUsers := make([]*userpb.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(&user)
	}

	return &userpb.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userpb.CreateUserResponse{
		User: toProtoUser(user),
	}, nil
}

func toProtoUser(user *entity.User) *userpb.User {
	return &userpb.User{
		Id:        int64(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

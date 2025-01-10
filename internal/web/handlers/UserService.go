package handlers

import (
	"context"
	"github.com/dyingvoid/pigeon-server/internal/application/models"
	"github.com/dyingvoid/pigeon-server/internal/mongodb/repositories"
	pb "github.com/dyingvoid/pigeon-server/internal/web/proto"
)

type UserService struct {
	rep *repositories.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(rep *repositories.UserRepository) *UserService {
	return &UserService{rep: rep}
}

func (s *UserService) AddUser(
	ctx context.Context, req *pb.AddUserRequest,
) (*pb.AddUserResponse, error) {
	user, err := models.NewUser(req.Username, req.PublicKey)
	if err != nil {
		return nil, err
	}

	err = s.rep.Add(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.AddUserResponse{Message: "great success"}, nil
}

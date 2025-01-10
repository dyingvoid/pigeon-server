package handlers

import (
	"context"
	"github.com/dyingvoid/pigeon-server/internal/web/authentication"
	pb "github.com/dyingvoid/pigeon-server/internal/web/proto"
)

type ChallengeService struct {
	auth *authentication.Authentication
	pb.UnimplementedChallengeServiceServer
}

func NewChallengeService(auth *authentication.Authentication) *ChallengeService {
	return &ChallengeService{auth: auth}
}

func (s *ChallengeService) GetChallenge(
	ctx context.Context, req *pb.ChallengeRequest,
) (*pb.ChallengeResponse, error) {
	challenge, err := s.auth.CreateChallenge(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ChallengeResponse{
		Nonce: challenge.Nonce,
	}, nil
}

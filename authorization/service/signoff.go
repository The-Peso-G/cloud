package service

import (
	"context"

	"github.com/go-ocf/cloud/authorization/pb"
	"github.com/go-ocf/kit/log"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SignOff invalidates device's Access Token.
func (s *Service) SignOff(ctx context.Context, request *pb.SignOffRequest) (*pb.SignOffResponse, error) {
	tx := s.persistence.NewTransaction(ctx)
	defer tx.Close()

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign off: %v", err))
	}
	userID, err := parseSubFromJwtToken(token)
	if err != nil {
		log.Debugf("cannot parse user from jwt token: %v", err)
		userID = request.UserId
	}

	if userID == "" {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign off: invalid UserId"))
	}

	d, ok, err := tx.Retrieve(request.DeviceId, userID)
	if err != nil {
		return nil, logAndReturnError(status.Errorf(codes.Internal, "cannot sign off: %v", err.Error()))
	}
	if !ok {
		return nil, logAndReturnError(status.Errorf(codes.NotFound, "cannot sign off: not found"))
	}
	if d.AccessToken != token {
		return nil, logAndReturnError(status.Errorf(codes.InvalidArgument, "cannot sign off: unexpected access token"))
	}

	err = tx.Delete(d.DeviceID, userID)
	if err != nil {
		return nil, logAndReturnError(status.Errorf(codes.Internal, "cannot sign off: %v", err.Error()))
	}

	return &pb.SignOffResponse{}, nil
}

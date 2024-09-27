package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/proto"
)

type ShortURLService struct {
	proto.UnimplementedShortURLServiceServer
	store  contracts.Storage
	zapLog *zap.SugaredLogger
}

func NewShortURLService(store contracts.Storage, zapLog *zap.SugaredLogger) *ShortURLService {
	return &ShortURLService{
		store:  store,
		zapLog: zapLog,
	}
}

func (sh *ShortURLService) InternalStats(ctx context.Context, _ *proto.InternalStatsRequest,
) (*proto.InternalStatsResponse, error) {
	sh.zapLog.Infoln("GRPC InternalStats called")
	var responseStats proto.InternalStatsResponse
	stats, err := sh.store.InternalStats(ctx)
	if err != nil {
		responseStats.Error = err.Error()
	}
	responseStats.Users = stats.Users
	responseStats.Urls = stats.Urls
	return &responseStats, nil
}

func (sh *ShortURLService) DeleteUserURLs(ctx context.Context, in *proto.DeleteUserURLsRequest,
) (*proto.DeleteUserURLsResponse, error) {
	sh.zapLog.Infoln("GRPC DeleteUserURLs called")
	var response proto.DeleteUserURLsResponse
	err := sh.store.DeleteURLs(ctx, in.Urls)
	if err != nil {
		response.Status = "no accepted"
		response.Error = fmt.Errorf("error deleting urls: %w", err).Error()
		return &response, err
	}

	response.Status = "accepted"
	return &response, nil
}

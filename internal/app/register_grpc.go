package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/peer"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/proto"
	"mmskazak/shorturl/internal/services/checkip"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/shorturlservice"
	"net"
)

type ShortURLService struct {
	proto.UnimplementedShortURLServiceServer
	cfg    *config.Config
	store  contracts.Storage
	zapLog *zap.SugaredLogger
}

func NewShortURLService(cfg *config.Config, store contracts.Storage, zapLog *zap.SugaredLogger) *ShortURLService {
	return &ShortURLService{
		cfg:    cfg,
		store:  store,
		zapLog: zapLog,
	}
}

func (sh *ShortURLService) InternalStats(ctx context.Context, _ *proto.InternalStatsRequest,
) (*proto.InternalStatsResponse, error) {
	sh.zapLog.Infoln("GRPC InternalStats called")
	var responseStats proto.InternalStatsResponse

	// Извлечение информации о peer (клиенте)
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("access forbiden")
	}

	addr := p.Addr
	sh.zapLog.Infof("Request from IP: %v", addr)

	// извлекаем только IP-адрес
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return nil, fmt.Errorf("invalid tcp address")
	}
	clientIP := tcpAddr.IP.String()
	sh.zapLog.Infof("Client IP: %s", clientIP)

	ok, err := checkip.CheckIPByCIDR(clientIP, sh.cfg.TrustedSubnet)
	if err != nil {
		return nil, fmt.Errorf("error checking ip address: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("access forbiden")
	}

	stats, err := sh.store.InternalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get internal stats: %w", err)
	}
	responseStats.Users = stats.Users
	responseStats.Urls = stats.Urls
	return &responseStats, nil
}

func (sh *ShortURLService) DeleteUserURLs(
	ctx context.Context,
	in *proto.DeleteUserURLsRequest,
) (*proto.DeleteUserURLsResponse, error) {
	sh.zapLog.Infoln("GRPC DeleteUserURLs called")
	var response proto.DeleteUserURLsResponse
	err := sh.store.DeleteURLs(ctx, in.Urls)
	if err != nil {
		response.Status = "not accepted"
		return nil, fmt.Errorf("error deleting urls: %w", err)
	}

	response.Status = "accepted"
	return &response, nil
}

func (sh *ShortURLService) FindUserURLs(
	ctx context.Context,
	in *proto.FindUserURLsRequest,
) (*proto.FindUserURLsResponse, error) {
	sh.zapLog.Infoln("GRPC FindUserURLs called")
	var response proto.FindUserURLsResponse

	urls, err := sh.store.GetUserURLs(ctx, in.UserId, sh.cfg.BaseHost)
	if err != nil {
		return nil, fmt.Errorf("error getting user urls: %w", err)
	}

	// Преобразуем результаты в слайс структур UserURLs
	for _, url := range urls {
		userURL := &proto.UserURLs{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
		}
		response.UserUrls = append(response.UserUrls, userURL)
	}

	return &response, nil
}

func (sh *ShortURLService) SaveShortenURLsBatch(
	ctx context.Context,
	in *proto.SaveShortenURLsBatchRequest,
) (*proto.SaveShortenURLsBatchResponse, error) {
	sh.zapLog.Infoln("GRPC SaveShortenURLsBatch called")
	var response proto.SaveShortenURLsBatchResponse

	// Преобразуем []*Incoming в []models.Incoming
	incomingModels := make([]models.Incoming, len(in.Incoming))

	// Сохранение пакета коротких URL
	for i, inc := range in.Incoming {
		incomingModels[i] = models.Incoming{
			// Здесь копируем поля структуры
			CorrelationID: inc.CorrelationId,
			OriginalURL:   inc.OriginalUrl,
		}
	}
	generator := genidurl.NewGenIDService()
	outputs, err := sh.store.SaveBatch(ctx, incomingModels, sh.cfg.BaseHost, in.UserId, generator)
	if err != nil {
		return nil, fmt.Errorf("error saving shorten urls: %w", err)
	}

	// Преобразуем outputs в слайс Output для ответа
	for _, output := range outputs {
		out := &proto.Output{
			CorrelationId: output.CorrelationID,
			ShortUrl:      output.ShortURL,
		}
		response.Incoming = append(response.Incoming, out)
	}

	return &response, nil
}

func (sh *ShortURLService) HandleCreateShortURL(
	ctx context.Context,
	in *proto.HandleCreateShortURLRequest,
) (*proto.HandleCreateShortURLResponse, error) {
	sh.zapLog.Infoln("GRPC HandleCreateShortURL called")
	var response proto.HandleCreateShortURLResponse

	generator := genidurl.NewGenIDService()
	dto := dtos.DTOShortURL{
		UserID:      in.UserId,
		OriginalURL: in.OriginalUrl,
		BaseHost:    sh.cfg.BaseHost,
		Deleted:     false,
	}

	shortURLService := shorturlservice.NewShortURLService()
	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, sh.store)
	if err != nil {
		return nil, fmt.Errorf("error creating shorten url: %w", err)
	}
	response.Result = shortURL
	return &response, nil
}

package app

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"

	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/contracts"
	"mmskazak/shorturl/internal/dtos"
	"mmskazak/shorturl/internal/models"
	"mmskazak/shorturl/internal/proto"
	"mmskazak/shorturl/internal/services/checkip"
	"mmskazak/shorturl/internal/services/genidurl"
	"mmskazak/shorturl/internal/services/jwttoken"
	"mmskazak/shorturl/internal/services/shorturlservice"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/peer"
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
	sh.zapLog.Debugln("GRPC InternalStats called")
	var responseStats proto.InternalStatsResponse

	// Извлечение информации о peer (клиенте)
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("access forbidden")
	}

	addr := p.Addr
	sh.zapLog.Debugf("Request from IP: %v\n", addr)

	// извлекаем только IP-адрес
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return nil, errors.New("invalid tcp address")
	}
	clientIP := tcpAddr.IP.String()
	sh.zapLog.Debugln("Client IP: %s", clientIP)

	ok, err := checkip.CheckIPByCIDR(clientIP, sh.cfg.TrustedSubnet)
	if err != nil {
		return nil, fmt.Errorf("error checking ip address: %w", err)
	}
	if !ok {
		return nil, errors.New("access forbidden")
	}

	stats, err := sh.store.InternalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get internal stats: %w", err)
	}
	responseStats.Users = wrapperspb.String(stats.Users)
	responseStats.Urls = wrapperspb.String(stats.Urls)
	return &responseStats, nil
}

func (sh *ShortURLService) DeleteUserURLs(
	ctx context.Context,
	in *proto.DeleteUserURLsRequest,
) (*proto.DeleteUserURLsResponse, error) {
	sh.zapLog.Debugln("GRPC DeleteUserURLs called")
	var response proto.DeleteUserURLsResponse

	// Преобразование []*Urls в []string
	urls := in.GetUrls()
	urlStrings := make([]string, len(urls))

	for i, u := range urls {
		urlStrings[i] = u.Urls.GetValue()
	}

	// Передаем преобразованный слайс []string в DeleteURLs
	err := sh.store.DeleteURLs(ctx, urlStrings)
	if err != nil {
		response.Status = &proto.Status{Status: wrapperspb.String("not accepted")}
		return nil, fmt.Errorf("error deleting urls: %w", err)
	}

	response.Status = &proto.Status{Status: wrapperspb.String("accepted")}
	return &response, nil
}

func (sh *ShortURLService) FindUserURLs(
	ctx context.Context,
	in *proto.FindUserURLsRequest,
) (*proto.FindUserURLsResponse, error) {
	sh.zapLog.Debugln("GRPC FindUserURLs called")
	var response proto.FindUserURLsResponse

	urls, err := sh.store.GetUserURLs(ctx, in.GetUserId().UserId.GetValue(), sh.cfg.BaseHost)
	if err != nil {
		return nil, fmt.Errorf("error getting user urls: %w", err)
	}

	// Преобразуем результаты в слайс структур UserURLs
	for _, url := range urls {
		userURL := &proto.UserURLs{
			ShortUrl:    wrapperspb.String(url.ShortURL),
			OriginalUrl: wrapperspb.String(url.OriginalURL),
		}
		response.UserUrls = append(response.UserUrls, userURL)
	}

	return &response, nil
}

func (sh *ShortURLService) SaveShortenURLsBatch(
	ctx context.Context,
	in *proto.SaveShortenURLsBatchRequest,
) (*proto.SaveShortenURLsBatchResponse, error) {
	sh.zapLog.Debugln("GRPC SaveShortenURLsBatch called")
	var response proto.SaveShortenURLsBatchResponse
	jwtString, err := sh.getOrCreateJWTToken(in.GetJwt().GetValue())
	if err != nil {
		return nil, fmt.Errorf("error getting jwt token: %w", err)
	}

	UserID, err := jwttoken.GetUserIDFromJWT(jwtString, sh.cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("user id not found: %w", err)
	}

	// Преобразуем []*Incoming в []models.Incoming
	incomingModels := make([]models.Incoming, len(in.GetIncoming()))

	// Сохранение пакета коротких URL
	for i, inc := range in.GetIncoming() {
		incomingModels[i] = models.Incoming{
			// Здесь копируем поля структуры
			CorrelationID: inc.GetCorrelationId().GetValue(),
			OriginalURL:   inc.GetOriginalUrl().GetValue(),
		}
	}
	generator := genidurl.NewGenIDService()
	outputs, err := sh.store.SaveBatch(ctx, incomingModels, sh.cfg.BaseHost, UserID, generator)
	if err != nil {
		return nil, fmt.Errorf("error saving shorten urls: %w", err)
	}

	// Преобразуем outputs в слайс Output для ответа
	for _, output := range outputs {
		out := &proto.Output{
			CorrelationId: wrapperspb.String(output.CorrelationID),
			ShortUrl:      wrapperspb.String(output.ShortURL),
		}
		response.Output = append(response.Output, out)
	}

	return &response, nil
}

func (sh *ShortURLService) HandleCreateShortURL(
	ctx context.Context,
	in *proto.HandleCreateShortURLRequest,
) (*proto.HandleCreateShortURLResponse, error) {
	sh.zapLog.Debugln("GRPC HandleCreateShortURL called")
	var response proto.HandleCreateShortURLResponse

	jwtString, err := sh.getOrCreateJWTToken(in.GetJwt().GetValue())
	if err != nil {
		return nil, fmt.Errorf("error getting jwt token: %w", err)
	}

	UserID, err := jwttoken.GetUserIDFromJWT(jwtString, sh.cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("user id not found: %w", err)
	}

	generator := genidurl.NewGenIDService()
	dto := dtos.DTOShortURL{
		UserID:      UserID,
		OriginalURL: in.GetOriginalUrl().GetValue(),
		BaseHost:    sh.cfg.BaseHost,
		Deleted:     false,
	}

	shortURLService := shorturlservice.NewShortURLService()
	shortURL, err := shortURLService.GenerateShortURL(ctx, dto, generator, sh.store)
	if err != nil {
		return nil, fmt.Errorf("error creating shorten url: %w", err)
	}

	response.Result = wrapperspb.String(shortURL)
	return &response, nil
}

func (sh *ShortURLService) getOrCreateJWTToken(jwt string) (string, error) {
	var err error
	var jwtString string
	if jwt == "" {
		// Создаем новый userID
		userID := uuid.New().String()
		jwtString, err = jwttoken.CreateNewJWTToken(userID, sh.cfg.SecretKey)
		if err != nil {
			return "", fmt.Errorf("error creating jwt token: %w", err)
		}
	} else {
		jwtString = jwt
	}

	return jwtString, nil
}

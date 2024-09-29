package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mmskazak/shorturl/internal/config"
	"mmskazak/shorturl/internal/proto"
	"mmskazak/shorturl/internal/storage/inmemory"
	"net"
	"testing"
)

func TestInternalStats(t *testing.T) {
	// Создаем тестовый сервер
	grpcServer := grpc.NewServer()
	zapLog := zap.NewNop().Sugar()
	cfg := &config.Config{
		TrustedSubnet: "127.0.0.0/24",
	}
	store, err := inmemory.NewInMemory(zapLog)
	require.NoError(t, err)
	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	// Запускаем сервер в отдельной горутине
	listener, err := net.Listen("tcp", ":0") // Слушаем на случайном порту
	require.NoError(t, err)
	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	// Создаем клиент для обращения к серверу
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck //пренебрежем этим в тесте

	client := proto.NewShortURLServiceClient(conn)

	// Создание пустого запроса
	req := &proto.InternalStatsRequest{}

	// Вызов метода
	resp, err := client.InternalStats(context.Background(), req)
	require.NoError(t, err)

	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil
	assert.Equal(t, "0", resp.Urls)
	assert.Equal(t, "0", resp.Users)

	// Закрываем сервер
	grpcServer.Stop()
}

func TestDeleteUserURLs(t *testing.T) {
	// Создаем тестовый сервер
	grpcServer := grpc.NewServer()
	zapLog := zap.NewNop().Sugar()
	cfg := &config.Config{
		TrustedSubnet: "127.0.0.0/24",
	}
	store, err := inmemory.NewInMemory(zapLog)
	require.NoError(t, err)
	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	// Запускаем сервер в отдельной горутине
	listener, err := net.Listen("tcp", ":0") // Слушаем на случайном порту
	require.NoError(t, err)
	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	// Создаем клиент для обращения к серверу
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck //пренебрежем этим в тесте

	client := proto.NewShortURLServiceClient(conn)

	req := &proto.DeleteUserURLsRequest{
		Urls: []string{"http://example.com/1", "http://example.com/2"},
	}

	// Вызов метода
	resp, err := client.DeleteUserURLs(context.Background(), req)
	require.NoError(t, err)

	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil
	assert.Equal(t, "accepted", resp.GetStatus())

	// Закрываем сервер
	grpcServer.Stop()
}

func TestFindUserURLs(t *testing.T) {
	// Создаем тестовый сервер
	grpcServer := grpc.NewServer()
	zapLog := zap.NewNop().Sugar()
	cfg := &config.Config{
		TrustedSubnet: "127.0.0.0/24",
	}
	store, err := inmemory.NewInMemory(zapLog)
	require.NoError(t, err)
	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	// Запускаем сервер в отдельной горутине
	listener, err := net.Listen("tcp", ":0") // Слушаем на случайном порту
	require.NoError(t, err)
	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	// Создаем клиент для обращения к серверу
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck //пренебрежем этим в тесте

	client := proto.NewShortURLServiceClient(conn)

	req := &proto.FindUserURLsRequest{
		UserId: "1",
	}

	// Вызов метода
	resp, err := client.FindUserURLs(context.Background(), req)
	require.NoError(t, err)

	var expected []*proto.UserURLs
	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil
	assert.Equal(t, expected, resp.GetUserUrls())

	// Закрываем сервер
	grpcServer.Stop()
}

func TestSaveShortenURLsBatch(t *testing.T) {
	// Создаем тестовый сервер
	grpcServer := grpc.NewServer()
	zapLog := zap.NewNop().Sugar()
	cfg := &config.Config{
		TrustedSubnet: "127.0.0.0/24",
	}
	store, err := inmemory.NewInMemory(zapLog)
	require.NoError(t, err)
	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	// Запускаем сервер в отдельной горутине
	listener, err := net.Listen("tcp", ":0") // Слушаем на случайном порту
	require.NoError(t, err)
	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	// Создаем клиент для обращения к серверу
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck //пренебрежем этим в тесте

	client := proto.NewShortURLServiceClient(conn)
	req := &proto.SaveShortenURLsBatchRequest{
		Incoming: []*proto.Incoming{
			{
				CorrelationId: "1",
				OriginalUrl:   "http://example.com/1",
			},
			{
				CorrelationId: "2",
				OriginalUrl:   "http://example.com/2",
			},
		},
	}

	// Вызов метода
	resp, err := client.SaveShortenURLsBatch(context.Background(), req)
	require.NoError(t, err)

	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil

	// Закрываем сервер
	grpcServer.Stop()
}

func TestHandleCreateShortURL(t *testing.T) {
	// Создаем тестовый сервер
	grpcServer := grpc.NewServer()
	zapLog := zap.NewNop().Sugar()
	cfg := &config.Config{
		TrustedSubnet: "127.0.0.0/24",
	}
	store, err := inmemory.NewInMemory(zapLog)
	require.NoError(t, err)
	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	// Запускаем сервер в отдельной горутине
	listener, err := net.Listen("tcp", ":0") // Слушаем на случайном порту
	require.NoError(t, err)
	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	// Создаем клиент для обращения к серверу
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck //пренебрежем этим в тесте

	client := proto.NewShortURLServiceClient(conn)
	req := &proto.HandleCreateShortURLRequest{
		OriginalUrl: "http://example.com/1",
	}

	// Вызов метода
	resp, err := client.HandleCreateShortURL(context.Background(), req)
	require.NoError(t, err)

	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil

	// Закрываем сервер
	grpcServer.Stop()
}

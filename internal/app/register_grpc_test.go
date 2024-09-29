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

func TestGetData(t *testing.T) {
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

	// Проверка результата
	assert.NoError(t, err)
	assert.NotNil(t, resp) // Проверка, что ответ не nil
}

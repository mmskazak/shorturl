package app

import (
	"context"
	"errors"
	"fmt"
	"mmskazak/shorturl/internal/proto"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"mmskazak/shorturl/internal/contracts"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mmskazak/shorturl/internal/config"
)

// App представляет приложение с HTTP сервером и логгером.
type App struct {
	server         *http.Server
	zapLog         *zap.SugaredLogger
	grpcServer     *grpc.Server
	grpcServerAddr string
}

// ErrStartingServer - ошибка старта сервера.
const ErrStartingServer = "error starting server"

// NewApp создает новый экземпляр приложения.
// Ctx - контекст для управления временем выполнения.
// Cfg - конфигурация приложения.
// Store - хранилище данных.
// ReadTimeout - таймаут чтения HTTP-запросов.
// WriteTimeout - таймаут записи HTTP-ответов.
// ZapLog - логгер.
func NewApp(
	ctx context.Context,
	cfg *config.Config,
	store contracts.Storage,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	zapLog *zap.SugaredLogger,
	shortURLService contracts.IShortURLService,
) *App {
	router := chi.NewRouter()

	router = registrationMiddleware(router, cfg, zapLog)
	router = registrationWEBRoutes(
		ctx,
		router,
		cfg,
		zapLog,
		store,
		shortURLService,
	)

	router = registrationAPIRoutes(
		ctx,
		router,
		cfg,
		zapLog,
		store,
		shortURLService,
	)

	manager := &autocert.Manager{
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist("localhost"),
	}

	return &App{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		},
		zapLog:         zapLog,
		grpcServer:     newGRPCServer(cfg, store, zapLog),
		grpcServerAddr: ":50051",
	}
}

// Start запускает HTTP сервер.
func (a *App) Start() error {
	a.zapLog.Infof("Server is running on %v\n", a.server.Addr)

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.zapLog.Infof("%v: %v", ErrStartingServer, err)
		return fmt.Errorf(ErrStartingServer+": %w", err)
	}
	return nil
}

func newGRPCServer(cfg *config.Config, store contracts.Storage, zapLog *zap.SugaredLogger) *grpc.Server {
	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервисы
	proto.RegisterShortURLServiceServer(grpcServer, NewShortURLService(cfg, store, zapLog))

	return grpcServer
}

// StartGRPC запускает GRPC сервер.
func (a *App) StartGRPC() error {
	a.zapLog.Infof("gRPC Server is running on %v\n", a.grpcServerAddr)

	lis, err := net.Listen("tcp", a.grpcServerAddr)
	if err != nil {
		a.zapLog.Errorf("Failed to listen on %v: %v", a.grpcServerAddr, err)
		return fmt.Errorf("failed to listen on %v: %w", a.grpcServerAddr, err)
	}

	err = a.grpcServer.Serve(lis)
	if err != nil {
		a.zapLog.Errorf("Error starting gRPC server: %v", err)
		return fmt.Errorf("error starting gRPC server: %w", err)
	}
	return nil
}

// Stop корректно завершает работу приложения.
func (a *App) Stop(ctx context.Context) error {
	// Закрытие сервера с учетом переданного контекста.
	if err := a.server.Shutdown(ctx); err != nil {
		a.zapLog.Errorf("Ошибка при остановке сервера: %v", err)
		return fmt.Errorf("err Shutdown server: %w", err)
	}

	a.zapLog.Infoln("HTTP сервер успешно остановлен.")

	// Остановка gRPC сервера.
	a.grpcServer.GracefulStop() // gRPC не поддерживает Shutdown с контекстом, поэтому используем GracefulStop
	a.zapLog.Infoln("gRPC сервер успешно остановлен.")

	// Дополнительные действия по завершению работы (например, закрытие подключений к БД и т.д.)
	return nil
}

func (a *App) StartAll() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := a.Start(); err != nil {
			a.zapLog.Errorf("HTTP server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.StartGRPC(); err != nil {
			a.zapLog.Errorf("gRPC server error: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

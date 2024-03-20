package http

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Karzoug/share_bot/internal/delivery/http/middleware"
	"github.com/Karzoug/share_bot/internal/usecase/debt"
	"github.com/Karzoug/share_bot/internal/usecase/user"
	"github.com/vorlif/spreak"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg Config,
	userService user.Service,
	debtService debt.Service,
	localizer *spreak.Localizer,
	logger *zap.Logger) *server {
	logger = logger.With(zap.String("source", "http-server"))

	srv := &server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			ErrorLog:     zap.NewStdLog(logger),
			IdleTimeout:  defaultIdleTimeout,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		logger: logger,
	}

	handler := http.Handler(newHandler(userService, debtService, localizer, logger))
	handler = middleware.Logger(logger)(handler)
	if cfg.HeaderAuth.Key != "" {
		handler = middleware.HeaderAuth(cfg.HeaderAuth.Key, func(token string) bool {
			return subtle.ConstantTimeCompare([]byte(cfg.HeaderAuth.Value), []byte(token)) == 1
		})(handler)
	}
	handler = middleware.RecoverPanic(logger)(handler)
	srv.httpServer.Handler = handler

	return srv
}

func (s *server) Run(ctx context.Context) error {
	shutdownErrorChan := make(chan error)

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- s.httpServer.Shutdown(ctx)
	}()

	s.logger.Info("starting server", zap.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", zap.String("addr", s.httpServer.Addr))

	return nil
}

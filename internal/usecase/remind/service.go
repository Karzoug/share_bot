package remind

import (
	"context"
	"time"

	"github.com/vorlif/spreak"
	"go.uber.org/zap"

	"github.com/Karzoug/share_bot/internal/api"
	rrepo "github.com/Karzoug/share_bot/internal/usecase/remind/repo"
	urepo "github.com/Karzoug/share_bot/internal/usecase/user/repo"
)

type Service struct {
	cfg        Config
	api        api.API
	remindRepo rrepo.Repo
	userRepo   urepo.Repo
	t          *spreak.Localizer
	logger     *zap.Logger
}

func New(cfg Config,
	api api.API,
	remindRepo rrepo.Repo,
	userRepo urepo.Repo,
	localizer *spreak.Localizer,
	logger *zap.Logger) Service {
	return Service{
		cfg:        cfg,
		api:        api,
		remindRepo: remindRepo,
		userRepo:   userRepo,
		t:          localizer,
		logger:     logger,
	}
}

func (s Service) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.RunFrequency)
	defer ticker.Stop()

	bgFn := func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, s.cfg.RunFrequency)
		defer cancel()

		s.remind(ctx)
	}

	for {
		select {
		case <-ticker.C:
			bgFn(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

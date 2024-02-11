package debt

import (
	"time"

	"go.uber.org/zap"

	"github.com/vorlif/spreak"

	"github.com/Karzoug/share_bot/internal/api"
	drepo "github.com/Karzoug/share_bot/internal/usecase/debt/repo"
	urepo "github.com/Karzoug/share_bot/internal/usecase/user/repo"
)

var defaultBackgroundTaskTimeout = 3 * time.Second

type Service struct {
	api      api.API
	userRepo urepo.Repo
	debtRepo drepo.Repo
	t        *spreak.Localizer
	logger   *zap.Logger
}

func New(api api.API,
	userRepo urepo.Repo,
	debtRepo drepo.Repo,
	localizer *spreak.Localizer,
	logger *zap.Logger) Service {
	return Service{
		api:      api,
		userRepo: userRepo,
		debtRepo: debtRepo,
		t:        localizer,
		logger:   logger,
	}
}

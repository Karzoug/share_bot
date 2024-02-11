package app

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vorlif/spreak"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	_ "modernc.org/sqlite"

	"github.com/Karzoug/share_bot/internal/api"
	"github.com/Karzoug/share_bot/internal/config"
	"github.com/Karzoug/share_bot/internal/delivery/http"
	"github.com/Karzoug/share_bot/internal/usecase/debt"
	dRepo "github.com/Karzoug/share_bot/internal/usecase/debt/repo"
	"github.com/Karzoug/share_bot/internal/usecase/remind"
	rRepo "github.com/Karzoug/share_bot/internal/usecase/remind/repo"
	"github.com/Karzoug/share_bot/internal/usecase/user"
	uRepo "github.com/Karzoug/share_bot/internal/usecase/user/repo"
)

func Run(ctx context.Context, cfg *config.Config, logger *zap.Logger) error {
	db, err := sqlx.Open("sqlite", cfg.DB.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	api, err := api.New(ctx, cfg.API, logger)
	if err != nil {
		return err
	}

	ur := uRepo.New(db)
	dr := dRepo.New(db)
	rr := rRepo.New(db)

	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.Russian),
		spreak.WithDomainPath(spreak.NoDomain, "./locale"),
		spreak.WithLanguage(language.English, language.Russian),
	)
	if err != nil {
		return err
	}

	t := spreak.NewLocalizer(bundle, language.Russian)

	userService := user.New(api, ur, t, logger)
	debtService := debt.New(api, ur, dr, t, logger)
	remindService := remind.New(cfg.Remind, api, rr, ur, t, logger)

	eg, ctx := errgroup.WithContext(ctx)
	srv := http.New(cfg.HTTP, userService, debtService, t, logger)
	eg.Go(func() error {
		return srv.Run(ctx)
	})
	eg.Go(func() error {
		return remindService.Run(ctx)
	})

	return eg.Wait()
}

package cli

import (
	"context"
	"errors"
	"time"

	"github.com/bajankristof/goweb/auth"
	"github.com/bajankristof/goweb/config"
	"github.com/bajankristof/goweb/jwt"
	"github.com/bajankristof/goweb/oidc"
	"github.com/bajankristof/goweb/user"
	"github.com/bajankristof/goweb/wellknown"
)

type app struct {
	jwts    *jwt.Signer
	idps    oidc.Registry
	auth    *auth.Service
	user    *user.Service
	closers []func() error
}

func newApp(ctx context.Context, cfg *config.Config) (*app, error) {
	a := &app{}
	fail := func(err error) (*app, error) {
		return nil, errors.Join(err, a.close())
	}

	db, err := openDB(ctx, cfg)
	if err != nil {
		return fail(err)
	}
	a.closers = append(a.closers, db.Close)

	jwts, err := jwt.NewSigner(cfg.JWT.SigningKey.Unwrap())
	if err != nil {
		return fail(err)
	}
	a.jwts = jwts

	a.idps = oidc.NewRegistry()
	for n, p := range cfg.OIDC {
		a.idps.Add(n, oidc.NewProvider(
			p.IssuerURL.String(),
			p.ClientID,
			p.ClientSecret,
		))
	}

	a.auth = auth.NewService(
		auth.Stores{UserStore: db.userStore, SessionStore: db.sessionStore},
		db.authUoW,
		a.jwts,
		a.idps,
		auth.WithAccessTokenTTL(time.Duration(cfg.Auth.AccessTokenTTL)),
		auth.WithRefreshTokenTTL(time.Duration(cfg.Auth.RefreshTokenTTL)),
	)

	a.user = user.NewService(db.userStore)

	return a, nil
}

func (a *app) info() wellknown.Info {
	return wellknown.Info{
		Version: Version,
		Auth: wellknown.AuthInfo{
			Providers: a.idps.Names(),
		},
	}
}

func (a *app) close() error {
	var err error
	for i := len(a.closers) - 1; i >= 0; i-- {
		err = errors.Join(err, a.closers[i]())
	}
	a.closers = nil
	return err
}

package sqlstore

import "context"

type Txr[Tx any] interface {
	RunInTx(ctx context.Context, f func(Tx) error) error
}

type UnitOfWork[Tx any, S any] struct {
	txr     Txr[Tx]
	factory func(Tx) S
}

func NewUnitOfWork[Tx any, S any](txr Txr[Tx], factory func(Tx) S) *UnitOfWork[Tx, S] {
	return &UnitOfWork[Tx, S]{
		txr:     txr,
		factory: factory,
	}
}

func (u *UnitOfWork[Tx, S]) Do(ctx context.Context, f func(S) error) error {
	return u.txr.RunInTx(ctx, func(tx Tx) error {
		store := u.factory(tx)
		return f(store)
	})
}

package smooth

import (
	"context"
)

type Engine interface {
	First(context.Context, interface{}, Query) error
	Get(context.Context, interface{}, Query) error
	Create(context.Context, interface{}) error
	Update(context.Context, interface{}) error
	Delete(context.Context, interface{}) error
	Raw(context.Context, interface{}, Query) error
	Exec(context.Context, string) error
	Health() (bool, error)
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

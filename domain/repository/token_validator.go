package repository

import "context"

type TokenValidator interface {
	Validate(ctx context.Context, token string) (bool, error)
}

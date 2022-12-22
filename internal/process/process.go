package process

import (
	"context"
)

type Process interface {
	Update(data any)
	Run(ctx context.Context)
}

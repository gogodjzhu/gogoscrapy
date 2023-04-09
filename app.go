package gogoscrapy

import "context"

type IApp interface {
	Start(ctx context.Context)
}

package err

import "github.com/pkg/errors"

var RetryAbleError = errors.New("retryable error")

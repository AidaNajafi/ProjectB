package provider

import (
	"context"
	"errors"
)

type FailureReason string

const (
	FailAlreadyReserved FailureReason = "already_reserved"
	FailProviderTimeOut FailureReason = "provider_timeout"
	FailProviderDown    FailureReason = "provider_down"
	FailBadRequest      FailureReason = "bad_request"
	FailUnAuthorized    FailureReason = "unauthorized"
	FailUnknown         FailureReason = "unknown"
)

func ClassifyProviderError(err error) FailureReason {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.Canceled) {
		return FailUnknown
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return FailProviderTimeOut
	}
	var pe *ProviderError
	if errors.As(err, &pe) {
		switch pe.StatusCode {
		case 409:
			return FailAlreadyReserved
		case 400:
			return FailBadRequest
		case 401, 403:
			return FailUnAuthorized
		case 429:
			return FailProviderDown
		case 500, 502, 503, 504:
			return FailProviderDown
		default:
			return FailUnknown
		}
	}
	return FailUnknown
}


package service

import "context"

type Service interface {
	MakeRequest(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error)
}

package request_id

import "context"

type key struct{}

func NewContext(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, key{}, uuid)
}

func FromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	uuid, ok := ctx.Value(key{}).(string)
	return uuid, ok
}

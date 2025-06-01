package app

import (
	"context"
	"strings"
)

var CurrentUrlKey = "current_url"

func IsActiveUrl(ctx context.Context, url string, strict bool) bool {
	currentUrl, ok := ctx.Value(CurrentUrlKey).(string)
	if !ok {
		return false
	}

	if strict {
		return currentUrl == url
	}

	return strings.HasPrefix(currentUrl, url) || currentUrl == url
}

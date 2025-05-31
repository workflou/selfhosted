package app

import (
	"context"
	"selfhosted/database/store"
)

func GetUserFromContext(ctx context.Context) *store.User {
	sess, ok := ctx.Value(SessionKey).(store.GetSessionByUuidRow)
	if !ok || sess.ID == 0 {
		return nil
	}

	return &sess.User
}

func GetSessionFromContext(ctx context.Context) *store.GetSessionByUuidRow {
	sess, ok := ctx.Value(SessionKey).(store.GetSessionByUuidRow)
	if !ok || sess.ID == 0 {
		return nil
	}

	return &sess
}

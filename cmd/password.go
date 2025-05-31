package cmd

import (
	"context"
	"log/slog"
	"selfhosted/database"
	"selfhosted/database/store"

	"golang.org/x/crypto/bcrypt"
)

func ChangeAdminPassword(email, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	slog.Debug("Changing admin password", "email", email, "hash", string(hash))

	err = store.New(database.DB).ChangeAdminPassword(context.Background(), store.ChangeAdminPasswordParams{
		Email:    email,
		Password: string(hash),
	})
	if err != nil {
		panic(err)
	}

	return nil
}

package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"selfhosted/cmd"
	"selfhosted/router"
)

var (
	addr = flag.String("addr", ":4000", "HTTP server address")
)

func main() {
	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "password":
			if len(os.Args) < 4 {
				slog.Error("Usage: selfhosted password <email> <new_password>")
				return
			}
			email := os.Args[2]
			newPassword := os.Args[3]
			if err := cmd.ChangeAdminPassword(email, newPassword); err != nil {
				slog.Error("Failed to change admin password", "error", err)
			}
			slog.Info("Admin password changed successfully", "email", email)
		}

		return
	}

	r := router.New()

	s := http.Server{
		Addr:    *addr,
		Handler: r,
	}

	s.ListenAndServe()
}

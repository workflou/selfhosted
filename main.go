package main

import (
	"flag"
	"fmt"
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

	slog.SetLogLoggerLevel(slog.LevelInfo)

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
			slog.Debug("Admin password changed successfully", "email", email)

		case "-h", "--help", "help":
			usage()
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

func usage() {
	fmt.Println("Usage: selfhosted [command] [args]")
	fmt.Println("Commands:")
	fmt.Println("  password <email> <new_password> - Change admin password")
	fmt.Println("  (no command) - Start the web server")
	fmt.Println("Options:")
	fmt.Println("  -addr <address> - HTTP server address (default :4000)")
	os.Exit(1)
}

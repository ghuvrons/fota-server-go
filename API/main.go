package API

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ghuvrons/fota-server-go/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(MyMiddleware)
	route(r)

	defaultHost := "127.0.0.1"
	defaultPort := 9000

	if host := os.Getenv("API_HOST"); host != "" {
		defaultHost = host
	}

	if port := os.Getenv("API_PORT"); port != "" {
		nerPort, err := strconv.Atoi(port)
		if err == nil {
			defaultPort = nerPort
		}
	}
	addr := fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	fmt.Println("API server started at", addr)
	http.ListenAndServe(addr, r)
	return
}

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := models.CreateDBConnection()
		if err != nil {
			panic(err.Error())
		}
		defer func() {
			db.Close()
		}()

		ctx := context.WithValue(r.Context(), "DB_MYSQL", db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

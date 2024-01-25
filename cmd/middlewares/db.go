package middlewares

import (
	"context"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/router"
)

/*
DB injects the database connection into the request context.
*/
func DB(next router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		db, err := db.Manager.GetDB()
		if err != nil {
			router.InternalError(w, r, p)
			return err
		}

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		return next(w, r, p)
	}
}


package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/Tus1688/kim-hackathon-2023-api/authutil"
)

func EnforceAuthentication(
	requiredRoles []string, expiredIn int, passUserId bool,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			access, err := r.Cookie("access")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claim, err := authutil.ExtractClaimAccessUser(access.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			//	check if the token still valid
			if claim.IssuedAt.Time.Add(time.Duration(expiredIn) * time.Minute).Before(time.Now()) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if len(requiredRoles) == 0 {
				if passUserId {
					ctx := context.WithValue(r.Context(), "uid", claim.Uid)
					r = r.WithContext(ctx)
				}
				next.ServeHTTP(w, r)
				return
			}
			//	check if the user has the required role
			if !verifyRoles(requiredRoles, claim.Roles) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if passUserId {
				ctx := context.WithValue(r.Context(), "uid", claim.Uid)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func verifyRoles(requiredRoles []string, userRoles []string) bool {
	for _, role := range requiredRoles {
		found := false
		for _, userRole := range userRoles {
			if role == userRole {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

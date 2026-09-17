package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type contextKey string

const (
	CtxStaffID contextKey = "staff_id"
	CtxRole    contextKey = "role"
)

// CORS mengizinkan request dari frontend Next.js (beda origin/port saat development).
// Di production, ganti "*" dengan domain frontend yang sebenarnya.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", " https://pos-fnb-nu.vercel.app/")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// StaffAuth memvalidasi JWT di header Authorization untuk route khusus staff
// (admin/kasir/kitchen). Klaim staff_id & role disimpan di request context
// untuk dipakai handler selanjutnya (role-based access control).
func StaffAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				utils.Error(w, http.StatusUnauthorized, "token tidak ditemukan")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				utils.Error(w, http.StatusUnauthorized, "token tidak valid atau kedaluwarsa")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.Error(w, http.StatusUnauthorized, "klaim token tidak valid")
				return
			}

			ctx := context.WithValue(r.Context(), CtxStaffID, claims["staff_id"])
			ctx = context.WithValue(ctx, CtxRole, claims["role"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

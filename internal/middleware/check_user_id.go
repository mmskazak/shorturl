package middleware

import (
	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
	"net/http"
)

func CheckUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
		// Проверка на успешное извлечение и наличие UserID
		if !ok || payload.UserID == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

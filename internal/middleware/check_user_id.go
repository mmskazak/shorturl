package middleware

import (
	"net/http"

	"mmskazak/shorturl/internal/ctxkeys"
	"mmskazak/shorturl/internal/services/jwtbuilder"
)

// CheckUserID создает middleware для проверки наличия и корректности UserID в контексте запроса.
// Если UserID отсутствует или пустой, возвращается ошибка 401 Unauthorized.
// В противном случае запрос передается следующему обработчику.
func CheckUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем полезную нагрузку JWT из контекста запроса
		payload, ok := r.Context().Value(ctxkeys.PayLoad).(jwtbuilder.PayloadJWT)
		// Проверяем успешное извлечение payload и наличие UserID
		if !ok || payload.UserID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

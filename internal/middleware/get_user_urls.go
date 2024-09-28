package middleware

import (
	"mmskazak/shorturl/internal/services/jwttoken"
	"net/http"

	"mmskazak/shorturl/internal/config"
)

// GetUserURLsForAuth создает middleware для проверки авторизации пользователя,
// особенно для запросов к пути "/api/user/urls". Если JWT токен не валиден,
// возвращается ошибка 401 Unauthorized. В противном случае запрос передается
// следующему обработчику.
func GetUserURLsForAuth(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := cfg.SecretKey
		cookie, err := r.Cookie(authorizationCookieName)
		jwt := cookie.Value
		_, err = jwttoken.GetSignedPayloadJWT(jwt, secretKey)
		// Проверка валидности токена для специфического пути
		if err != nil && r.URL.Path == "/api/user/urls" {
			// Возвращаем ошибку 401 Unauthorized, если токен недействителен
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

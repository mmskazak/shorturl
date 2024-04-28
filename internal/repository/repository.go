package repository

var urlMap map[string]string

func InitUrlMap() {
	// Инициализация UrlMap
	urlMap = make(map[string]string) // Используйте urlMap напрямую, без скобок
}

func GetUrlMap() map[string]string {
	return urlMap // Возвращайте urlMap, а не str, которая не определена в этом контексте
}

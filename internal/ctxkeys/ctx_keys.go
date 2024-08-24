package ctxkeys

type key int

// Ключи для передачи данных в контекстрах.
const (
	// PayLoad - подписанная полезная нагрузка JWT.
	PayLoad key = iota
)

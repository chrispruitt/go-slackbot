package bot

import (
	"os"
	"strconv"
)

func getenv(key string, fallback interface{}) interface{} {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	switch fallback.(type) {
	case string:
		v := os.Getenv(key)
		if len(value) == 0 {
			return fallback
		}
		return v
	case int:
		s := os.Getenv(key)
		v, err := strconv.Atoi(s)
		if err != nil {
			return fallback
		}
		return v

	case bool:
		s := os.Getenv(key)
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fallback
		}
		return v
	default:
		return value
	}
}

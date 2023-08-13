package server

import (
	"fmt"
	"math/rand"
	"time"
)

func someBusinessLogic(duration, failChance int64) (string, error) {
	if duration < 1 || duration > 60 {
		return "", fmt.Errorf("invalid duration %d", duration)
	}

	if failChance < 0 || failChance > 100 {
		return "", fmt.Errorf("invalid fail chance %d", failChance)
	}

	// Имитация выполнения
	time.Sleep(time.Second * time.Duration(duration))

	// Результат выполнения [1, 100]
	result := rand.Int63n(100) + 1

	if failChance >= result {
		return "", fmt.Errorf("failed - result %d ", result)
	}

	return fmt.Sprintf(
		"Успешно обработано\nПродолжительность: %d секунд\nШанс успеха: %d%%",
		duration, 100-failChance,
	), nil
}

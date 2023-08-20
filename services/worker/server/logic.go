package server

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var (
	errFailed           = errors.New("failed")
	errInvalidInputData = errors.New("invalid")
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func someBusinessLogic(duration, failChance int64) (int64, string, error) {
	if duration < 1 || duration > 60 {
		return 0, "", fmt.Errorf("%w duration %d", errInvalidInputData, duration)
	}

	if failChance < 0 || failChance > 100 {
		return 0, "", fmt.Errorf("%w fail chance %d", errInvalidInputData, failChance)
	}

	totalSleep := rand.Int63n(max(duration*80/100, 1)) + 1 + duration*60/100

	runtime.Gosched()
	// Имитация выполнения
	time.Sleep(time.Second * time.Duration(totalSleep))

	// Результат выполнения [1, 100]
	result := rand.Int63n(100) + 1

	if failChance >= result {
		return result, "", fmt.Errorf("%w: result %d ", errFailed, result)
	}

	return result, fmt.Sprintf(
		"Успешно обработано\nПродолжительность: %d секунд\nШанс успеха: %d%%",
		totalSleep, 100-failChance,
	), nil
}

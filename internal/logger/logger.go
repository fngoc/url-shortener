package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log = zap.NewNop()

var loglevel = "INFO"

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize() error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(loglevel)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData: &responseData{
				status: 0,
				size:   0,
			},
		}

		start := time.Now()
		handlerFunc.ServeHTTP(&lw, r)
		duration := time.Since(start)

		Log.Info("HTTP request",
			zap.String("method", r.Method),
			zap.Int("status", lw.responseData.status),
			zap.String("path", r.URL.Path),
			zap.String("duration", duration.String()),
			zap.Int("size", lw.responseData.size),
		)
	}
}

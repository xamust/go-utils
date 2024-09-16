# Logger tutorial

### Логер подстроен под требования команды BCONNECT

Копия хранится и передается с использованием `Context'a`, разработан на основе `log/slog`  
Интерфейс взаимодействия:
```go
type Logger interface {
	Init(opts ...Option) error

	Log(ctx context.Context, level Level, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Fatal(ctx context.Context, msg string, args ...any)

	String() string
	Fields(fields map[string]interface{}) Logger
}
```

Для полноценной работы ему необходимо извлечение данных из запроса, которое выполняет метод `ExtractHeader` из пакета `go-utils/metadata/middleware`, сохранение логера в контекст. 
Все это выполняется в middlewares

Имеет несколько уровней логирования: ```[trace, debug, info, warn, error, fatal]```, которыми можно управлять опциями, к примеру:

```json
// конфиг с дебаг уровнем логирования
"logger": {
    "loglevel": "debug"
}
```
```go
package main

import (
    "github.com/xamust/go-utils/logger"
    "github.com/xamust/go-utils/metadata/middleware"
)

type Config struct {
    Logger *logger.Config `json:"logger"`
}

func main(){
    cfg := Config{} //с данными о необходимом уровне логирования
    slog := logger.NewLogger(cfg.Logger) // метод создания и инициализации логгера
	
    //мидлвар, который позволяет прокинуть копию инициализированного логгера в контекст
	// handleTestLog := logger.InjectLogger(TestLogHandler, slog) 
  
	// для того, что бы логгер заполнился данными из header запроса необходимо его обернуть его еще одним мидлваром для извлечения данных
    handleTestLog := middleware.ExtractHeader(logger.InjectLogger(TestLogHandler, slog))

    mux := http.NewServeMux()
    mux.Handle("/testlog", handleTestLog)
}

func TestLogHandler(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	log := logger.FromContextLogger(ctx)
	
	log.Info(ctx, "Test message")
}
```

Доп опции для гибкой настройки логера:
```go
    logger.NewLogger(cfg.Logger,
            WithOutput(io.Writer), 				// Перенаправить поток вывода данных
            WithContext(context.Context), 		// Добавить контекст
            WithFields(map[string]interface{}), // Добавить дополнительные необходимые поля со значениями
            WithLevel(Level), 					// сменить уровень логирования
            WithSource(), 						// включить вывод строки вызовы метода логера
            WithCallerSkipCount(int), 			// управление количества пропуска строк в стеке, для корректировки вывода
            MaxBytesMessage(int), 				// установление предела длины сообщения
    )
```

Мануальный инит и настройка логгера:
```go
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),				// Init InfoLevel
		WithOutput(buf),
		WithSource(),
	)
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Debug(context.TODO(), "Test message") //вызов не отработает, тк включен уровень Info
```

В мидлвар для фреймоврка `Echo` встроена возможность выводить стектрейс(по ключу `stacktrace`) ошибки, для этого необходимо ошибку обернуть соответствующим методом или имплементировать интерфейс
```go
package main

import "github.com/pkg/errors"

type stackTracer interface {
    StackTrace() errors.StackTrace
}

errors.WithStack(errors.New("Error message"))
```

В рамках поддержки динамически изменяемых значений по ключам `EventSource` и `EventReciever` управление ими вынесено в контекст, к примеру:
```go
// func NewContextEvent(ctx context.Context, source, receiver string) context.Context
ctx := r.Context()
log := logger.FromContextLogger(ctx)

ctx = logger.NewContextEvent(ctx, "Reciever", "Source")
log.Info(ctx, "message")
// output
// {..., "msg": "message", "EventReciever":"Source", "EventSource":"Reciever", ....}

ctx = logger.NewContextEvent(ctx, "Source", "Reciever")
log.Info(ctx, "message")
// output
// {..., "msg": "message", "EventReciever":"Reciever", "EventSource":"Source", ....}
```

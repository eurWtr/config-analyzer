# Config Analyzer

**Config Analyzer** — это быстрая и расширяемая утилита на Go для статического анализа конфигурационных файлов (JSON, YAML).

Утилита может работать как классический CLI-инструмент, так и как HTTP REST или gRPC API.

---

## Возможности

- **Поддержка форматов:** Автоматическое определение и парсинг `.json`, `.yaml`, `.yml`.
- **Поиск частых уязвимостей:** Пароли в открытом виде, отключенный TLS, привязка к `0.0.0.0`, слабые алгоритмы шифрования и оставленный `debug` режим.
- **Рекурсивное сканирование:** Рекурсивная проверка директорий одной командой.
- **Встроенные серверы:** Запуск в режиме HTTP REST или gRPC.
- **Расширяемость:** Легко написать и добавить собственные правила проверок.

---

## Быстрый старт

### Требования
* Установленный Go (версия 1.23 или выше)
* *(Опционально)* `make` для быстрой сборки

### Установка и сборка

1. Склонируйте репозиторий и перейдите в папку:
   ```bash
   git clone https://github.com/eurWtr/config-analyzer.git
   cd config-analyzer
   ```

2. Загрузите зависимости:
   ```bash
   go mod tidy
   ```
3. Установите плагины
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```
   ```bash
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. Выполните
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/analyzer.proto
   ```

5. Скомпилируйте приложение:
    * **Для Linux / macOS:**
      ```bash
      go build -o analyzer ./cmd/analyzer/
      ```
    * **Для Windows:**
      ```bash
      go build -o analyzer.exe ./cmd/analyzer/
      ```

---

## Как использовать (CLI)

*Примечание: Если вы не хотите компилировать бинарный файл, любую команду ниже можно запустить через `go run ./cmd/analyzer/main.go` вместо `./analyzer`.*
*Для вывода списка доступных команд используйте флаг ```--help```.*


### 1. Проверка одного файла
Самый простой способ проверить конкретный конфиг. Программа вернет `exit code 1`, если найдет уязвимости.
```bash
./analyzer config.yaml
```

### 2. Рекурсивная проверка папки
Ищет все файлы конфигураций внутри указанной директории.
```bash
./analyzer -r ./my-project/configs/
```

### 3. Тихий режим (Silent)
Выводит отчет, но завершает работу без ошибки (`exit code 0`).
```bash
./analyzer -s config.json
```

### 4. Чтение из потока (Stdin)
```bash
echo '{"database":{"password":"123"}}' | ./analyzer --stdin
```

---

## Использование как API-сервис

Утилиту можно запустить HTTP или gRPC сервер для проверки конфигов по сети.

### HTTP REST API
Запуск сервера:
```bash
./analyzer --http :8080
```
Тестовый запрос c JSON:
```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "config": "{\"server\":{\"tls\":false},\"db\":{\"password\":\"secret\"}}"
  }'
```
Тестовый запрос c YAML:
```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "config": "version: 1.4\nserver:\n  host: 0.0.0.0\n  port: 8080\n  tls:\n    enabled: false\n    skip_verify: true\nlog:\n  level: debug\nauth:\n  password: secret\nstorage:\n  digest-algorithm: MD5\n"
  }'
```
### gRPC API
Запуск сервера:
```bash
./analyzer --grpc :9090
```

---

## Встроенные правила проверок

Анализатор изначально ищет следующие проблемы:

| Правило | Уровень | Описание |
| :--- |:-------:| :--- |
| **Plain Password** |  HIGH   | Ищет пароли, токены, ключи API и приватные RSA-ключи в открытом виде (игнорирует ссылки на ENV и Vault). |
| **TLS Disabled** |  HIGH   | Находит явно отключенную проверку сертификатов (`insecure_skip_verify: true`) или выключенный SSL. |
| **Weak Algorithm** |  HIGH   | Ищет использование устаревших алгоритмов шифрования/хэширования (MD5, SHA1, RC4, DES). |
| **Bind Address** | MEDIUM  | Предупреждает, если сервис слушает все сетевые интерфейсы (`0.0.0.0`), открывая порты наружу. |
| **File Permissions** | MEDIUM  | *(Только при сканировании файлов)* Предупреждает, если конфиг доступен для чтения/записи всем пользователям ОС. |
| **Debug Mode** |  LOW  | Находит оставленные флаги `debug: true` или `log_level: trace`, которые могут привести к утечке данных в логи. |

---

## Расширяемость

Для добавления нового правила:
1. Создайте файл в internal/rules/
2. Реализуйте интерфейс Rule:

```go

package rules

import "config-analyzer/internal/models"

type MyNewRule struct{}

func (r *MyNewRule) Name() string {
return "my-new-rule"
}

func (r *MyNewRule) Check(config map[string]interface{}, filePath string) []models.Issue {
// Логика проверки
return nil
}
```

3. Зарегистрируйте в NewRegistry() в registry.go:

```go
    r.Register(&MyNewRule{})
```

Для добавления нового формата — расширьте ```internal/parser/parser.go```, добавив новую функцию парсинга и соответствующее определение формата.

### Работа с БД

На прокте используется PostgreSQL. Для создания новой таблицы или изменения старой требуется редактировать файлы в папке [/src/tables](/src/tables). Изменения применяются при перезапуске платформы.

Все функции для работы с сущностями должны храниться в папке [/src/core/database](/src/core/database). \
Как правило, мы хотим иметь основные функции, такие как поиск, обновление и удаление.

Для построения запросов рекомендуется использовать функцию `CreateColumnsQuery`, которая позволяет создавать SQL-строку с колонками и их значениями. \
Для работы необходимо определить структуру с тегами "db" (названия колонок в PostgreSQL).

Пример функции для изменения пользователя:

```go
type UpdateUserOptionsFields struct {
	FirstName *string  `db:"first_name"`
	LastName  *string  `db:"last_name"`
	AvatarId  **string `db:"avatar_id"`
}

type UpdateUserOptions struct {
	Id     int
	Fields *UpdateUserOptionsFields
}

func (dbClient *DatabaseClient) UpdateUser(options *UpdateUserOptions) error {
	currentTime := time.Now().UTC()

	query := fmt.Sprintf("UPDATE users SET updated_at = $1, %s WHERE id = $2", CreateColumnsQuery(options.Fields))

	_, err := dbClient.p.Exec(context.Background(), query, currentTime, options.Id)

	return err
}
```

### Работа с хендлерами

Существует две основные функции для:

- NewPublicHttpHandler (функция для регистрации публичного хендлера)
- NewPrivateHttpHandler (функция для регистрации приватного хендлера)

Приватные хендлеры имеют встроенную проверку токена и могут вернуть код 401 в случае его невалидности.

Для получения ID пользователя, от имени которого выполняется запрос, можно использовать функцию `ExtractIntParam`.

Пример:

```go
userId := core.ExtractUserId(ctx) // ctx - core.HttpCtx
```

> В случае использования функции внутри публичного хендлера требуются дополнительные проверки.

Каждый запрос состоит из трёх этапов:

1) Конвертация данных, поленных с фронта (Prepare)
2) Валидация (Validate)
3) Обработка (Handle)

```go
type httpHandler[T any] interface {
	Prepare([]byte, HttpCtx) *T
	Validate(*T, HttpCtx) error
	Handle(*T, HttpCtx) (statusCode HttpStatus, resultData []byte, err error)
}
```

> Организация запросов должна быть семантической. Например, для запроса обновления пользователей следует использовать метод PATCH, а для передачи ID следует использовать url-параметр.

Также, на проекте используется библиотека [FlatBuffers](https://flatbuffers.dev/). При получении данных с фронта и отправке ответа на фронт требуется выполнять конвертацию, используя сгенерированный код.

[Пример получения списка пользователей](/src/handlers/users/get/index.go)

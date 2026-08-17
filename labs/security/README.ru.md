# Лаборатория безопасности

## Назначение

`labs/security` содержит краткие заметки по исходным материалам, чек-листы и небольшие примеры для
Stage 7. Это не второе приложение и не общий пакет аутентификации для `book-social`.

Текущая волна — **Stage 7A**: аутентификация `book-social` v0.2.5–v0.2.6 вместе с первым блоком TDD
foundations. Stage 7B (CORS, API rate limiting и API-specific controls) не начинается, пока не
появится стабильный контракт `/api/*`.

Связанный план: [`docs/private/plan-stage7.ru.md`](../../docs/private/plan-stage7.ru.md).

## Глава 8 *Let's Go*: Stateful HTTP

Источник: Alex Edwards, *Let's Go*, 2nd edition (2025), глава 8: «Choosing a session manager»,
«Setting up the session manager» и «Working with session data».

### Краткое резюме

HTTP не сохраняет состояние между запросами. Сессия связывает несколько запросов одного пользователя
браузера: cookie передаёт случайный session token/ID, а server-side store по нему находит данные
сессии и её срок действия. Сам token не должен содержать flash message, user identity или другой
session payload.

Книга сравнивает два подхода:

| Подход                                                      | Когда уместен                                                      | Ограничение Stage 7A                                                                                              |
|-------------------------------------------------------------|--------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------|
| Client-side session data в подписанной/зашифрованной cookie | Нужны именно client-side данные и приняты их ограничения.          | Не использовать для private session data без явного решения о размере, exposure и lifecycle.                      |
| Server-side session store                                   | Нужны controlled invalidation, expiry и private server-side state. | Для auth flow предпочтителен manager, который умеет renew session ID после login и снижает риск session fixation. |

В Snippetbox книга выбирает `alexedwards/scs`: manager получает store, lifetime и middleware
`LoadAndSave`. Middleware загружает session по cookie перед handler и сохраняет изменения после
него. Это полезный ориентир, но не готовое решение для `book-social`: MySQL table и 12-hour lifetime
из книги не переносятся автоматически.

### Что применимо к `book-social`

- Выбрать server-side store и session manager после фиксации auth contract. При server-side model
  manager должен уметь renew или rotate session ID после успешной authentication.
- Хранить в cookie только непрозрачный high-entropy identifier. User ID и другие данные остаются в
  server-side session store; password, CSRF values и private content не попадают в logs или
  templates.
- Явно определить lifecycle: создание после успешного login, renewal при authentication, expiry/idle
  policy, deletion или invalidation при logout и поведение expired/missing session.
- Подключать load/save middleware только к dynamic MPA routes, которые читают или меняют session
  state. Static files и health endpoint не должны без причины создавать работу session store.
- Передавать session state handler-ам через request context, предоставленный middleware, а не через
  global variable, query parameter или доверяемый client header.
- Использовать одноразовую семантику `Pop` для flash messages; обычный `Get` — только когда значение
  должно пережить следующий request.

### Решения, требующие отдельной проработки

Глава даёт модель, но не заменяет security decisions Stage 7A:

- конкретный store `book-social` и его migration/cleanup policy;
- cookie `HttpOnly`, `Secure`, `SameSite`, name, path и development-versus-HTTPS behavior;
- absolute и idle expiration, renewal semantics и concurrent-session policy;
- login/register/logout routes, redirect и error contract, CSRF integration;
- ownership/authorization boundary для private resource;
- retention и cleanup истёкших server-side sessions, а также operational observability без записи
  session IDs.

Зафиксировать эти решения в auth contract до реализации; не выводить их только из выбранной library.

## Минимальный flow и проверяемые инварианты

```text
успешный login
  → создать/renew server-side session
  → response устанавливает opaque cookie
  → protected route загружает current user из session
  → logout invalidates session и очищает browser state
```

| Риск                                            | Минимальная проверка перед handoff                                 |
|-------------------------------------------------|--------------------------------------------------------------------|
| Anonymous user получает private resource        | HTTP test: redirect или refusal без session.                       |
| ID, зафиксированный до login, остаётся валидным | Test: successful login renews session identifier.                  |
| Logout не прекращает access                     | HTTP flow: login → protected action → logout → refusal.            |
| Flash message показывается повторно             | Test: set flash → first render shows it → refresh does not.        |
| Secret попадает к client или в logs             | Handler/log assertion: no password, raw session ID или CSRF value. |
| Static/health request создаёт session           | Router/middleware test: no session cookie или store write.         |

Использовать `httptest`/`ServeHTTP` и deterministic test store при проверке HTTP/session boundary.
Настоящий database integration test нужен только после появления конкретного persistence risk и
должен использовать явно disposable database.

## Следующий учебный шаг

До кода v0.2.5 описать session contract: actors, routes, session contents, store, cookie policy,
lifetime, renewal, logout invalidation и точные test cases. Затем изучить главу 10 *Let's Go* для
password handling, authorization и CSRF. Не добавлять login form или user model только на основании
главы 8.

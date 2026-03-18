# Link Shortener

Небольшой сервис для сокращения ссылок.

Вы вводите обычную ссылку (c https:// обязательно), сервис сохраняет её в Postgres и возвращает короткий “хэш”.
Переход по `/{hash}` делает редирект на исходный URL.

## Что умеет

- UI-страница на `/` (простая HTML-форма).
- Сокращение ссылки через `POST /short`.
- Редирект по `GET /{hash}`.
- Хранение соответствий `original_link -> hash_link` в Postgres.
- Запуск через Docker Compose (Postgres + приложение + Nginx).

## Запуск 

Проект ориентирован на запуск через Docker Compose - Вам не нужно устанавливать Go или Postgres на хосте — достаточно Docker и Docker Compose.

### 1) Подготовить переменные окружения

Скопируйте пример и при необходимости поменяйте значения:

```bash
cp .env.example .env
```

Поля `DB_...` по умолчанию настроены так, чтобы Postgres работал в отдельном контейнере (`DB_HOST=postgres`).

### 2) Запустить контейнеры

```bash
docker compose up --build
```

Или в фоне:

```bash
docker compose up -d --build
```

После старта сервисы будут доступны по умолчанию:

- Приложение (backend): http://localhost:8080
- Nginx (если настроен для вашего домена/cert): http://localhost и https://localhost

Если вы хотите работать с приложением без Nginx — используйте `http://localhost:8080`.

## Переменные окружения

Файл `.env` читается при старте приложения.

- `LOG_LEVEL` — уровень логов (сейчас в коде используется `debug` по умолчанию)
- `DB_USER` — пользователь Postgres
- `DB_PASSWORD` — пароль Postgres
- `DB_NAME` — имя базы
- `DB_HOST` — хост (в Docker Compose это обычно `postgres`)
- `DB_PORT` — порт (обычно `5432`)

## API (как пользоваться)

### Через веб‑интерфейс (UI)

Проще всего пользоваться готовой страницей в `static/web.html` — она подключается автоматически при запуске через Docker + Nginx или доступна напрямую на порту `8080`.

Шаги:

1. Откройте в браузере `http://localhost` (если используете nginx) или `http://localhost:8080` (если обращаетесь к приложению напрямую).
2. Введите полный URL в поле ввода (обязательно `http://...` или `https://...`).
3. Нажмите кнопку `Short`.
4. В поле результата появится короткая ссылка; нажмите на неё, чтобы открыть исходный URL, либо используйте кнопку копирования, чтобы вставить ссылку где нужно.

Поведение и подсказки:

- UI делает `POST /short` и показывает ответ (короткий хэш) в удобной форме.
- Если сайт не отвечает или блокирует HEAD‑запросы, UI покажет сообщение об ошибке.
- Сервис не принимает ссылки вида `http(s)://localhost:8080/...` — если вы пытаетесь сократить локальные адреса, получите ошибку.
- После получения хэша UI формирует полную ссылку по `window.location.origin + '/' + hash`.

Пример: в браузере нажмите `Short`, получите `aZ3f`, затем откройте `http://localhost:8080/aZ3f` — произойдёт редирект на исходный URL.

### Через терминал

Запрос:

```bash
curl -X POST http://localhost:8080/short \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'url=https://example.com'
```

Ответ — это строка-хэш, например `aZ3f`.
Полная короткая ссылка получается как `http://<host>/<hash>`.

### Перейти по короткой ссылке

```bash
curl -I http://localhost:8080/aZ3f
```

Вернётся `302 Found` и заголовок `Location` с исходной ссылкой.

## Про Nginx и HTTPS

В compose есть сервис `nginx`, который:

- слушает `80` и делает редирект на `443`
- на `443` проксирует запросы в `app:8080`
- отдаёт `/static/` из отдельной директории

Важно: текущий конфиг Nginx заточен под конкретный домен и сертификаты Let’s Encrypt.
Если вы запускаете проект у себя локально или на другом домене — обновите `server_name` и пути к сертификатам в `nginx/nginx.conf`, либо используйте приложение напрямую на `:8080`.

## TLS / SSL сертификаты (Let's Encrypt)

Коротко — для публичного домена удобно использовать Let's Encrypt через `certbot`. Сертификаты должны быть доступны в контейнере Nginx по тем путям, что указаны в `nginx/nginx.conf` (в этом проекте монтируем `/etc/letsencrypt` из хоста).

1) Убедитесь, что домен указывает на ваш хост 

2) На хосте установите `certbot` (или `snap install --classic certbot`).

3) Получить сертификат 

```bash
# standalone (certbot запустит временный сервер на 80)
sudo certbot certonly --standalone -d yourdomain.example.com

# после получения сертификатов они появятся в /etc/letsencrypt/live/yourdomain.example.com/
```

4) Проверить наличие файлов:

```bash
ls -l /etc/letsencrypt/live/yourdomain.example.com/
# fullchain.pem  privkey.pem  cert.pem  chain.pem
```

5) В `docker-compose.yaml` уже есть монтирование:

```yaml
nginx:
  volumes:
    - /etc/letsencrypt:/etc/letsencrypt:ro
```

Если вы получили сертификаты на хосте — они сразу станут доступны nginx-контейнеру после (пере)старта.

6) Перезапустите compose и перезагрузите nginx внутри контейнера:

```bash
docker compose up -d --build
docker compose exec nginx nginx -s reload
```

7) Проверка автообновления:

```bash
sudo certbot renew --dry-run
```

Заметки и безопасность:
- Не копируйте приватные ключи в репозиторий. Держите их в `/etc/letsencrypt` и монтируйте в контейнер.

## Структура проекта

- `main.go` — входная точка: чтение конфига, подключение к БД, инициализация таблицы, запуск HTTP-сервера
- `pkg/config` — загрузка `.env`
- `pkg/server` — роутинг и запуск `http.ListenAndServe`
- `pkg/handlers` — HTTP-обработчики (`/short`, `/{hash}`)
- `pkg/hash` — генерация base62-хэша
- `pkg/helpers` — строка подключения к БД + создание таблицы
- `static/` — UI-страница и картинки
- `docker-compose.yaml` — Postgres + app + nginx
- `Dockerfile` — сборка приложения в два этапа (builder + runtime)

## Run checklist

Короткий набор команд, чтобы быстро поднять проект (copy-paste):

```bash
# Собрать и запустить контейнеры (в фоне)
docker compose up -d --build

# Проверить статус контейнеров
docker compose ps

# Посмотреть логи (nginx / app)
docker compose logs -f nginx
docker compose logs -f app

# Проверить, что приложение отвечает
curl -I http://localhost:8080/
curl -X POST http://localhost:8080/short -d 'url=https://example.com'

# --- Остановка/перезапуск: часто полезно перед повторным `up`
# Остановить конкретные сервисы (если надо):
docker compose stop nginx postgres || true

# Корректно удалить проект (остановит и удалит контейнеры, сети, удалит "брошенные" контейнеры):
docker compose down --remove-orphans

# Если остаются контейнеры с фиксированными именами (например `postgres`), можно удалить их явно:
docker rm -f postgres || true

# Опционально (ДЕСТРУКТИВНО): остановить и удалить ВСЕ контейнеры на хосте
# Используйте с осторожностью — удалит все локально запущенные контейнеры
docker stop $(docker ps -q) 2>/dev/null || true
docker rm -f $(docker ps -aq) 2>/dev/null || true
```

Если docker пишет про конфликт имени контейнера `postgres`, удалите или переименуйте старый контейнер:

```bash
docker ps -a --filter "name=postgres"
docker rm -f postgres
# или
docker rename postgres postgres_old
```

## API examples

Примеры запросов и ожидаемого поведения.

- Shorten (request):

```bash
curl -X POST http://localhost:8080/short \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'url=https://example.com'
```

- Shorten (response): просто тело с хэшем, например:

```
1aZf
```

- Follow short link (GET):

```bash
curl -I http://localhost:8080/1aZf
```

- Expected response headers (on redirect):

```
HTTP/1.1 302 Found
Location: https://example.com
```

## Troubleshooting

### Конфликт имени контейнера `postgres`

Если Docker пишет, что контейнер с именем `postgres` уже существует, посмотрите, кто его создал:

```bash
docker ps -a --filter "name=postgres"
docker inspect --format='Project={{index .Config.Labels "com.docker.compose.project"}} Service={{index .Config.Labels "com.docker.compose.service"}}' postgres
```

И при необходимости удалите/переименуйте:

```bash
docker rm -f postgres
# или
docker rename postgres postgres_old
```

---

Подробнее про устройство проекта — в `ARCHITECTURE.md`.

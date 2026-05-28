        # redis — RDB-снэпшоты

        Homework-шаблон для урока **l1_rdb_snapshots** (RDB-снэпшоты) на платформе Vibe Learn.

        ## Что делать

        Напиши Go-программу, которая измеряет memory bloat и latency при BGSAVE.

**Окружение:** docker-compose с одним Redis (RDB включён, AOF выключен).

**Шаги:**
1. Загрузи ~1M ключей через pipeline (значения по 100-500 байт — реалистичный размер).
2. Запусти фоновую write-heavy нагрузку: ~1000 `SET` существующих ключей в секунду
   (именно перезапись — она триггерит COW на уже занятых страницах).
3. Триггерни `BGSAVE` и параллельно семпли `INFO persistence` / `INFO memory`:
   `used_memory_rss`, `rdb_last_cow_size`, `latest_fork_usec`, `rdb_last_bgsave_status`.
4. Построй отчёт: пиковый прирост RSS относительно baseline (в %), длительность fork,
   длительность всего bgsave.

**CI-проверки в template repo:**
- `assert rss_spike_pct > 50` — под write-heavy COW зафиксирован рост RSS.
- `assert rdb_last_bgsave_status == "ok"` — дамп прошёл успешно.
- `assert latest_fork_usec > 0` — fork замерен.

## Контекст (из transfer-задачи урока)

Тебе достался Redis для primary-store: пользовательские сессии (TTL=24ч) и счётчики
кликов. SLA: «потеря не более 60 секунд». Текущий конфиг: `save 3600 1, save 300 100,
save 60 10000`. Что плохо в этом конфиге, как починить, и стоит ли включать AOF?

## Recap из урока

- RDB — periodic binary snapshot. BGSAVE = fork() + child пишет dump.rdb. Атомарно через temp+rename.
- Триггер: `save N M` в конфиге (N секунд + M изменений). Сработало любое — запускается BGSAVE.
- Главный риск — **удвоение памяти** при COW на write-heavy нагрузке. Следи за `fork_failed_count`.
- Окно потерь = интервал между save. Худший случай: kill -9 за 1мс до планового дампа = весь интервал данных потерян.
- RDB достаточно для кэша. Для primary-store с SLA «секунды потерь» — AOF или RDB+AOF гибрид.

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go`, прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - Docker + docker-compose — `docker compose up -d` поднимает single-node Redis 7 на `localhost:6379` (с включёнными keyspace-notifications и AOF). Адрес переопределяется через env `REDIS_ADDR`.

        ## Запуск

        ```bash
        # Поднять локальный Redis
        docker compose up -d

        # Прогнать тесты (интеграционный включается через REDIS_INTEGRATION=1)
        go test ./...
        REDIS_INTEGRATION=1 go test ./...

        # Запустить main (печатает marker; замени stub на реализацию)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.

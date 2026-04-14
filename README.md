# Task Service

РЎРµСЂРІРёСЃ РґР»СЏ СѓРїСЂР°РІР»РµРЅРёСЏ Р·Р°РґР°С‡Р°РјРё СЃ HTTP API РЅР° Go.

## РўСЂРµР±РѕРІР°РЅРёСЏ

- Go `1.23+`
- Docker Рё Docker Compose

## Р‘С‹СЃС‚СЂС‹Р№ Р·Р°РїСѓСЃРє С‡РµСЂРµР· Docker Compose

```bash
docker compose up --build
```

РџРѕСЃР»Рµ Р·Р°РїСѓСЃРєР° СЃРµСЂРІРёСЃ Р±СѓРґРµС‚ РґРѕСЃС‚СѓРїРµРЅ РїРѕ Р°РґСЂРµСЃСѓ `http://localhost:9000`.

Р•СЃР»Рё `postgres` СѓР¶Рµ Р·Р°РїСѓСЃРєР°Р»СЃСЏ СЂР°РЅРµРµ СЃРѕ СЃС‚Р°СЂРѕР№ СЃС…РµРјРѕР№, РїРµСЂРµСЃРѕР·РґР°Р№ volume:

```bash
docker compose down -v
docker compose up --build
```

РџСЂРёС‡РёРЅР° РІ С‚РѕРј, С‡С‚Рѕ SQL-С„Р°Р№Р» РёР· `migrations/0001_create_tasks.up.sql` РјРѕРЅС‚РёСЂСѓРµС‚СЃСЏ РІ `docker-entrypoint-initdb.d` Рё РїСЂРёРјРµРЅСЏРµС‚СЃСЏ С‚РѕР»СЊРєРѕ РїСЂРё РёРЅРёС†РёР°Р»РёР·Р°С†РёРё РїСѓСЃС‚РѕРіРѕ data volume.

## Swagger

Swagger UI:

```text
http://localhost:9000/swagger/
```

OpenAPI JSON:

```text
http://localhost:9000/swagger/openapi.json
```

## API

Р‘Р°Р·РѕРІС‹Р№ РїСЂРµС„РёРєСЃ API:

```text
/api/v1
```

РћСЃРЅРѕРІРЅС‹Рµ РјР°СЂС€СЂСѓС‚С‹:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`

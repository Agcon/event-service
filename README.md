# Инструкция по загрузке данных в MongoDB с использованием Docker

Этот документ содержит пошаговую инструкцию по развертыванию MongoDB в контейнере Docker и загрузке данных из файла `events.json` в базу данных MongoDB.

## Необходимые инструменты

Для выполнения всех шагов вам понадобятся:

- **Docker** — для запуска контейнера с MongoDB.
- **MongoDB** — система управления базами данных.
- **mongoimport** — инструмент для импорта данных в MongoDB.
- **mongosh** — оболочка для работы с MongoDB через командную строку.

## Шаги

Сначала необходимо скачать и запустить контейнер с MongoDB.

После установки Docker выполните следующую команду для запуска контейнера с MongoDB:

```bash
docker run -d --name mongodb -p 27017:27017 mongo:latest
```

Теперь вам нужно передать файл с данными (например, events.json) в контейнер. Для этого используем команду docker cp, которая копирует файлы с вашей машины в контейнер Docker.

```bash
docker cp events.json <container_id>:/data/events.json
```

Теперь вам нужно подключиться к контейнеру и выполнить команду импорта данных с помощью mongoimport.

```bash
docker exec -it <container_id> bash
```

Далее используйте команду mongoimport для загрузки данных из файла в MongoDB.

```bash
mongoimport --db events_db --collection events --file /data/events.json --jsonArray --drop
```

После успешного выполнения команды данные из файла будут загружены в коллекцию events базы данных events_db.

## Сервис-роутер для обмена данными с Minio Object Storage.


### Содержание
- [Сборка и тестирование](https://github.com/blablatov/s3minio-router#сборка-локально)
- [Схема сервиса](https://github.com/blablatov/s3minio-router#блок-схема-сервиса)
- [Создание бакета](https://github.com/blablatov/s3minio-router#создание-бакета)


### Описание
* Endpoint сервиса.
Точка доступа для клиентов к сервису реализована через `http`-запрос вида:
```sh
http://localhost:8080
```

* Загрузчик (`downloader`) для загрузки медиа файлов из бакета.
Скачиваемый медиа-контент, стримится в браузере, если последний поддерживает отображение видеоформата.
Форматы - `.gif`, `.mp4` загружаются с бакета `minio` в локальную папку `download` и оттуда стримятся в браузере. Форматы - `.avi`, `.mkv` только загружаются в эту папку. Для их отображение необходимо использовать сторонние кодеки.

Точка доступа к этому загрузчику GET, Headers (uuid: имя бакета):
```sh
http://localhost:8080/stream/<filename.extension>
```
* Загрузчик (`uploader`) для выгрузки медиа файлов в бакет.
Выгружает выбранный файл из локальной папки `upload` и сохраняет с нужным названием в бакет `minio`.
Точка доступа к этому загрузчику POST, Headers (uuid: ид юзера):
```sh
http://localhost:8080/upstream/<filename.extension>
```

### Сборка локально
[Содержание](https://github.com/blablatov/s3minio-router#содержание)

```sh
go build -race
```

### Тестирование локально
```sh
go test -v .
```

### Использование
```sh
./s3minio-router
```


### Блок-схема сервиса
[Содержание](https://github.com/blablatov/s3minio-router#содержание)


```mermaid
graph TB

  SubGraph1Flow
  subgraph "Minio remote service"
  SubGraph1Flow(Object Storage)
  end

  subgraph "HTTP-server router"
  SubGraph2Flow(HTTP-Module `s3minio-router`) --> S3-Module`downloader` -- S3/REST_HTTP/1.1 <--> SubGraph1Flow
  SubGraph2Flow(HTTP-Module `s3minio-router`) --> S3-Module`uploader` -- S3/REST_HTTP/1.1 <--> SubGraph1Flow
  end

  subgraph "HTTP-clients"
  Node1[Requests from clients] -- REST_HTTP/1.1 <--> SubGraph2Flow
end
```

---

### Создание бакета
[Содержание](https://github.com/blablatov/s3minio-router#содержание)

* Получение `uuid` пользователя из заголовка `HTTP`-запроса.
Модулем `maker-bucket` выполняется проверка, если в `minio` есть бакет у которого название соответствует значению этого `uuid`, выполняется обмен данными с этим бакетом.
Если соответствующий бакет не найден, создается новый бакет, которому присваивается имя вновь полученного `uuid`.
Каждый пользователь является владельцем своего бакета с именем как у своего `uuid`.

### Блок-схема
[Содержание](https://github.com/blablatov/s3minio-router#содержание)

```mermaid
graph TB

  SubGraph1Flow
  subgraph "Minio remote service"
  SubGraph1Flow(Object Storage)
  end

  subgraph "HTTP-server router"
  SubGraph2Flow(HTTP-Module `s3minio-router`) --> S3-Module`maker-bucket` --> S3-Module`downloader` -- S3/REST_HTTP/1.1 <--> SubGraph1Flow
  SubGraph2Flow(HTTP-Module `s3minio-router`) --> S3-Module`maker-bucket` --> S3-Module`uploader` -- S3/REST_HTTP/1.1 <--> SubGraph1Flow
  SubGraph2Flow(HTTP-Module `s3minio-router`) --> S3-Module`maker-bucket` -- S3/REST_HTTP/1.1 <--> SubGraph1Flow
  end

  subgraph "HTTP-clients"
  Node1[Requests from clients] -- REST_HTTP/1.1 <--> SubGraph2Flow
end
```

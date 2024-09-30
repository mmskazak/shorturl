# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `http://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template http://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](http://github.com/Yandex-Practicum/go-autotests).

## Оценка покрытия кода

##### Для вычисления покрытия кода используйте следующие команды:

##### Запустить тесты и сгенерировать отчет покрытия
go test -v -coverpkg=./... -coverprofile=profile.cov ./...

##### Проанализировать отчет покрытия
go tool cover -func profile.cov

##### Посмотреть покрытие в веб-браузере
go tool cover -html=profile.cov

##### Покрытие с исключением сгенерированного кода
go tool cover -func cover.out
cat cover.out.tmp | grep -v "pb.go" > cover.out
cat cover.out | grep -v "pb.go" > cover.out
go tool cover -func cover.out



## Запуск кастомного линтера
Переходим в папку с линтером.
``````
cd .\cmd\staticlint\
``````
Запускаем билд линтера.
``````
go build -o mylint.exe
``````
Копируем сбилженый линт в корень проекта.
Запускаем линтер.
``````
.\mylint.exe -config staticlint.json .\internal\...
.\mylint.exe -config staticlint.json .\...
``````
Файл с конфигурацией есть в папку с линтером, а так же его копия в корне проекта.

## Генерация protoc
``````
 protoc --proto_path=internal/proto --go_out=internal/proto --go_opt=paths=source_relative --go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative shorturl.proto
``````
``````
 protoc --proto_path internal/proto --go_out internal/proto --go_opt paths=source_relative --go-grpc_out internal/proto --go-grpc_opt paths=source_relative shorturl.proto
``````


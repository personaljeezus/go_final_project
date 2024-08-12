# Что это?
WEB-server для todo листа, написан на GIN web framework, бд sqlite
# Сколько выполнено дополнительных заданий
Из заданий со звездочкой только .env файл для порта и адреса бд, * *реализован только базовый api* *
# Инстукция запуска кода локально
***Запуск командой go run .***
localhost:7540 рабочий порт
# Запуск тестов
***go test ./tests***
## Параметры из settings.go
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = ``
# Docker and etc.
Я не знаю как убрать строчку cgo под .env, поэтому она всё ещё здесь
## Docker build
***docker build -tag go-final:v1 .***
## Docker run
***docker run -d -p 7540:7540 --name go_final go-final:v1***
## Docker stop
***docker stop go_final***
## Docker container remove
***docker rm -f go_final***
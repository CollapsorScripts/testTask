FROM --platform=linux/amd64 amd64/golang:latest AS builder

SHELL ["/bin/bash", "-c"]

# Устанавливаем переменные
ENV GOARCH=amd64
ENV TARGETOS=linux
ENV ENTRY_POINT=./cmd/entrypoint
ENV PROGRAM=./service
ENV WKDIR=/build

# Рабочая директория
WORKDIR ${WKDIR}
COPY . ${WKDIR}

# Компилируем
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${GOARCH} go build -o ${PROGRAM} ${ENTRY_POINT}

# Создаем финальный образ
FROM amd64/alpine:latest

# Устанавливаем переменные
ENV PROGRAM=service
ENV WKDIR=/app
ENV BUILDIR=/build
#Порт для прослушки
ENV PORT=8080

# Рабочая директория
WORKDIR ${WKDIR}

RUN apk add nano

# Копируем исполняемый файл из предыдущего образа
COPY --from=builder ${BUILDIR}/${PROGRAM} ./${PROGRAM}

# Добавляем сертификаты
RUN apk add --upgrade --no-cache ca-certificates && update-ca-certificates

# Устанавливаем время
RUN apk add tzdata && echo "Europe/Moscow" > /etc/timezone && ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
#

# Копируем файл конфигурации в контейнер
COPY config/prod.yaml .
RUN update-ca-certificates

# Открываем порты
EXPOSE ${PORT}
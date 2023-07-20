FROM alpine:latest

ENV LANG="C.UTF-8" \
    TZ=Asia/Shanghai \
    PUID=1000 \
    PGID=1000 \
    ANIMEGO_CONFIG="/data/animego.yaml" \
    ANIMEGO_DATA_PATH="/data" \
    ANIMEGO_DOWNLOAD_PATH="/download" \
    ANIMEGO_SAVE_PATH="/anime" \
    ANIMEGO_CONFIG_BACKUP=1

COPY AnimeGo /app/AnimeGo

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/localtime \
    && apk del tzdata

WORKDIR /app

ENTRYPOINT ["sh", "-c", "/app/AnimeGo"]

EXPOSE 7991
VOLUME ["/data", "/download", "/anime"]

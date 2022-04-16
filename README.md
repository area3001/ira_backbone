# Go-IRA
Ira server

## Generating API docs
```shell
~/go/bin/swag init -g main.go --output docs
```

## Building docker image
```shell
docker buildx build --platform linux/amd64,linux/arm64 --push -t calmera/goira .
```
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

## Run downloaded binary on MacOS

This app is not notarized, so you need to manually open it as follows:

- Make binary executable with `chmod +x goira-LATESTVERSION+darwin-YOURARCH`.
- Right click on the binary in Finder and select 'Open With' then 'Other'.
- In the "Enable" drop-down menu select "All Applications".
- Go to 'Utilities' and select "Terminal".

You can now run the binary in your terminal.

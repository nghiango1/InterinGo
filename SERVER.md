# Server mode

It better to have the server public InterinGo interpreter so that others can easier access, evaluating the language

## Make tools

I setup Makefile to handle CLI operation, use `make build-run` to rebuild and start the server
- `make` or `make all` or `make help`: Show all option command
- `make build`: Build/Rebuild `templ` template files and generating tailwindcss stylesheet
- `make run`: Run built the server
- `make build-run`: Do both

Example
```sh
make
```

## Dev mode

Golang doesn't have watch mode, but `templ` and `tailwindcss` have it
- `make tailwind-watch`: Tailwind watch mode - Auto rebuild when files change
- `make templ-watch`: Templ watch mode - Auto rebuild when files change
- `go run . -s` or `make go-run`: Run the server without build

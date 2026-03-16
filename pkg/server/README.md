# `server` - Docs

This folder contain all the server file. Which kept it as a seperated component of InterinGo project.

## Structure

### Source file

- `index.templ`: All website component and template
- `server.go`: Server page handler and listener

### Built file

> This mean we don't have to install the build tool on other machine when we clone the repo and just need to test REPL

- `index_templ.go`: generated from `templ generate` command
- `assets/stylesheet.css`: generated from `tailwind` CLI

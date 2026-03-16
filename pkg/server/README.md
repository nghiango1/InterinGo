# `server` - Docs

This folder contain all the server file. Which kept it as a seperated component of InterinGo project.

## Structure

### Source file

- `index.templ`: All website component and template
- `server.go`: Server page handler and listener

### Built file

> This mean we don't have to install the build tool on other machine when we clone the repo and just need to test REPL

- `index_templ.go`: generated from `templ generate` command


## Notes

embed useage note: `content/**/*` to cover all file

I even try to do thing like this as I belived `content` is actually good enough

- Tested a lot with fs.ReadDir() to understand that the file just not there
- Gin community also help with middleware to handle file that not just stole
the whole "/" path (unlike gin.StaticFS which I can use for `/assets`)

Here is what I tested
```go
func traversal(r *gin.Engine, dist string, curr string) {
	webpage, err := static.EmbedFolder(embedContent, dist)
	log.Printf("[INFO] Server static FS `%v` at `%v`\n", curr, dist)
	r.Use(static.Serve(curr, webpage))

	list, err := embedContent.ReadDir(dist)
	if err != nil {
		log.Printf("[ERROR] err: %v\n", err.Error())
		return
	}
	log.Printf("[huh] %v\n", list)
	for _, d := range list {
		if d.IsDir() {
			traversal(r, filepath.Join(dist, d.Name()), path.Join(curr, d.Name()))
		}
	}
}
```

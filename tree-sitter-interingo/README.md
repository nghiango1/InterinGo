# tree-sitter-interingo

Providing highlighting and language support for Iteringo by using treesitter

## Using tree-sitter cli

Prerequire install
```sh
apt install make
apt install node
```

Install tree-sitter cli globally
```sh
npm install -g tree-sitter-cli
```

Config tree-sitter, adding language file to place
```sh
tree-sitter init-config
mkdir ~/dev
ln -s ~/path/to/interpreter/tree-sitter-interingo ~/dev/
```

Try parsing, you can test with any other file in `../test/`
```sh
# "parse": "tree-sitter parse ../test/function-02.iig",
npm run parse
```

Try hightlighting, you can also test with any other file in `../test/`
```sh
# "highlight": "tree-sitter highlight ../test/function-02.iig",
npm run highlight
```

## Neovim config for nvim-treesitter

### Adding parsers

Interingo and it associate file `iig` is not a default supported languages, so to get syntax hightlighting you will have to config both `nvim-treesitter` and vim to recognize it.

```lua
local parser_config = require "nvim-treesitter.parsers".get_parser_configs()
parser_config.interingo = {
  install_info = {
    url = "~/path/to/tree-sitter-interingo", -- local path of cloned repo
    files = {"src/parser.c"},
    branch = "main", -- default branch in case of git repo if different from master
    generate_requires_npm = false, -- if stand-alone parser without npm dependencies
    requires_generate_from_grammar = false, -- if folder contains pre-generated src/parser.c
  },
}

vim.treesitter.language.register('interingo', 'interingo')
vim.filetype.add({
  extension = {
    iig = 'interingo',
  },
})
```

After that, start `nvim` and run command `:TSInstall interingo`.

### Adding queries

`:TSInstall` will not copy query files from the grammar repository. To make the installed InterinGo grammar to be useful, we must manually add query files to nvim-treesitter installation directory.

```sh
ln -s ~/path/to/interpreter/tree-sitter-interingo/queries ~/.local/share/nvim/lazy/nvim-treesitter/interingo
#      ^^^^^^^^^^^^^^^^^^^^ change to repo directory                          ^^^^^^^^^^^^^^^^^^^^ repace this path if you not using lazy package manager
```

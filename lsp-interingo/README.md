# LSP

LSP for interingo language

## Build and install

### Get LSP server executable 

Build the code with

```sh
go build .
```

Add build file to default local bin path 

```sh
ln -s /path/to/build/file/interingo-lsp ~/.local/bin
#     ^^^^^^^^^^^^^^^^^^^ change to your machine path
export PATH="$PATH:~/.local/bin"
```

### Setup LSP for Neovim with `iig` file

Add this to Neovim `after/plugin/iig.lua` file

```lua
vim.filetype.add({ extension = { iig = 'interingo', }, })

local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

local function custom_root_dir(filename, bufnr)
    return vim.fn.getcwd()
end

configs.interingo = {
  default_config = {
    cmd = { 'interingo-lsp' },
    filetypes = { 'interingo' },
    -- root_dir = lspconfig.util.root_pattern('.git', 'deluge'),
    root_dir = custom_root_dir,
    settings = {}
  }
}

lspconfig.interingo.setup({
  on_attach = on_attach,
  capabilities = capabilities,
  filetypes = { "interingo" },
})
```

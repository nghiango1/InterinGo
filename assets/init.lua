vim.opt.number = true

vim.diagnostic.config({
	virtual_text = true
})

vim.api.nvim_create_autocmd('LspAttach', {
	group = vim.api.nvim_create_augroup('UserLspConfig', {}),
	callback = function(ev)
		-- Enable completion triggered by <c-x><c-o>
		vim.bo[ev.buf].omnifunc = 'v:lua.vim.lsp.omnifunc'

		-- Buffer local mappings.
		-- See `:help vim.lsp.*` for documentation on any of the below functions
		local opts = { buffer = ev.buf }
		vim.keymap.set('n', 'gD', vim.lsp.buf.declaration, opts)
		vim.keymap.set('n', 'gd', vim.lsp.buf.definition, opts)
		vim.keymap.set('n', 'K', vim.lsp.buf.hover, opts)
		vim.keymap.set('n', 'gi', vim.lsp.buf.implementation, opts)
		vim.keymap.set('n', '<C-k>', vim.lsp.buf.signature_help, opts)
		vim.keymap.set('n', '<space>wa', vim.lsp.buf.add_workspace_folder, opts)
		vim.keymap.set('n', '<space>wr', vim.lsp.buf.remove_workspace_folder, opts)
		vim.keymap.set('n', '<space>wl', function()
			print(vim.inspect(vim.lsp.buf.list_workspace_folders()))
		end, opts)
		vim.keymap.set('n', '<space>D', vim.lsp.buf.type_definition, opts)
		vim.keymap.set('n', '<space>rn', vim.lsp.buf.rename, opts)
		vim.keymap.set({ 'n', 'v' }, '<space>ca', vim.lsp.buf.code_action, opts)
		vim.keymap.set('n', 'gr', vim.lsp.buf.references, opts)
		vim.keymap.set('n', '<space>f', function()
			vim.lsp.buf.format { async = true }
		end, opts)
		vim.keymap.set('n', '<space>a', function()
			vim.lsp.inlay_hint.enable(not vim.lsp.inlay_hint.is_enabled())
		end, opts)
	end,
})

-- local lspconfig = require('lspconfig')
-- local configs = require('lspconfig.configs')

local function custom_root_dir(filename, bufnr)
	return vim.fn.getcwd()
end

-- REF: https://neovim.io/doc/user/lsp/#vim.lsp.Config
vim.filetype.add({
	extension = { iig = 'interingo', },
})
-- Some how doesn't auto attach to the buffer
vim.lsp.config['interingo_ls'] = {
	cmd = { 'interingo-lsp' },
	filetypes = { 'interingo' },
	-- root_dir = lspconfig.util.root_pattern('.git', 'deluge'),
	root_dir = custom_root_dir,
	settings = {}
}
vim.lsp.enable('interingo_ls')

-- This help manually 
vim.api.nvim_create_autocmd('FileType', {
	pattern = 'interingo',
	callback = function()
		vim.lsp.start(vim.lsp.config['interingo_ls'])
	end,
})

-- Some coloring, should be overided by colorscheme
vim.api.nvim_set_hl(0, '@lsp.type.keyword',     { fg = '#569CD6' })
vim.api.nvim_set_hl(0, '@lsp.type.type',        { fg = '#4EC9B0' })
vim.api.nvim_set_hl(0, '@lsp.type.class',       { fg = '#4EC9B0' })
vim.api.nvim_set_hl(0, '@lsp.type.interface',   { fg = '#4EC9B0' })
vim.api.nvim_set_hl(0, '@lsp.type.enum',        { fg = '#4EC9B0' })
vim.api.nvim_set_hl(0, '@lsp.type.function',    { fg = '#DCDCAA' })
vim.api.nvim_set_hl(0, '@lsp.type.method',      { fg = '#DCDCAA' })
vim.api.nvim_set_hl(0, '@lsp.type.variable',    { fg = '#9CDCFE' })
vim.api.nvim_set_hl(0, '@lsp.type.property',    { fg = '#9CDCFE' })
vim.api.nvim_set_hl(0, '@lsp.type.parameter',   { fg = '#9CDCFE' })
vim.api.nvim_set_hl(0, '@lsp.type.string',      { fg = '#CE9178' })
vim.api.nvim_set_hl(0, '@lsp.type.number',      { fg = '#B5CEA8' })
vim.api.nvim_set_hl(0, '@lsp.type.enumMember',  { fg = '#4FC1FF' })
vim.api.nvim_set_hl(0, '@lsp.type.comment',     { fg = '#6A9955' })
vim.api.nvim_set_hl(0, '@lsp.type.macro',       { fg = '#C586C0' })
vim.api.nvim_set_hl(0, '@lsp.type.decorator',   { fg = '#C586C0' })
vim.api.nvim_set_hl(0, '@lsp.type.namespace',   { fg = '#4EC9B0' })
vim.api.nvim_set_hl(0, '@lsp.type.operator',    { fg = '#D4D4D4' })

vim.cmd[[colorscheme retrobox]]

require('vim.lsp.log').set_format_func(vim.inspect)

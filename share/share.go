package share

// Hack to share flag value as a global variable
var (
	// Verbose cause the REPL to print out debug value for Lexer and Parser
	VerboseMode bool

	// Hotload flag cause page to be populated via os.ReadFile() instead of cached
	// at Init time
	HotLoad bool
)

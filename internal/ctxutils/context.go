package ctxutils

var (
	// ContextDbKey is the key used to store the database connection in the context.Context
	ContextDbKey = "db"
	// ContextConfigKey is the key used to store the configuration in the context.Context
	ContextConfigKey = "cfg"
	// ContextTraktKey is the key used to store the Trakt client in the context.Context
	ContextPlexKey = "plex"
)

package resources

import "embed"

//go:embed dist app.wasm profile-default.png
var Resources embed.FS

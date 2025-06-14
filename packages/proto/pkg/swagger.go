package assets

import _ "embed"

//go:embed proto/docs/swagger/api.swagger.json
var SwaggerJSON []byte

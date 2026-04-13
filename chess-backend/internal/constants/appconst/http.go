package appconst

const (
	HeaderContentType = "Content-Type"
	MimeTypeJSON      = "application/json"
)

const (
	MsgInvalidRequest = "invalid request"
	MsgUnauthorized   = "unauthorized"
	MsgInvalidToken   = "invalid token"
)

const (
	HeaderAuthorization = "Authorization"
	AuthSchemeBearer    = "Bearer"
	QueryParamToken     = "token"

	MsgMissingAuthHeader = "missing authorization header"
	MsgInvalidAuthHeader = "invalid authorization header"
	MsgMissingToken      = "missing token"
)

const (
	WSClientBufferSize = 256
)

const (
	MsgWebSocketUpgradeFailed = "websocket upgrade failed"
	UsernameUnknown           = "Unknown"
)

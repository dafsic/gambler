package gamblerhandler

import "github.com/dafsic/gambler/version"

type errcodeT struct {
	//200-300:success;400-500:client error;500-600:server error
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Ver  string `json:"config"`
}

var ver = version.TransferVersion.String()

func responseSuccess(msg string) *errcodeT {
	return &errcodeT{Code: 200, Msg: msg, Ver: ver}
}

var (
	ErrNoData      = errcodeT{201, "Ok, but no data", ver}
	ErrNoCompleted = errcodeT{202, "Not completed yet", ver}

	ErrAuth            = errcodeT{401, "Auth failed", ver}
	ErrPermission      = errcodeT{403, "Permission Denied", ver}
	ErrNotFound        = errcodeT{404, "Not found", ver}
	ErrSign            = errcodeT{413, "Sign not match", ver}
	ErrLimit           = errcodeT{414, "Over the limit", ver}
	ErrTimestamp       = errcodeT{415, "Invalid request timestamp", ver}
	ErrSwitch          = errcodeT{426, "Switch is closed", ver}
	ErrIncorrectFormat = errcodeT{421, "Post body format incorrect", ver}
	ErrClientIP        = errcodeT{424, "Invalid client ip", ver}
	ErrValue           = errcodeT{425, "Incorrect type or value", ver}
	ErrRange           = errcodeT{426, "Invalid time range", ver}

	ErrInternalError = errcodeT{500, "Internal server error", ver}
)

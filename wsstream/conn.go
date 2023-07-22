package wsstream

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/websocket"
)

var (
	// connectionUpgradeRegex matches any Connection header value that includes upgrade
	connectionUpgradeRegex = regexp.MustCompile(`(^|.*,\\s*)upgrade($|\\s*,)`)
)

// IsWebSocketRequest returns true if the incoming request contains connection upgrade headers
// for WebSockets.
func IsWebSocketRequest(req *http.Request) bool {
	if !strings.EqualFold(req.Header.Get("Upgrade"), "websocket") {
		return false
	}
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection")))
}

// handshake ensures the provided user protocol matches one of the allowed protocols. It returns
// no error if no protocol is specified.
func handshake(config *websocket.Config, req *http.Request, allowed []string) error {
	protocols := config.Protocol
	if len(protocols) == 0 {
		protocols = []string{""}
	}

	for _, protocol := range protocols {
		for _, allow := range allowed {
			if allow == protocol {
				config.Protocol = []string{protocol}
				return nil
			}
		}
	}

	return fmt.Errorf("requested protocol(s) are not supported: %v; supports %v", config.Protocol, allowed)
}

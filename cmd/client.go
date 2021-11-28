package cmd

import "net/http"

type Client struct {
	// scheme sets the scheme for the client
	scheme string
	// host holds the server address to connect to
	host string
	// proto holds the client protocol i.e. unix.
	proto string
	// addr holds the client address.
	addr string
	// basePath holds the path to prepend to the requests.
	basePath string
	// client used to send and receive http requests.
	client *http.Client
	// version of the server to talk to.
	version string
	// custom http headers configured by users.
	customHTTPHeaders map[string]string
	// manualOverride is set to true when the version was set by users.
	manualOverride bool

	// negotiateVersion indicates if the client should automatically negotiate
	// the API version to use when making requests. API version negotiation is
	// performed on the first request, after which negotiated is set to "true"
	// so that subsequent requests do not re-negotiate.
	negotiateVersion bool

	// negotiated indicates that API version negotiation took place
	negotiated bool
}

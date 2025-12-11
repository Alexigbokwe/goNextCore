package core

type HttpStatusType struct{}

const (
	HttpStatusContinue           = 100 // RFC 9110, 15.2.1
	HttpStatusSwitchingProtocols = 101 // RFC 9110, 15.2.2
	HttpStatusProcessing         = 102 // RFC 2518, 10.1
	HttpStatusEarlyHints         = 103 // RFC 8297

	HttpStatusOK                          = 200 // RFC 9110, 15.3.1
	HttpStatusCreated                     = 201 // RFC 9110, 15.3.2
	HttpStatusAccepted                    = 202 // RFC 9110, 15.3.3
	HttpStatusNonAuthoritativeInformation = 203 // RFC 9110, 15.3.4
	HttpStatusNoContent                   = 204 // RFC 9110, 15.3.5
	HttpStatusResetContent                = 205 // RFC 9110, 15.3.6
	HttpStatusPartialContent              = 206 // RFC 9110, 15.3.7
	HttpStatusMultiStatus                 = 207 // RFC 4918, 11.1
	HttpStatusAlreadyReported             = 208 // RFC 5842, 7.1
	HttpStatusIMUsed                      = 226 // RFC 3229, 10.4.1

	HttpStatusMultipleChoices   = 300 // RFC 9110, 15.4.1
	HttpStatusMovedPermanently  = 301 // RFC 9110, 15.4.2
	HttpStatusFound             = 302 // RFC 9110, 15.4.3
	HttpStatusSeeOther          = 303 // RFC 9110, 15.4.4
	HttpStatusNotModified       = 304 // RFC 9110, 15.4.5
	HttpStatusUseProxy          = 305 // RFC 9110, 15.4.6
	HttpStatusSwitchProxy       = 306 // RFC 9110, 15.4.7 (Unused)
	HttpStatusTemporaryRedirect = 307 // RFC 9110, 15.4.8
	HttpStatusPermanentRedirect = 308 // RFC 9110, 15.4.9

	HttpStatusBadRequest                   = 400 // RFC 9110, 15.5.1
	HttpStatusUnauthorized                 = 401 // RFC 9110, 15.5.2
	HttpStatusPaymentRequired              = 402 // RFC 9110, 15.5.3
	HttpStatusForbidden                    = 403 // RFC 9110, 15.5.4
	HttpStatusNotFound                     = 404 // RFC 9110, 15.5.5
	HttpStatusMethodNotAllowed             = 405 // RFC 9110, 15.5.6
	HttpStatusNotAcceptable                = 406 // RFC 9110, 15.5.7
	HttpStatusProxyAuthRequired            = 407 // RFC 9110, 15.5.8
	HttpStatusRequestTimeout               = 408 // RFC 9110, 15.5.9
	HttpStatusConflict                     = 409 // RFC 9110, 15.5.10
	HttpStatusGone                         = 410 // RFC 9110, 15.5.11
	HttpStatusLengthRequired               = 411 // RFC 9110, 15.5.12
	HttpStatusPreconditionFailed           = 412 // RFC 9110, 15.5.13
	HttpStatusRequestEntityTooLarge        = 413 // RFC 9110, 15.5.14
	HttpStatusRequestURITooLong            = 414 // RFC 9110, 15.5.15
	HttpStatusUnsupportedMediaType         = 415 // RFC 9110, 15.5.16
	HttpStatusRequestedRangeNotSatisfiable = 416 // RFC 9110, 15.5.17
	HttpStatusExpectationFailed            = 417 // RFC 9110, 15.5.18
	HttpStatusTeapot                       = 418 // RFC 9110, 15.5.19 (Unused)
	HttpStatusMisdirectedRequest           = 421 // RFC 9110, 15.5.20
	HttpStatusUnprocessableEntity          = 422 // RFC 9110, 15.5.21
	HttpStatusLocked                       = 423 // RFC 4918, 11.3
	HttpStatusFailedDependency             = 424 // RFC 4918, 11.4
	HttpStatusTooEarly                     = 425 // RFC 8470, 5.2.
	HttpStatusUpgradeRequired              = 426 // RFC 9110, 15.5.22
	HttpStatusPreconditionRequired         = 428 // RFC 6585, 3
	HttpStatusTooManyRequests              = 429 // RFC 6585, 4
	HttpStatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	HttpStatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	HttpStatusInternalServerError           = 500 // RFC 9110, 15.6.1
	HttpStatusNotImplemented                = 501 // RFC 9110, 15.6.2
	HttpStatusBadGateway                    = 502 // RFC 9110, 15.6.3
	HttpStatusServiceUnavailable            = 503 // RFC 9110, 15.6.4
	HttpStatusGatewayTimeout                = 504 // RFC 9110, 15.6.5
	HttpStatusHTTPVersionNotSupported       = 505 // RFC 9110, 15.6.6
	HttpStatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	HttpStatusInsufficientStorage           = 507 // RFC 4918, 11.5
	HttpStatusLoopDetected                  = 508 // RFC 5842, 7.2
	HttpStatusNotExtended                   = 510 // RFC 2774, 7
	HttpStatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

var HttpStatus = struct {
	Continue           int
	SwitchingProtocols int
	Processing         int
	EarlyHints         int

	OK                          int
	Created                     int
	Accepted                    int
	NonAuthoritativeInformation int
	NoContent                   int
	ResetContent                int
	PartialContent              int
	MultiStatus                 int
	AlreadyReported             int
	IMUsed                      int

	MultipleChoices   int
	MovedPermanently  int
	Found             int
	SeeOther          int
	NotModified       int
	UseProxy          int
	SwitchProxy       int
	TemporaryRedirect int
	PermanentRedirect int

	BadRequest                   int
	Unauthorized                 int
	PaymentRequired              int
	Forbidden                    int
	NotFound                     int
	MethodNotAllowed             int
	NotAcceptable                int
	ProxyAuthRequired            int
	RequestTimeout               int
	Conflict                     int
	Gone                         int
	LengthRequired               int
	PreconditionFailed           int
	RequestEntityTooLarge        int
	RequestURITooLong            int
	UnsupportedMediaType         int
	RequestedRangeNotSatisfiable int
	ExpectationFailed            int
	Teapot                       int
	MisdirectedRequest           int
	UnprocessableEntity          int
	Locked                       int
	FailedDependency             int
	TooEarly                     int
	UpgradeRequired              int
	PreconditionRequired         int
	TooManyRequests              int
	RequestHeaderFieldsTooLarge  int
	UnavailableForLegalReasons   int

	InternalServerError           int
	NotImplemented                int
	BadGateway                    int
	ServiceUnavailable            int
	GatewayTimeout                int
	HTTPVersionNotSupported       int
	VariantAlsoNegotiates         int
	InsufficientStorage           int
	LoopDetected                  int
	NotExtended                   int
	NetworkAuthenticationRequired int
}{
	Continue:           HttpStatusContinue,
	SwitchingProtocols: HttpStatusSwitchingProtocols,
	Processing:         HttpStatusProcessing,
	EarlyHints:         HttpStatusEarlyHints,

	OK:                          HttpStatusOK,
	Created:                     HttpStatusCreated,
	Accepted:                    HttpStatusAccepted,
	NonAuthoritativeInformation: HttpStatusNonAuthoritativeInformation,
	NoContent:                   HttpStatusNoContent,
	ResetContent:                HttpStatusResetContent,
	PartialContent:              HttpStatusPartialContent,
	MultiStatus:                 HttpStatusMultiStatus,
	AlreadyReported:             HttpStatusAlreadyReported,
	IMUsed:                      HttpStatusIMUsed,

	MultipleChoices:   HttpStatusMultipleChoices,
	MovedPermanently:  HttpStatusMovedPermanently,
	Found:             HttpStatusFound,
	SeeOther:          HttpStatusSeeOther,
	NotModified:       HttpStatusNotModified,
	UseProxy:          HttpStatusUseProxy,
	SwitchProxy:       HttpStatusSwitchProxy,
	TemporaryRedirect: HttpStatusTemporaryRedirect,
	PermanentRedirect: HttpStatusPermanentRedirect,

	BadRequest:                   HttpStatusBadRequest,
	Unauthorized:                 HttpStatusUnauthorized,
	PaymentRequired:              HttpStatusPaymentRequired,
	Forbidden:                    HttpStatusForbidden,
	NotFound:                     HttpStatusNotFound,
	MethodNotAllowed:             HttpStatusMethodNotAllowed,
	NotAcceptable:                HttpStatusNotAcceptable,
	ProxyAuthRequired:            HttpStatusProxyAuthRequired,
	RequestTimeout:               HttpStatusRequestTimeout,
	Conflict:                     HttpStatusConflict,
	Gone:                         HttpStatusGone,
	LengthRequired:               HttpStatusLengthRequired,
	PreconditionFailed:           HttpStatusPreconditionFailed,
	RequestEntityTooLarge:        HttpStatusRequestEntityTooLarge,
	RequestURITooLong:            HttpStatusRequestURITooLong,
	UnsupportedMediaType:         HttpStatusUnsupportedMediaType,
	RequestedRangeNotSatisfiable: HttpStatusRequestedRangeNotSatisfiable,
	ExpectationFailed:            HttpStatusExpectationFailed,
	Teapot:                       HttpStatusTeapot,
	MisdirectedRequest:           HttpStatusMisdirectedRequest,
	UnprocessableEntity:          HttpStatusUnprocessableEntity,
	Locked:                       HttpStatusLocked,
	FailedDependency:             HttpStatusFailedDependency,
	TooEarly:                     HttpStatusTooEarly,
	UpgradeRequired:              HttpStatusUpgradeRequired,
	PreconditionRequired:         HttpStatusPreconditionRequired,
	TooManyRequests:              HttpStatusTooManyRequests,
	RequestHeaderFieldsTooLarge:  HttpStatusRequestHeaderFieldsTooLarge,
	UnavailableForLegalReasons:   HttpStatusUnavailableForLegalReasons,

	InternalServerError:           HttpStatusInternalServerError,
	NotImplemented:                HttpStatusNotImplemented,
	BadGateway:                    HttpStatusBadGateway,
	ServiceUnavailable:            HttpStatusServiceUnavailable,
	GatewayTimeout:                HttpStatusGatewayTimeout,
	HTTPVersionNotSupported:       HttpStatusHTTPVersionNotSupported,
	VariantAlsoNegotiates:         HttpStatusVariantAlsoNegotiates,
	InsufficientStorage:           HttpStatusInsufficientStorage,
	LoopDetected:                  HttpStatusLoopDetected,
	NotExtended:                   HttpStatusNotExtended,
	NetworkAuthenticationRequired: HttpStatusNetworkAuthenticationRequired,
}

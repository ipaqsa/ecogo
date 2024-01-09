package ecogo

const (
	userAgent = "eco-client"
	mediaType = "application/json"

	internalHeaderRetryAttempts = "X-ECO-Retry-Attempts"

	defaultRetryMax     = 3
	defaultRetryWaitMax = 30
	defaultRetryWaitMin = 1
)

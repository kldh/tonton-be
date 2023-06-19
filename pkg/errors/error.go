package errors

type reason string

const (
	DomainUM    = "UM"
	DomainOrder = "ORDER"
)

const (
	ReasonSuccess       = "SUCCESS"
	ReasonInternalError = "INTERNAL_ERROR"
)

package main

import (
	"errors"
	"fmt"
)

var (
	ErrOpsNotFound   = errors.New("operations record not found")
	ErrOpsConflict   = errors.New("operations revision conflict")
	ErrOpsInvalid    = errors.New("operations request is invalid")
	ErrOpsTransition = errors.New("operations status transition is not allowed")
	ErrOpsPolicy     = errors.New("operations policy rejected the request")
)

// OpsError wraps an underlying sentinel with the operation that produced it.
// Cause is preserved via Unwrap so errors.Is can traverse the chain even when
// the sentinel was returned by a lower layer (store, state machine, policy).
type OpsError struct {
	Code      string
	Operation string
	Cause     error
}

func (e *OpsError) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Operation, e.Cause)
}
func (e *OpsError) Unwrap() error { return e.Cause }

// wrapOps preserves the sentinel chain: the cause is reachable via errors.Is,
// so callers can still distinguish not_found / conflict / transition / policy
// after the service has annotated the error with the failing operation.
func wrapOps(code, operation string, cause error) error {
	return &OpsError{Code: code, Operation: operation, Cause: cause}
}

// opsCode classifies an error into a stable string code for the response.
// errors.Is is required because store/state/policy layers wrap the sentinels.
func opsCode(err error) string {
	switch {
	case errors.Is(err, ErrOpsNotFound):
		return "not_found"
	case errors.Is(err, ErrOpsConflict):
		return "conflict"
	case errors.Is(err, ErrOpsInvalid):
		return "invalid"
	case errors.Is(err, ErrOpsTransition):
		return "transition"
	case errors.Is(err, ErrOpsPolicy):
		return "policy"
	default:
		return "internal"
	}
}
func opsIsNotFound(err error) bool   { return errors.Is(err, ErrOpsNotFound) }
func opsIsConflict(err error) bool   { return errors.Is(err, ErrOpsConflict) }
func opsIsInvalid(err error) bool    { return errors.Is(err, ErrOpsInvalid) }
func opsIsTransition(err error) bool { return errors.Is(err, ErrOpsTransition) }
func opsIsPolicy(err error) bool     { return errors.Is(err, ErrOpsPolicy) }

// opsHTTPStatus maps a classified error code to the matching HTTP status.
// Conflict and illegal transitions are 409; policy rejection is 422.
func opsHTTPStatus(code string) int {
	switch code {
	case "not_found":
		return 404
	case "conflict", "transition":
		return 409
	case "invalid":
		return 400
	case "policy":
		return 422
	default:
		return 500
	}
}

package model

import "errors"

var (
	ErrInternal = errors.New("internal error")

	// * Data errors.
	ErrDataNotFound    = errors.New("data not found")
	ErrNoUpdatedData   = errors.New("no data to update")
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")

	// * Communication errors.
	ErrBadRequest = errors.New("error parsing request")

	// * Authority errors.
	ErrForbidden = errors.New("user is forbidden to access the resource")
)

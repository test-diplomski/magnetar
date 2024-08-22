package domain

import "errors"

var (
	ErrNodeClaimed = errors.New("node has been already claimed and is not in the node pool anymore")
	ErrServerSide  = errors.New("an unexpected server-side error occurred")
	ErrForbidden   = errors.New("you are not authorized to perform this operation")
)

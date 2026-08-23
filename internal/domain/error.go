package domain

import "errors"

var (
	ErrContractNotFound                 = errors.New("contract not found")
	ErrForbidden                        = errors.New("forbidden")
	ErrConflict                         = errors.New("conflict")
	ErrAlreadyExists                    = errors.New("already exists")
	ErrNotAllowedCurrentStatusToPublish = errors.New("current status isn't allowed to publish")
	ErrClientNotDriver                  = errors.New("client isn't driver of this contract")
	ErrStatusIsPublishedAlready         = errors.New("status is published already")
)

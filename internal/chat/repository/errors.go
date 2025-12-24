package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrSameParticipant = errors.New("Same participant")
var ErrDuplicateMessage = errors.New("Duplicate message")
package model

import (
	"errors"
	"strings"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

func NewPriority(p string) (Priority, error) {
	p = strings.ToUpper(p)
	if !Priority(p).Valid() {
		return "", errors.New("invalid priority")
	}
	return Priority(p), nil
}

func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

func (s Priority) String() string {
	return string(s)
}

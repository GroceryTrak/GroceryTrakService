package models

import "time"

type Priority int

const (
	DefaultPriority Priority = 0
	HighPriority    Priority = 1
	LowPriority     Priority = -1
)

type QueueItem struct {
	ItemID    uint      `json:"item_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Priority  Priority  `json:"priority"`
}

package storage

import "github.com/Quard/poindexter/internal/reading_list"

type Storage interface {
	Add(*reading_list.Record) error
	List(userID string, processor func(record *reading_list.Record)) error
	MarkAsRead(userID string, ID string) (reading_list.Record, error)
}

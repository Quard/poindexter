package internal_api

import (
	context "context"

	"github.com/Quard/poindexter/internal/reading_list"
	"github.com/getsentry/sentry-go"
)

func (srv internalAPIServer) Add(ctx context.Context, url *URL) (*ReadingListRecord, error) {
	record, err := reading_list.NewRecordByURL(url.UserID, url.Url)
	if err != nil {
		return &ReadingListRecord{}, err
	}
	if err := srv.storage.Add(&record); err != nil {
		return &ReadingListRecord{}, err
	}

	respRecord := &ReadingListRecord{
		Id:       record.ID.Hex(),
		UserID:   record.UserID.Hex(),
		Title:    record.Title,
		Url:      record.URL,
		ImageUrl: record.ImageURL,
		Created:  record.Created.Unix(),
		IsRead:   record.IsRead,
	}

	return respRecord, nil
}

func (srv internalAPIServer) List(user *User, stream InternalAPI_ListServer) error {
	return srv.storage.List(user.ID, func(record *reading_list.Record) {
		err := stream.Send(&ReadingListRecord{
			Id:       record.ID.Hex(),
			Title:    record.Title,
			Url:      record.URL,
			ImageUrl: record.ImageURL,
			Created:  record.Created.Unix(),
			IsRead:   record.IsRead,
		})
		if err != nil {
			sentry.CaptureException(err)
		}
	})
}

func (srv internalAPIServer) MarkAsRead(ctx context.Context, id *ID) (*ReadingListRecord, error) {
	record, err := srv.storage.MarkAsRead(id.UserID, id.Id)

	return convertToReadingListRecord(record), err
}

func (srv internalAPIServer) Del(ctx context.Context, id *ID) (*ReadingListRecord, error) {
	return &ReadingListRecord{}, nil
}

func convertToReadingListRecord(record reading_list.Record) *ReadingListRecord {
	return &ReadingListRecord{
		Id:       record.ID.Hex(),
		UserID:   record.UserID.Hex(),
		Title:    record.Title,
		Url:      record.URL,
		ImageUrl: record.ImageURL,
		Created:  record.Created.Unix(),
		IsRead:   record.IsRead,
	}
}

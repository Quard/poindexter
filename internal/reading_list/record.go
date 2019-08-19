package reading_list

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	readability "github.com/philipjkim/goreadability"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   primitive.ObjectID `bson:"user_id"`
	Title    string             `bson:"title"`
	URL      string             `bson:"url"`
	ImageURL string             `bson:"image_url"`
	Created  time.Time          `bson:"created"`
	IsRead   bool               `bson:"is_read"`
}

func NewRecordByURL(userID string, url string) (Record, error) {
	userObjectID, errID := primitive.ObjectIDFromHex(userID)
	if errID != nil {
		sentry.CaptureException(errID)
		return Record{}, errors.New("unable to parse user id")
	}

	record := Record{
		UserID:  userObjectID,
		URL:     url,
		Created: time.Now(),
		IsRead:  false,
	}

	content, err := readability.Extract(url, readability.NewOption())
	if err != nil {
		sentry.CaptureException(err)
		return record, errors.New("unable to retrieve info by given URL")
	}

	record.Title = content.Title
	record.ImageURL = content.Images[0].URL

	return record, nil
}

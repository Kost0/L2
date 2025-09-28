package business

import (
	"L2/18/models"
	"encoding/json"
	"errors"
	"time"
)

func NewEvent(buf []byte) (*models.Event, error) {
	event := &models.Event{}

	err := json.Unmarshal(buf, event)
	if err != nil {
		return nil, err
	}

	_, err = time.Parse("2006-01-02", event.Date)
	if err != nil {
		return nil, err
	}

	models.Events = append(models.Events, event)

	return event, nil
}

func UpdateEvent(buf []byte) (*models.Event, error) {
	event := &models.Event{}

	err := json.Unmarshal(buf, event)
	if err != nil {
		return nil, err
	}

	_, err = time.Parse("2006-01-02", event.Date)
	if err != nil {
		return nil, err
	}

	for i, ev := range models.Events {
		if ev.UserID == event.UserID {
			if ev.Event == event.Event || ev.Date == event.Date {
				models.Events[i] = event
				return event, nil
			}
		}
	}

	return nil, errors.New("event not found")
}

func DeleteEvent(buf []byte) (*models.Event, error) {
	event := &models.Event{}

	err := json.Unmarshal(buf, event)
	if err != nil {
		return nil, err
	}

	for i, ev := range models.Events {
		if *ev == *event {
			models.Events = append(models.Events[:i], models.Events[i+1:]...)
			return event, nil
		}
	}

	return nil, errors.New("event not found")
}

func FindEventForTime(date time.Time, userID string, duration time.Duration) ([]*models.Event, error) {
	necessaryEvents := make([]*models.Event, 0)

	timeFromNow := date.Add(duration)

	for _, event := range models.Events {
		if event.UserID == userID {
			eventDate, err := time.Parse("2006-01-02", event.Date)
			if err != nil {
				return nil, err
			}
			if (eventDate.Equal(date) || eventDate.After(date)) && (eventDate.Equal(timeFromNow) || eventDate.Before(timeFromNow)) {
				necessaryEvents = append(necessaryEvents, event)
			}
		}
	}

	return necessaryEvents, nil
}

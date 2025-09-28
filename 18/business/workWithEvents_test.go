package business

import (
	"L2/18/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent_Success(t *testing.T) {
	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	newEvent, err := NewEvent(buf)
	assert.NoError(t, err)
	assert.Equal(t, event, newEvent)
}

func TestNewEvent_WrongDate(t *testing.T) {
	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-74",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	_, err = NewEvent(buf)
	assert.Error(t, err)
}

func TestNewEvent_WrongStruct(t *testing.T) {
	event := struct {
		Number int
	}{
		Number: 1,
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	_, err = NewEvent(buf)
	assert.Error(t, err)
}

func TestUpdateEvent_Success(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}}

	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-10",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := UpdateEvent(buf)
	assert.NoError(t, err)
	assert.Equal(t, event, updatedEvent)
}

func TestUpdateEvent_WrongDate(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}}

	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-74",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := UpdateEvent(buf)
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)
}

func TestUpdateEvent_NothingToUpdate(t *testing.T) {
	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-74",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := UpdateEvent(buf)
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)
}

func TestUpdateEvent_WrongStruct(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}}

	event := struct {
		Number int
	}{
		Number: 1,
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	updatedEvent, err := UpdateEvent(buf)
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)
}

func TestDeleteEvent_Success(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}}

	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}

	buf, err := json.Marshal(models.Events[0])
	assert.NoError(t, err)

	deletedEvent, err := DeleteEvent(buf)
	assert.NoError(t, err)
	assert.Equal(t, event, deletedEvent)
}

func TestDeleteEvent_NothingToDelete(t *testing.T) {
	event := &models.Event{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}

	buf, err := json.Marshal(event)
	assert.NoError(t, err)

	deletedEvent, err := DeleteEvent(buf)
	assert.Error(t, err)
	assert.Nil(t, deletedEvent)
}

func TestFindEventForTime_Success_Day(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-09",
		Event:  "test event",
	}}

	date, err := time.Parse("2006-01-02", models.Events[0].Date)
	assert.NoError(t, err)

	events, err := FindEventForTime(date, "1", time.Hour*24)
	assert.NoError(t, err)
	assert.Equal(t, models.Events, events)
}

func TestFindEventForTime_Success_Week(t *testing.T) {
	models.Events = []*models.Event{
		{
			UserID: "1",
			Date:   "2025-09-09",
			Event:  "test event",
		},
		{
			UserID: "1",
			Date:   "2025-09-10",
			Event:  "test event",
		},
	}

	date, err := time.Parse("2006-01-02", models.Events[0].Date)
	assert.NoError(t, err)

	events, err := FindEventForTime(date, "1", time.Hour*24*7)
	assert.NoError(t, err)
	assert.Equal(t, models.Events, events)
}

func TestFindEventForTime_Success_Month(t *testing.T) {
	models.Events = []*models.Event{
		{
			UserID: "1",
			Date:   "2025-09-09",
			Event:  "test event",
		},
		{
			UserID: "1",
			Date:   "2025-09-10",
			Event:  "test event",
		},
		{
			UserID: "1",
			Date:   "2025-10-01",
			Event:  "test event",
		},
	}

	date, err := time.Parse("2006-01-02", models.Events[0].Date)
	assert.NoError(t, err)

	events, err := FindEventForTime(date, "1", time.Hour*24*30)
	assert.NoError(t, err)
	assert.Equal(t, models.Events, events)
}

func TestFindEventForTime_WrongDate(t *testing.T) {
	models.Events = []*models.Event{{
		UserID: "1",
		Date:   "2025-09-74",
		Event:  "test event",
	}}

	date := time.Now()

	events, err := FindEventForTime(date, "1", time.Hour*24)
	assert.Error(t, err)
	assert.Nil(t, events)
}

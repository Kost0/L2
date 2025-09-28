package handlers

import (
	"L2/18/business"
	"L2/18/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type result struct {
	Res []*models.Event `json:"result"`
}

func newResult(res []*models.Event, w http.ResponseWriter) {
	newRes := result{Res: res}

	buf, err := json.Marshal(&newRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type eventErr struct {
	ErrString string `json:"error"`
}

func newEventErr(err error, w http.ResponseWriter) {
	evErr := &eventErr{ErrString: err.Error()}

	buf, err := json.Marshal(evErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if strings.Contains(evErr.ErrString, "not found") {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.WriteHeader(http.StatusBadRequest)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newEventErr(err, w)
		return
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			newEventErr(err, w)
			return
		}
	}()

	event, err := business.NewEvent(buf)
	if err != nil {
		newEventErr(err, w)
	}

	newResult([]*models.Event{event}, w)

	w.WriteHeader(http.StatusOK)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newEventErr(err, w)
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	event, err := business.UpdateEvent(buf)
	if err != nil {
		newEventErr(err, w)
	}

	newResult([]*models.Event{event}, w)

	w.WriteHeader(http.StatusOK)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newEventErr(err, w)
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	event, err := business.DeleteEvent(buf)
	if err != nil {
		newEventErr(err, w)
	}

	newResult([]*models.Event{event}, w)

	w.WriteHeader(http.StatusOK)
}

func GetEventForDay(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
	}

	queryParams := parsedURL.Query()

	date := queryParams.Get("date")
	user := queryParams.Get("user_id")

	layout := "2006-01-02"

	day, err := time.Parse(layout, date)
	if err != nil {
		newEventErr(err, w)
	}

	necessaryEvents, err := business.FindEventForTime(day, user, time.Hour*24)
	if err != nil {
		newEventErr(err, w)
	}

	newResult(necessaryEvents, w)

	w.WriteHeader(http.StatusOK)
}

func GetEventForWeek(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
	}

	queryParams := parsedURL.Query()

	date := queryParams.Get("date")
	user := queryParams.Get("user_id")

	layout := "2006-01-02"

	day, err := time.Parse(layout, date)
	if err != nil {
		newEventErr(err, w)
	}

	necessaryEvents, err := business.FindEventForTime(day, user, time.Hour*24*7)
	if err != nil {
		newEventErr(err, w)
	}
	newResult(necessaryEvents, w)

	w.WriteHeader(http.StatusOK)
}

func GetEventForMonth(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
	}

	queryParams := parsedURL.Query()

	date := queryParams.Get("date")
	user := queryParams.Get("user_id")

	layout := "2006-01-02"

	day, err := time.Parse(layout, date)
	if err != nil {
		newEventErr(err, w)
	}

	necessaryEvents, err := business.FindEventForTime(day, user, time.Hour*24*30)
	if err != nil {
		newEventErr(err, w)
	}

	newResult(necessaryEvents, w)

	w.WriteHeader(http.StatusOK)
}

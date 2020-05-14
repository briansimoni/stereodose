package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/sessions"
)

// FeedbackController has REST methods to perform operations
// on the Feedback resource
type FeedbackController struct {
	DB    *models.StereoDoseDB
	Store *sessions.CookieStore
}

// NewFeedbackController returns a pointer to NewFeedbackController
func NewFeedbackController(db *models.StereoDoseDB, store *sessions.CookieStore) *FeedbackController {
	return &FeedbackController{DB: db, Store: store}
}

// CreateFeedback will take a JSON payload of Feedback and write to the database
// The intention is that we leave off the authentication middleware because
// I don't want to require that users are authenticated to provide feedback
// But if they are, I am going to record who they are and other user-specific attributes
func (f *FeedbackController) CreateFeedback(w http.ResponseWriter, r *http.Request) error {
	var feedback = &models.Feedback{}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(feedback)
	if err != nil {
		return &util.StatusError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	feedback.DetectedUserAgent = r.Header.Get("user-agent")

	session, err := f.Store.Get(r, "stereodose_session")
	if err == nil {
		ID, ok := session.Values["User_ID"].(uint)
		if ok {
			feedback.UserID = ID
		}
	}

	err = f.DB.Feedback.CreateFeedback(feedback)
	if err != nil {
		return err
	}
	return nil
}

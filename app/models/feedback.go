package models

import "github.com/jinzhu/gorm"

// FeedbackService is an interface used to describe all of the behavior
// of some kind of service that a "Feedback Service" should offer
// Useful for mocks/fakes when unit testing
type FeedbackService interface {
	CreateFeedback(feedback *Feedback) error
}

// Feedback is the data structure that data submitted from user surveys
// The idea is that Stereodose applications like the web app and iOS app
// can provide views where we can survey the users and ask about their experience.
// GORM adds an "s" at the end of the table name even though the plural form of feedback is feedback
type Feedback struct {
	gorm.Model
	UserID            uint   `json:"userID"`
	DetectedUserAgent string `json:"detectedUserAgent"`
	OtherComments     string `json:"otherComments"`
	GoodExperience    bool   `json:"goodExperience"`
}

// StereodoseFeedbackService is an implementation of Feedback Service
type StereodoseFeedbackService struct {
	db *gorm.DB
}

// CreateFeedback will save feedback to the database
func (s *StereodoseFeedbackService) CreateFeedback(feedback *Feedback) error {
	return s.db.Create(feedback).Error
}

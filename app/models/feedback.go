package models

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

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
	OtherComments     string `json:"otherComments" gorm:"type:varchar(10000);"`
	GoodExperience    bool   `json:"goodExperience"`
}

// StereodoseFeedbackService is an implementation of Feedback Service
type StereodoseFeedbackService struct {
	db  *gorm.DB
	SNS *sns.SNS
}

// NewFeedbackService will create a new FeedBack service and return a pointer
func NewFeedbackService(db *gorm.DB) *StereodoseFeedbackService {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	SNS := sns.New(session)
	return &StereodoseFeedbackService{
		SNS: SNS,
		db:  db,
	}
}

// CreateFeedback will save feedback to the database
// It will also attempt to send a message to the Stereodose OPS SNS ARN.
// If it fails to send to SNS it will only log warnings. It will not report errors to end users
func (s *StereodoseFeedbackService) CreateFeedback(feedback *Feedback) error {
	err := s.db.Create(feedback).Error
	if err != nil {
		return err
	}

	m, err := json.MarshalIndent(feedback, "", "	")
	if err != nil {
		log.Warn("Unable to JSON marshal for publishing to SNS")
	}

	_, err = s.SNS.Publish(&sns.PublishInput{
		Message:          aws.String(string(m)),
		TopicArn:         aws.String("arn:aws:sns:us-east-1:502859415194:stereodose-ops"),
	})
	if err != nil {
		log.Warn(err.Error())
	}
	return nil
}

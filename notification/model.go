package common

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// NotificationType represents the type of a notification.
type NotificationType string

const (
	NotificationTypeInfo     NotificationType = "info"
	NotificationTypeWarning  NotificationType = "warning"
	NotificationTypeCritical NotificationType = "critical"
)

// Notification represents a notification sent from the application.
type Notification struct {
	gorm.Model
	Time              time.Time         `gorm:"not null"`      // Time when the notification was sent.
	Type              NotificationType  `gorm:"not null"`      // Type of the notification (e.g., info, warning, critical).
	Recipient         string            `gorm:"not null"`      // Recipient of the notification (e.g., email address, phone number).
	Message           string            `gorm:"not null"`      // Message content of the notification.
	IsRead            bool              `gorm:"default:false"` // Flag indicating if the notification has been read.
	IsArchived        bool              `gorm:"default:false"` // Flag indicating if the notification is archived.
	From              uuid.UUID         // ID of the user who sent the notification.
	To                uuid.UUID         // ID of the user who received the notification.
	AttachmentIDs     []int64           // IDs of attachments associated with the notification.
	RelatedCalendarID int64             // ID of the calendar related to the notification (if applicable).
	RelatedEventID    int64             // ID of the event related to the notification (if applicable).
	Metadata          map[string]string `gorm:"type:json"` // Additional metadata associated with the notification.
}

// MarkAsRead marks the notification as read.
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// Archive archives the notification.
func (n *Notification) Archive() {
	n.IsArchived = true
}

// AddAttachmentID adds a new attachment ID to the notification.
func (n *Notification) AddAttachmentID(attachmentID int64) {
	n.AttachmentIDs = append(n.AttachmentIDs, attachmentID)
}

// RemoveAttachmentID removes an attachment ID from the notification based on the provided index.
func (n *Notification) RemoveAttachmentID(index int) {
	if index >= 0 && index < len(n.AttachmentIDs) {
		n.AttachmentIDs = append(n.AttachmentIDs[:index], n.AttachmentIDs[index+1:]...)
	}
}

// SetRelatedCalendarID sets the ID of the related calendar for the notification.
func (n *Notification) SetRelatedCalendarID(calendarID int64) {
	n.RelatedCalendarID = calendarID
}

// SetRelatedEventID sets the ID of the related event for the notification.
func (n *Notification) SetRelatedEventID(eventID int64) {
	n.RelatedEventID = eventID
}

// AddMetadata adds metadata to the notification.
func (n *Notification) AddMetadata(key, value string) {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string)
	}
	n.Metadata[key] = value
}

// Validate checks if the notification struct is valid and returns an error if any validation rule is not satisfied.
func (n *Notification) Validate() error {
	validate := validator.New()
	err := validate.Struct(n)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
		return errors.New("validation failed: " + strings.Join(validationErrors, ", "))
	}
	return nil
}

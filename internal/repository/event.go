package repository

import (
	"errors"
	"gorm.io/gorm"

	"etna-notification/internal/database"
	"etna-notification/internal/domain"
)

type ICalendarEventRepository interface {
	Save(event *domain.CalendarEvent) (*domain.CalendarEvent, error)
	FindAll() ([]*domain.CalendarEvent, error)
	IsNotified(event *domain.CalendarEvent) (bool, error)
	Migrate() error
}

type calendarEventRepository struct {
	database.Client
}

func NewCalendarEventRepository(client database.Client) ICalendarEventRepository {
	return &calendarEventRepository{
		client,
	}
}

func (nr *calendarEventRepository) Save(event *domain.CalendarEvent) (*domain.CalendarEvent, error) {
	return event, nr.DB.Create(event).Error
}

func (nr *calendarEventRepository) FindAll() ([]*domain.CalendarEvent, error) {
	var events []*domain.CalendarEvent
	err := nr.DB.Find(&events).Error
	return events, err
}

func (nr *calendarEventRepository) IsNotified(event *domain.CalendarEvent) (bool, error) {
	result := nr.DB.Take(&event, "external_id = ? and user_id = ?", event.ExternalID, event.UserID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	return true, result.Error
}

func (nr *calendarEventRepository) Migrate() error {
	return nr.DB.AutoMigrate(&domain.CalendarEvent{})
}

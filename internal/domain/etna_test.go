package domain

import (
	"testing"

	"gorm.io/gorm"

	"github.com/maxatome/go-testdeep/td"
)

func TestBuildMessageFromEtnaCalendarEvent(t *testing.T) {
	event := EtnaCalendarEvent{
		ID:           1056971,
		Event:        22683,
		Name:         "Suivi relatif à l'étape 5",
		ActivityName: "Management de la qualité 1ère année",
		SessionName:  "2023_Master - Octobre_CMG-MGQ5_10_0",
		Type:         "suivi",
		Location:     "Contacter 5 minutes avant via Google Chat salle G15 pour isolement",
		Start:        "2023-02-17 11:00:00",
		End:          "2023-02-17 11:10:00",
		Group:        EtnaCalendarEventGroup{},
		Registration: EtnaCalendarEventRegistration{},
		UvName:       "CMG-MGQ5",
	}

	td.Cmp(t,
		event.BuildCalendarMessage(),
		":date: **CMG-MGQ5** Suivi relatif à l'étape 5 : Management de la qualité 1ère année. "+
			"Contacter 5 minutes avant via Google Chat salle G15 pour isolement "+
			": 2023-02-17 11:00:00 - 2023-02-17 11:10:00")
}

func TestBuildCalendarEventWithUUIDAsId(t *testing.T) {
	event := EtnaCalendarEvent{
		ID:           "uuid",
		Event:        22683,
		Name:         "Suivi relatif à l'étape 5",
		ActivityName: "Management de la qualité 1ère année",
		SessionName:  "2023_Master - Octobre_CMG-MGQ5_10_0",
		Type:         "suivi",
		Location:     "Contacter 5 minutes avant via Google Chat salle G15 pour isolement",
		Start:        "2023-02-17 11:00:00",
		End:          "2023-02-17 11:10:00",
		Group:        EtnaCalendarEventGroup{},
		Registration: EtnaCalendarEventRegistration{},
		UvName:       "CMG-MGQ5",
	}

	td.Cmp(t, event.BuildCalendarEvent(&User{
		Model: gorm.Model{
			ID: 1,
		},
	}), &CalendarEvent{
		ExternalID: "uuid",
		UserID:     1,
	})
}

func TestBuildCalendarEventWithIntegerAsId(t *testing.T) {
	event := EtnaCalendarEvent{
		ID:           123,
		Event:        22683,
		Name:         "Suivi relatif à l'étape 5",
		ActivityName: "Management de la qualité 1ère année",
		SessionName:  "2023_Master - Octobre_CMG-MGQ5_10_0",
		Type:         "suivi",
		Location:     "Contacter 5 minutes avant via Google Chat salle G15 pour isolement",
		Start:        "2023-02-17 11:00:00",
		End:          "2023-02-17 11:10:00",
		Group:        EtnaCalendarEventGroup{},
		Registration: EtnaCalendarEventRegistration{},
		UvName:       "CMG-MGQ5",
	}

	td.Cmp(t, event.BuildCalendarEvent(&User{
		Model: gorm.Model{
			ID: 1,
		},
	}), &CalendarEvent{
		ExternalID: "123",
		UserID:     1,
	})
}

func TestIsNotifiableCalendarEvent(t *testing.T) {
	td.CmpTrue(t, EtnaCalendarEvent{Type: "suivi"}.IsNotifiable())
	td.CmpTrue(t, EtnaCalendarEvent{Type: "soutenance"}.IsNotifiable())
	td.CmpFalse(t, EtnaCalendarEvent{Type: "seminaire"}.IsNotifiable())
}

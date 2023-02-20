package domain

import (
	"github.com/maxatome/go-testdeep/td"
	"testing"
	"time"
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

func TestIsNotifiableCalendarEvent(t *testing.T) {
	td.CmpTrue(t, EtnaCalendarEvent{Type: "suivi"}.IsNotifiable())
	td.CmpTrue(t, EtnaCalendarEvent{Type: "soutenance"}.IsNotifiable())
	td.CmpFalse(t, EtnaCalendarEvent{Type: "seminaire"}.IsNotifiable())
}

func TestEtnaCalendarEventIsInNext30Minutes(t *testing.T) {
	now := time.Now()
	eventStart := now.Add(15 * time.Minute)
	td.CmpTrue(t, EtnaCalendarEvent{Start: eventStart.Format("2006-01-02 15:04:05")}.IsInNext30Minutes())
	eventStart = eventStart.Add(14 * time.Minute)
	td.CmpTrue(t, EtnaCalendarEvent{Start: eventStart.Format("2006-01-02 15:04:05")}.IsInNext30Minutes())
	eventStart = eventStart.Add(2 * time.Minute)
	td.CmpFalse(t, EtnaCalendarEvent{Start: eventStart.Format("2006-01-02 15:04:05")}.IsInNext30Minutes())

	eventStart = now.Add(-5 * time.Minute)
	td.CmpFalse(t, EtnaCalendarEvent{Start: eventStart.Format("2006-01-02 15:04:05")}.IsInNext30Minutes())
}

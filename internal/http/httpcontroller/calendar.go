package httpcontroller

import (
	"fmt"
	ics "github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"net/http"

	"etna-notification/internal/controller"
)

type HTTPCalendarController struct {
	controller controller.ICalendarController
}

// NewCalendarController boostrap quiz endpoints with their controllers.
// These routes are located under /quiz prefix.
func NewCalendarController(e *gin.Engine, controller controller.ICalendarController) *HTTPCalendarController {
	httpController := &HTTPCalendarController{
		controller: controller,
	}

	group := e.Group("/calendar")
	group.GET("", httpController.Calendar)

	return httpController
}

func (q *HTTPCalendarController) Calendar(c *gin.Context) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	events := q.controller.GetCalendarEvent()
	for _, calEvent := range events {
		event := cal.AddEvent(fmt.Sprintf("etna-notification-%v", calEvent.ID))
		event.SetStartAt(calEvent.EventStartTime())
		event.SetEndAt(calEvent.EventStopTime())
		event.SetSummary(fmt.Sprintf("%s %s %s", calEvent.UvName, calEvent.Name, calEvent.ActivityName))
		event.SetLocation(calEvent.Location)
		event.SetDescription(calEvent.BuildCalendarMessage())
		event.SetURL(calEvent.Location)
	}

	buffer := cal.Serialize()
	c.Header("Content-Disposition", "attachment; filename=calendar.ics")
	c.Data(http.StatusOK, "application/octet-stream", []byte(buffer))
}

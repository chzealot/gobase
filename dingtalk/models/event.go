package models

import "time"

type EventTime struct {
	Date     string    `json:"date"`
	DateTime time.Time `json:"dateTime"`
	TimeZone string    `json:"timeZone"`
}

type EventRecurrencePattern struct {
	Type       string `json:"type"`
	DayOfMonth int    `json:"dayOfMonth"`
	DaysOfWeek string `json:"daysOfWeek"`
	Index      string `json:"index"`
	Interval   int    `json:"interval"`
}

type EventRecurrenceRange struct {
	Type                string    `json:"type"`
	EndDate             time.Time `json:"endDate"`
	NumberOfOccurrences int       `json:"numberOfOccurrences"`
}

type EventRecurrence struct {
	Pattern EventRecurrencePattern `json:"pattern"`
	Range   EventRecurrenceRange   `json:"range"`
}

type EventAttendee struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	ResponseStatus string `json:"responseStatus"`
	Self           bool   `json:"self"`
	IsOptional     bool   `json:"isOptional"`
}
type EventAttendeeList []EventAttendee

type EventOrganizer struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	ResponseStatus string `json:"responseStatus"`
	Self           bool   `json:"self"`
}

type EventLocation struct {
	DisplayName  string   `json:"displayName"`
	MeetingRooms []string `json:"meetingRooms"`
}

type EventOnlineMeetingInfo struct {
	Type         string `json:"type"`
	ConferenceId string `json:"conferenceId"`
	Url          string `json:"url"`
}

type EventReminder struct {
	Method  string `json:"method"`
	Minutes string `json:"minutes"`
}
type EventReminderList []EventReminder

type EventSharedProperties struct {
	SourceOpenCid string `json:"sourceOpenCid"`
	BelongCorpId  string `json:"belongCorpId"`
}
type EventExtendedProperties struct {
	SharedProperties EventSharedProperties `json:"sharedProperties"`
}

type EventMeetingRoom struct {
	RoomId         string `json:"roomId"`
	ResponseStatus string `json:"responseStatus"`
	DisplayName    string `json:"displayName"`
}
type EventMeetingRoomList []EventMeetingRoom

type EventCategory struct {
	DisplayName string `json:"displayName"`
}
type EventCategoryList []EventCategory

type CalendarEvent struct {
	ID                 string                  `json:"id"`
	Summary            string                  `json:"summary"`
	Description        string                  `json:"description"`
	Start              EventTime               `json:"start"`
	OriginStart        EventTime               `json:"originStart"`
	End                EventTime               `json:"end"`
	IsAllDay           bool                    `json:"isAllDay"`
	Recurrence         EventRecurrence         `json:"recurrence"`
	Attendees          EventAttendeeList       `json:"attendees"`
	Organizer          EventOrganizer          `json:"organizer"`
	Location           EventLocation           `json:"location"`
	SeriesMasterId     string                  `json:"seriesMasterId"`
	CreateTime         time.Time               `json:"createTime"`
	UpdateTime         time.Time               `json:"updateTime"`
	Status             string                  `json:"status"`
	OnlineMeetingInfo  EventOnlineMeetingInfo  `json:"onlineMeetingInfo"`
	Reminders          EventReminderList       `json:"reminders"`
	ExtendedProperties EventExtendedProperties `json:"extendedProperties"`
	MeetingRooms       EventMeetingRoomList    `json:"meetingRooms"`
	Categories         EventCategoryList       `json:"categories"`
}

type CalendarEventList []*CalendarEvent

type EventResponse struct {
	NextToken string            `json:"nextToken"`
	Events    CalendarEventList `json:"events"`
	SyncToken string            `json:"syncToken"`
}

package models

type Calendar struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	TimeZone    string `json:"timeZone"`
	ETag        string `json:"eTag"`
	Type        string `json:"type"`
	Privilege   string `json:"privilege"`
}

type CalendarList []*Calendar

type CalendarResponse struct {
	CalendarOriginResponse *CalendarOriginResponse `json:"response"`
}

type CalendarOriginResponse struct {
	Calendars CalendarList `json:"calendars"`
}

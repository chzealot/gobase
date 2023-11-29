package models

type CreateTodoTaskRequestDetailUrl struct {
	AppUrl string `json:"appUrl"`
	PcUrl  string `json:"pcUrl"`
}
type CreateTodoTaskRequestNotifyConfigs struct {
	DingNotify string `json:"dingNotify"`
}
type CreateTodoTaskRequest struct {
	SourceID           string                         `json:"sourceId"`
	Subject            string                         `json:"subject"`
	CreatorID          string                         `json:"creatorId"`
	Description        string                         `json:"description"`
	DueTime            int64                          `json:"dueTime"`
	ExecutorIds        []string                       `json:"executorIds"`
	ParticipantIds     []string                       `json:"participantIds"`
	DetailUrl          CreateTodoTaskRequestDetailUrl `json:"detailUrl"`
	IsOnlyShowExecutor bool                           `json:"isOnlyShowExecutor"`
	//Priority           int                                `json:"priority"`
	NotifyConfigs CreateTodoTaskRequestNotifyConfigs `json:"notifyConfigs"`
}

type CreateTodoTaskResponse struct {
	ID             string   `json:"id"`
	BizTag         string   `json:"bizTag"`
	CreatedTime    int64    `json:"createdTime"`
	CreatorID      string   `json:"creatorId"`
	Done           bool     `json:"done"`
	DueTime        int64    `json:"dueTime"`
	FinishTime     int      `json:"finishTime"`
	ModifiedTime   int64    `json:"modifiedTime"`
	ModifierId     string   `json:"modifierId"`
	ParticipantIds []string `json:"participantIds"`
	Priority       int      `json:"priority"`
	RequestId      string   `json:"requestId"`
	Source         string   `json:"source"`
	StartTime      int      `json:"startTime"`
	Subject        string   `json:"subject"`
	TenantId       string   `json:"tenantId"`
	TenantType     string   `json:"tenantType"`
}

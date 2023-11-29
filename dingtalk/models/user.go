package models

import "time"

type ContactUser struct {
	Nick      string `json:"nick"`
	AvatarUrl string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

type TopGetByUnionIdResponse struct {
	ContactType int    `json:"contact_type"`
	UserID      string `json:"userid"`
}

type TopUserRole struct {
	GroupName string `json:"group_name"`
	Name      string `json:"name"`
	ID        int    `json:"id"`
}
type TopUserDeptOrder struct {
	DeptID int   `json:"dept_id"`
	Order  int64 `json:"order"`
}
type TopUserLeaderInDept struct {
	Leader bool `json:"leader"`
	DeptID int  `json:"dept_id"`
}
type TopUser struct {
	Extension        string                `json:"extension"`
	Boss             bool                  `json:"boss"`
	UnionID          string                `json:"unionid"`
	RoleList         []TopUserRole         `json:"role_list"`
	ExclusiveAccount bool                  `json:"exclusive_account"`
	Admin            bool                  `json:"admin"`
	Remark           string                `json:"remark"`
	Title            string                `json:"title"`
	Userid           string                `json:"userid"`
	WorkPlace        string                `json:"work_place"`
	DeptOrderList    []TopUserDeptOrder    `json:"dept_order_list"`
	RealAuthed       bool                  `json:"real_authed"`
	DeptIdList       []int                 `json:"dept_id_list"`
	JobNumber        string                `json:"job_number"`
	Email            string                `json:"email"`
	LeaderInDept     []TopUserLeaderInDept `json:"leader_in_dept"`
	CreateTime       time.Time             `json:"create_time"`
	Mobile           string                `json:"mobile"`
	Active           bool                  `json:"active"`
	Telephone        string                `json:"telephone"`
	Avatar           string                `json:"avatar"`
	HideMobile       bool                  `json:"hide_mobile"`
	Senior           bool                  `json:"senior"`
	Name             string                `json:"name"`
	StateCode        string                `json:"state_code"`
}

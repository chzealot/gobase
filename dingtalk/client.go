package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chzealot/gobase/dingtalk/models"
	"github.com/chzealot/gobase/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	url2 "net/url"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 60

type Client struct {
	ClientID     string
	ClientSecret string
	mutex        sync.Mutex
	expireAt     int64
	AccessToken  string
}

func NewDingTalkClient(clientId, clientSecret string) *Client {
	return &Client{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}
}

func (c *Client) GetUserAccessToken(code string) (*models.UserAccessTokenResponse, error) {
	req := models.UserAccessTokenRequest{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Code:         code,
		RefreshToken: "",
		GrantType:    "authorization_code",
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", "https://api.dingtalk.com/v1.0/oauth2/userAccessToken", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := &models.UserAccessTokenResponse{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	resp.ExpireTime = time.Now().Unix() + resp.ExpireIn
	return resp, nil
}

func (c *Client) GetContactUser(token string, unionId string) (*models.ContactUser, error) {
	url := fmt.Sprintf("https://api.dingtalk.com/v1.0/contact/users/%s", url2.QueryEscape(unionId))
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("x-acs-dingtalk-access-token", token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := &models.ContactUser{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetEvents(token string, unionId string) (models.CalendarEventList, error) {
	if unionId == "me" || unionId == "" {
		validUnionId, err := c.GetMyUnionID(token)
		if err != nil {
			return nil, err
		}
		unionId = validUnionId
	}

	cals, err := c.GetCalendars(token, unionId)
	if err != nil {
		return nil, err
	}
	var allEvents models.CalendarEventList
	for _, cal := range cals {
		if cal.Type != "primary" {
			// TODO: 当前仅支持主日历，其他日历暂不考虑
			continue
		}
		events, err := c.GetCalendarEvents(token, unionId, cal.ID)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}
	if len(allEvents) == 0 {
		if err != nil {
			return nil, err
		}
	}
	return allEvents, nil
}

func (c *Client) GetMyUnionID(token string) (string, error) {
	p, err := c.GetContactUser(token, "me")
	if err != nil {
		return "", err
	}
	return p.UnionID, nil
}

func (c *Client) GetCalendars(token string, unionId string) (models.CalendarList, error) {
	if unionId == "me" || unionId == "" {
		validUnionId, err := c.GetMyUnionID(token)
		if err != nil {
			return nil, err
		}
		unionId = validUnionId
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("https://api.dingtalk.com/v1.0/calendar/users/%s/calendars", url2.QueryEscape(unionId)), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("x-acs-dingtalk-access-token", token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := &models.CalendarResponse{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	if resp.CalendarOriginResponse == nil || resp.CalendarOriginResponse.Calendars == nil {
		return nil, errors.New("maybe permission deny")
	}
	return resp.CalendarOriginResponse.Calendars, nil
}

func (c *Client) GetCalendarEvents(token, unionId, calendarId string) (models.CalendarEventList, error) {
	if unionId == "me" || unionId == "" {
		validUnionId, err := c.GetMyUnionID(token)
		if err != nil {
			return nil, err
		}
		unionId = validUnionId
	}

	today := time.Now()
	timeMin := today.Format("2006-01-02") + "T00:00:00+08:00"
	timeMax := today.Format("2006-01-02") + "T23:59:59+08:00"

	url := fmt.Sprintf("https://api.dingtalk.com/v1.0/calendar/users/%s/calendars/%s/events?timeMin=%s&timeMax=%s",
		url2.QueryEscape(unionId), url2.QueryEscape(calendarId),
		url2.QueryEscape(timeMin), url2.QueryEscape(timeMax))
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("x-acs-dingtalk-access-token", token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := &models.EventResponse{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	if resp.Events == nil {
		return nil, errors.New("maybe permission deny")
	}
	return resp.Events, nil
}

func (c *Client) GetAccessToken() (string, error) {
	accessToken := ""
	{
		// 先查询缓存
		c.mutex.Lock()
		now := time.Now().Unix()
		if c.expireAt > 0 && c.AccessToken != "" && (now+60) < c.expireAt {
			// 预留一分钟有效期避免在Token过期的临界点调用接口出现401错误
			accessToken = c.AccessToken
		}
		c.mutex.Unlock()
	}
	if accessToken != "" {
		return accessToken, nil
	}

	tokenResult, err := c.getAccessTokenFromAPI()
	if err != nil {
		return "", err
	}

	{
		// 更新缓存
		c.mutex.Lock()
		c.AccessToken = tokenResult.AccessToken
		c.expireAt = time.Now().Unix() + int64(tokenResult.ExpiresIn)
		c.mutex.Unlock()
	}
	return tokenResult.AccessToken, nil
}

func (c *Client) getAccessTokenFromAPI() (*models.GetTokenResponse, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	const apiUrl = "https://oapi.dingtalk.com/gettoken"
	query := url2.Values{}
	query.Add("appkey", c.ClientID)
	query.Add("appsecret", c.ClientSecret)
	fullUrl := apiUrl + "?" + query.Encode()

	// Send the HTTP request and parse the response body as JSON
	httpClient := http.Client{Timeout: defaultTimeout}
	res, err := httpClient.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &models.GetTokenResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK || response.ErrorCode != 0 {
		logger.Errorw("dingtalk.Client, getAccessTokenFromAPI failed",
			zap.Int("statusCode", res.StatusCode),
			zap.Any("response", response))
		return nil, errors.New(response.ErrorMessage)
	}
	return response, nil
}

func (c *Client) GetUserIDByUnionID(unionId string) (string, error) {
	appAccessToken, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=%s", url2.QueryEscape(appAccessToken))
	params := make(map[string]string, 0)
	params["unionid"] = unionId
	reqBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/json")
	client := http.Client{Timeout: defaultTimeout}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resp := models.TopResult[models.TopGetByUnionIdResponse]{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}

	return resp.Result.UserID, nil
}

func (c *Client) GetUserFromTop(userId string) (*models.TopUser, error) {
	appAccessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/v2/user/get?access_token=%s", url2.QueryEscape(appAccessToken))
	params := make(map[string]string, 0)
	params["userid"] = userId
	reqBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	client := http.Client{Timeout: defaultTimeout}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := models.TopResult[models.TopUser]{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp.Result, nil
}

func (c *Client) CreateTodoTask(creator, subject string, dueTime time.Time) (*models.CreateTodoTaskResponse, error) {
	appAccessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://api.dingtalk.com/v1.0/todo/users/%s/tasks?operatorId=%s",
		url2.QueryEscape(creator),
		url2.QueryEscape(creator))

	req := models.CreateTodoTaskRequest{
		Subject:        subject,
		DueTime:        dueTime.UnixMilli(),
		CreatorID:      creator,
		ExecutorIds:    []string{creator},
		ParticipantIds: []string{creator},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-acs-dingtalk-access-token", appAccessToken)
	client := http.Client{Timeout: defaultTimeout}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := models.CreateTodoTaskResponse{}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

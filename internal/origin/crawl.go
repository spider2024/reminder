package origin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"reminder/etc"
	"reminder/internal/logger"
	"strings"
	"time"
)

const (
	sessionUrl = "/api-getsessionid.json"
	loginUrl   = "/user-login.json?account={0}&password={1}&zentaosid={2}"
	bugUrl     = "/api.php/v1/products/{0}/bugs?limit=26"
	//bugViewUrl  = "/bug-view-{}.html"
	headerToken = "token"
	sessionKey  = "sessionKey"
	userKey     = "userKey"
)

func sessionKeeper(ctx context.Context) (string, error) {
	token, err := etc.Rdb.Get(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		logger.InfoF("Key:%s does not exist\n", sessionKey)
	} else if err != nil {
		logger.ErrorF("loggerin-fail:%s\n", err)
		return "", err
	}
	if token != "" {
		return token, err
	}
	// get sessionId
	resp, err := http.Get(etc.AppConfig.ZenTao.Url + sessionUrl)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorF("close sessionid error: %v\n", err)
		}
	}(resp.Body)
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]any)
	err = json.Unmarshal(all, &m)
	if err != nil {
		logger.ErrorF("unmarchal sessionid error: %v\n", err)
		return "", err
	}
	dataJson, ok := m["data"].(string)
	if !ok {
		logger.ErrorF("unmarchal sessionid data error: %v\n", err)
		return "", errors.New("unmarchal sessionid data error")
	}
	dataMap := make(map[string]any)
	err = json.Unmarshal([]byte(dataJson), &dataMap)
	if err != nil {
		logger.ErrorF("unmarshal data 2 map err:%s\n", err)
		return "", err
	}
	sessionID := dataMap["sessionID"].(string)
	etc.Rdb.Set(ctx, sessionKey, sessionID, time.Hour*8)
	return sessionID, nil
}

func Login(ctx context.Context, username, password string) (userId string, token string, err error) {
	token, err = sessionKeeper(ctx)
	if err != nil {
		return "", "", err
	}
	userId, err = etc.Rdb.Get(ctx, userKey).Result()
	if errors.Is(err, redis.Nil) {
		logger.InfoF("Key:%s does not exist\n", userKey)
	} else if err != nil {
		return "", "", err
	}
	if userId != "" {
		return userId, token, nil
	}
	userId, err = login(ctx, username, password, token)
	if err != nil {
		logger.ErrorF("login.fail:%s", err)
		return "", "", err
	}
	return userId, token, nil
}

func login(ctx context.Context, username, password, token string) (userId string, err error) {
	usernameDone := strings.Replace(etc.AppConfig.ZenTao.Url+loginUrl, "{0}", username, 1)
	passwordDone := strings.Replace(usernameDone, "{1}", password, 1)
	realUrl := strings.Replace(passwordDone, "{2}", token, 1)
	getUserInfoResp, err := http.Get(realUrl)
	body := getUserInfoResp.Body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.InfoF("close getUserInfo error: %v\n", err)
		}
	}(body)
	if err != nil {
		return "", err
	}
	readAll, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	m := make(map[string]any)

	if err := json.Unmarshal(readAll, &m); err != nil {
		logger.ErrorF("Error:", err)
		return "", err
	}
	userMap, ok := m["user"].(map[string]any)
	if !ok {
		logger.ErrorF("body.user.unmarshal err")
		return "", errors.New("user not return")
	} else {
		id, ok := userMap["id"].(string)
		if !ok {
			return "", errors.New("id not return")
		} else {
			etc.Rdb.Set(ctx, userKey, id, time.Hour*8)
			return id, nil
		}
	}
}

type bugs struct {
	Bugs []BugView `json:"bugs"`
}

type BugView struct {
	Id         int        `json:"id"`
	Title      string     `json:"title"`
	Severity   int        `json:"severity"`
	Url        string     `json:"url"`
	OpenedDate time.Time  `json:"openedDate"`
	AppendDate time.Time  `json:"appendDate"`
	AssignedTo AssignedTo `json:"assignedTo"`
	Status     string     `json:"status"`
}

type AssignedTo struct {
	Id           int       `json:"id"`
	Account      string    `json:"account"`
	AssignedDate time.Time `json:"assignedDate"`
}

// Bugs search bugs in project which
func Bugs(token, projectId, userId string) error {
	realUrl := strings.Replace(etc.AppConfig.ZenTao.Url+bugUrl, "{0}", projectId, 1)
	request, err := http.NewRequest("GET", realUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Set(headerToken, token)
	client := &http.Client{}
	do, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorF("close bugs error: %v\n", err)
		}
	}(do.Body)
	body, err := io.ReadAll(do.Body)
	if err != nil {
		logger.ErrorF("Error reading response body:", err)
		return nil
	}
	var bugsList bugs
	//bugsResp := make(map[string]any)
	err = json.Unmarshal(body, &bugsList)
	if err != nil {
		logger.ErrorF("parse bugs resp error: %v\n", err)
		return err
	}
	for _, bug := range bugsList.Bugs {
		if fmt.Sprintf("%s", bug.AssignedTo.Id) == userId {
			logger.InfoF("bugId:%d,account:%s,title:%s,openedDate:%s,status:%s", bug.AssignedTo.Id, bug.AssignedTo.Account, bug.Title, bug.OpenedDate, bug.Status)
		} else {
			sprintf := fmt.Sprintf("bugId:%d,account:%s,title:%s,openedDate:%s,status:%s", bug.AssignedTo.Id, bug.AssignedTo.Account, bug.Title, bug.OpenedDate, bug.Status)
			logger.InfoF(sprintf)
			logger.InfoF("bugId:%d,account:%s,title:%s,openedDate:%s,status:%s", bug.AssignedTo.Id, bug.AssignedTo.Account, bug.Title, bug.OpenedDate, bug.Status)
		}
	}
	//severity
	return nil
}

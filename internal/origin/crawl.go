package origin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reminder/etc"
	"strings"
	"time"
)

const (
	url         = "https://zentao.youpinsanyue.com"
	loginUrl    = "https://zentao.youpinsanyue.com/user-login.json?account={0}&password={1}&zentaosid={2}"
	bugUrl      = "https://zentao.youpinsanyue.com/api.php/v1/products/{0}/bugs?limit=26"
	bugViewUrl  = "https://zentao.youpinsanyue.com/bug-view-{}.html"
	headerToken = "token"
	sessionKey  = "sessionKey"
	userKey     = "userKey"
)

func sessionKeeper(ctx context.Context) (string, error) {
	userId, err := etc.Rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		return "", err
	}
	if userId != "" {
		return userId, err
	}
	// get sessionId
	resp, err := http.Get(url + "/api-getsessionid.json ")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("close sessionid error: %v\n", err)
		}
	}(resp.Body)
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sessionID := string(all)
	etc.Rdb.Set(ctx, sessionKey, sessionID, time.Hour*8)
	return sessionID, nil
}

func Login(ctx context.Context, username, password string) (userId string, token string, err error) {
	token, err = sessionKeeper(ctx)
	if err != nil {
		return "", "", err
	}
	userId, err = etc.Rdb.Get(ctx, userKey).Result()
	if err != nil {
		return "", "", err
	}
	if userId != "" {
		return userId, token, nil
	}
	return login(ctx, username, password)
}

func login(ctx context.Context, username, password string) (userId string, token string, err error) {
	usernameDone := strings.Replace(loginUrl, "{0}", username, 1)
	passwordDone := strings.Replace(usernameDone, "{1}", password, 1)
	realUrl := strings.Replace(passwordDone, "{2}", token, 1)

	getUserInfoResp, err := http.Get(realUrl)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("close getUserInfo error: %v\n", err)
		}
	}(getUserInfoResp.Body)
	if err != nil {
		return "", "", err
	}
	readAll, err := io.ReadAll(getUserInfoResp.Body)
	if err != nil {
		return "", "", err
	}
	var m map[string]any

	if err := json.Unmarshal(readAll, &m); err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}
	userMap, ok := m["user"].(map[string]any)
	if !ok {
		return "", "", errors.New("user not return")
	} else {
		id, ok := userMap["id"].(string)
		if !ok {
			return "", "", errors.New("id not return")
		} else {
			etc.Rdb.Set(ctx, sessionKey, id, time.Hour*8)
			return id, token, nil
		}
	}
}

type BugView struct {
	id         string
	title      string
	severity   int
	url        string
	appendDate string
	//	assignedTo
}

// Bugs search bugs in project which
func Bugs(token, projectId, userId string) error {
	realUrl := strings.Replace(bugUrl, "{0}", projectId, 1)
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
			fmt.Printf("close bugs error: %v\n", err)
		}
	}(do.Body)
	body, err := io.ReadAll(do.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	fmt.Printf(string(body))
	//severity
	return nil
}
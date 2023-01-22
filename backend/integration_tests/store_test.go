package integrationtest_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (suite *BackendTestSuite) Test_GetStoreInfoWithSSE() {
	httpClient := http.Client{Timeout: 3 * time.Second}

	// ========================== open store ==========================
	encodedPassword := base64.StdEncoding.EncodeToString([]byte("im_password"))
	params := map[string]interface{}{
		"email":    "test@gmail.com",
		"password": encodedPassword,
		"name":     "test",
		"timezone": "Asia/Taipei",
		"queue_names": []string{
			"queue_test_1", "queue_test_2",
		},
	}
	jsonParams, _ := json.Marshal(params)
	response, _ := httpClient.Post(suite.ServerBaseURL+"/stores", "application/json", bytes.NewBuffer(jsonParams))

	suite.Equal(200, response.StatusCode)

	var decodedResponse map[string]interface{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)
	suite.Equal(1, int(decodedResponse["id"].(float64)))
	suite.Equal(2, len(decodedResponse["queues"].([]interface{})))

	// ========================== signin store ==========================
	params = map[string]interface{}{
		"email":    "test@gmail.com",
		"password": encodedPassword,
	}
	jsonParams, _ = json.Marshal(params)
	response, _ = httpClient.Post(suite.ServerBaseURL+"/stores/signin", "application/json", bytes.NewBuffer(jsonParams))

	suite.Equal(200, response.StatusCode)

	decodedResponse = map[string]interface{}{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)
	suite.Equal(1, int(decodedResponse["id"].(float64)))
	suite.Equal(2, len(response.Cookies())) // refresh, refreshable, two cookies

	refreshCookie := response.Cookies()[0]

	// ========================== get store jwt token ==========================
	request, _ := http.NewRequest(http.MethodPut, suite.ServerBaseURL+"/stores/token", nil)
	request.AddCookie(refreshCookie)
	response, _ = httpClient.Do(request)

	decodedResponse = map[string]interface{}{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)

	sessionToken := decodedResponse["session_token"].(string)
	// token := decodedResponse["token"].(string)

	// ========================== get session sse ==========================
	sessionSSEContext, _ := context.WithTimeout(context.Background(), time.Duration(300*time.Millisecond))
	var sessionSSEResponse *http.Response
	sessionSSEDoneChan := make(chan bool)

	go func() {
		request, _ := http.NewRequestWithContext(sessionSSEContext, http.MethodGet, suite.ServerBaseURL+"/sessions/sse", nil)
		q := request.URL.Query()
		q.Add("session_token", sessionToken)
		request.URL.RawQuery = q.Encode()
		sessionSSEResponse, _ = httpClient.Do(request)
		sessionSSEDoneChan <- true
	}()
	<-sessionSSEDoneChan

	sessionSSEResponseBytes, _ := ioutil.ReadAll(sessionSSEResponse.Body)
	sessionSSEResponseString := string(sessionSSEResponseBytes)
	replacer := strings.NewReplacer("\n", "", "data: ", "")
	sessionSSEResponseString = replacer.Replace(sessionSSEResponseString)

	decodedResponse = map[string]interface{}{}
	json.Unmarshal([]byte(sessionSSEResponseString), &decodedResponse)
	sessionID := decodedResponse["id"].(string)

	// ========================== scan session ==========================
	params = map[string]interface{}{
		"store_id": 1,
	}
	jsonParams, _ = json.Marshal(params)
	request, _ = http.NewRequest(http.MethodPut, suite.ServerBaseURL+"/sessions/"+sessionID, bytes.NewBuffer(jsonParams))
	request.Header.Add("Authorization", sessionID)
	request.Header.Set("Content-Type", "application/json")
	response, _ = httpClient.Do(request)

	decodedResponse = map[string]interface{}{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)

	// ========================== get store sse ==========================
	getStoreSSEContext, _ := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	var getStoreSSEResponse *http.Response
	getStoreSSEDoneChan := make(chan bool)

	go func() {
		request, _ := http.NewRequestWithContext(getStoreSSEContext, http.MethodGet, suite.ServerBaseURL+"/stores/1/sse", nil)
		getStoreSSEResponse, _ = httpClient.Do(request)
		getStoreSSEDoneChan <- true
	}()

	// ========================== create customers ==========================
	params = map[string]interface{}{
		"store_id": 1,
		"customers": []map[string]interface{}{
			{
				"name":     "kida",
				"phone":    "0932000000",
				"queue_id": 1,
			},
			{
				"name":     "kidb",
				"phone":    "0932000000",
				"queue_id": 1,
			},
			{
				"name":     "kidc",
				"phone":    "0932000000",
				"queue_id": 2,
			},
		},
	}

	time.Sleep(400 * time.Millisecond) // smaller than getStoreSSEContext timeout (400 < 500)
	jsonParams, _ = json.Marshal(params)
	request, _ = http.NewRequest(http.MethodPost, suite.ServerBaseURL+"/customers", bytes.NewBuffer(jsonParams))
	request.Header.Add("Authorization", sessionID)
	request.Header.Set("Content-Type", "application/json")
	response, _ = httpClient.Do(request)

	decodedResponseSlice := []map[string]interface{}{}
	json.NewDecoder(response.Body).Decode(&decodedResponseSlice)
	suite.Equal(3, len(decodedResponseSlice))

	<-getStoreSSEDoneChan
	re := regexp.MustCompile("data: ")
	getStoreSSEResponseBytes, _ := ioutil.ReadAll(getStoreSSEResponse.Body)
	getStoreSSEResponseString := string(getStoreSSEResponseBytes)
	matches := re.FindAllStringIndex(getStoreSSEResponseString, -1)
	suite.Equal(2, len(matches))

	// ========================== check store_sesisons in db ==========================
	query := `SELECT state 
				FROM store_sessions
				WHERE id=$1`
	rows := suite.db.QueryRow(query, sessionID)
	sessionState := ""
	err := rows.Scan(&sessionState)
	suite.NoError(err)
	suite.Equal("used", sessionState)

}

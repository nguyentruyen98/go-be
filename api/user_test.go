package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/nguyentruyen98/go-be/db/mock"
	db "github.com/nguyentruyen98/go-be/db/sqlc"
	"github.com/nguyentruyen98/go-be/util"
	"github.com/stretchr/testify/require"
)

func randomUser() (db.Users, string) {
	return db.Users{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}, "mypassword"

}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		body: gin.H{
			"username": user.Username,
			"email":    user.Email,
			"fullName": user.FullName,
			"password": password,
		},
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
		},
		checkResponse: func(recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchUser(t, recorder.Body, user)
		},
	}}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// *build stubs
			tc.buildStubs(store)

			// * star test server and send request

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)
			// * Send request and record its response in recorder
			server.router.ServeHTTP(recorder, request)

			// * Check response
			tc.checkResponse(recorder)
		})
	}

}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.Users) {
	data, err := ioutil.ReadAll(body)

	require.NoError(t, err)

	var getUser db.Users

	err = json.Unmarshal(data, &getUser)

	require.NoError(t, err)

	require.Equal(t, user, getUser)
}

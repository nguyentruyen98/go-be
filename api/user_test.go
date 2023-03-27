package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)

}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
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
			arg := db.CreateUserParams{
				Username: user.Username,
				Email:    user.Email,
				FullName: user.FullName,
			}
			store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
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

			server := newTestServer(t, store)
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

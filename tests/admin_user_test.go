package tests

import (
	"covid19kalteng/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGetUserList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	auth.GET("/admin/users").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test query found
	obj := auth.GET("/admin/users").WithQuery("username", "adminkey").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// test query invalid
	obj = auth.GET("/admin/users").WithQuery("username", "should not found this").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 0)
}

func TestNewUser(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"username": "test user",
		"email":    "testuser@covid19kalteng.id",
		"phone":    "08111",
		"status":   "active",
		"roles":    []int{3},
	}
	// normal scenario reporter user
	obj := auth.POST("/admin/users").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("username").ValueEqual("username", "test user")

	// normal scenario administrator
	payload = map[string]interface{}{
		"username": "test user1",
		"email":    "testuser1@covid19kalteng.id",
		"phone":    "08112",
		"status":   "active",
		"roles":    []int{1},
	}
	obj = auth.POST("/admin/users").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("username").ValueEqual("username", "test user1")

	// unique test
	auth.POST("/admin/users").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test invalid
	payload = map[string]interface{}{
		"username": "",
	}
	auth.POST("/admin/users").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

func TestGetUserbyID(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	// valid response
	obj := auth.GET("/admin/users/1").Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/users/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestPatchUser(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	adminToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+adminToken)
	})

	payload := map[string]interface{}{
		"status": "inactive",
	}

	// valid response
	obj := auth.PATCH("/admin/users/1").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("status").ValueEqual("status", "inactive")

	// test change roles of reporter to core system
	payload = map[string]interface{}{
		"roles": []int{1},
	}
	auth.PATCH("/admin/users/3").WithJSON(payload).
		Expect().Status(http.StatusUnprocessableEntity).JSON().Object()

	// uniques
	auth.PATCH("/admin/users/2").WithJSON(map[string]interface{}{
		"username": "adminkey",
	}).Expect().Status(http.StatusInternalServerError).JSON().Object()

	auth.PATCH("/admin/users/2").WithJSON(map[string]interface{}{
		"email": "admin@covid19kalteng.com",
	}).Expect().Status(http.StatusInternalServerError).JSON().Object()

	auth.PATCH("/admin/users/2").WithJSON(map[string]interface{}{
		"phone": "081234567890",
	}).Expect().Status(http.StatusInternalServerError).JSON().Object()

	// test invalid token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrong token")
	})
	auth.PATCH("/admin/users/1").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	authController "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/controllers/authController"
	orgController "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/controllers/orgController"
	"github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/middleware/isAuth"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
)

type Org struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OrgId    string `json:"org_id"`
}

type UserDelete struct {
	Username string `json:"username"`
}

type OrgDelete struct {
	Name string `json:"name"`
}

func Test_LoginUserWithNoBody(t *testing.T) {
	userServer := fiber.New()
	// ! for now we will use default admin user that we created in db
	userServer.Post("/api/auth/login", authController.Login)

	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", nil)
	userReq.Header.Set("Content-Type", "application/json")

	userRes, _ := userServer.Test(userReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&adminBody)

	assert.Equal(t, http.StatusBadRequest, userRes.StatusCode)
	assert.Equal(t, "fail", adminBody["status"])
}

func Test_LoginUserWithValidBody(t *testing.T) {
	userServer := fiber.New()
	// ! for now we will use default admin user that we created in db
	userServer.Post("/api/auth/login", authController.Login)

	userInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	userReqBody, _ := json.Marshal(userInput)
	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq.Header.Set("Content-Type", "application/json")

	userRes, _ := userServer.Test(userReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&adminBody)

	assert.Equal(t, http.StatusOK, userRes.StatusCode)
	assert.Equal(t, "success", adminBody["status"])

	// ! check if cookie is set with token
	cookie := userRes.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie)
}

func Test_CreateOrgWithoutAdminLogin(t *testing.T) {
	server := fiber.New()
	server.Use(isAuth.IsAdminCheck)
	server.Post("/api/org", orgController.CreateOrg)

	input := Org{
		Name: "Org-2",
		Head: "Rohit-1",
	}

	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/org", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res, _ := server.Test(req, -1)

	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "fail", body["status"])
}

func Test_CreateOrgWithAdminLogin(t *testing.T) {
	userServer := fiber.New()
	// ! first login as admin to get token
	userServer.Post("/api/auth/login", authController.Login)

	userInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	userReqBody, _ := json.Marshal(userInput)
	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq.Header.Set("Content-Type", "application/json")

	userRes, _ := userServer.Test(userReq, -1)
	var userBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&userBody)

	assert.Equal(t, http.StatusOK, userRes.StatusCode)
	assert.Equal(t, "success", userBody["status"])

	// ! check if cookie is set with token
	cookie := userRes.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	// ! now create org with token
	userServer.Use(isAuth.IsAdminCheck)
	userServer.Post("/api/org", orgController.CreateOrg)

	orgInput := Org{
		Name: "Org-2",
		Head: "Rohit-1",
	}

	orgReqBody, _ := json.Marshal(orgInput)
	orgReq, _ := http.NewRequest(http.MethodPost, "/api/org", bytes.NewBuffer(orgReqBody))
	orgReq.Header.Set("Content-Type", "application/json")
	orgReq.Header.Set("Cookie", cookie)

	orgRes, _ := userServer.Test(orgReq, -1)
	var orgBody map[string]interface{}
	json.NewDecoder(orgRes.Body).Decode(&orgBody)

	assert.Equal(t, http.StatusCreated, orgRes.StatusCode)
	assert.Equal(t, "success", orgBody["status"])
}

func Test_CreateOrgWithNoBody(t *testing.T) {
	userServer := fiber.New()
	// ! first login as admin to get token
	userServer.Post("/api/auth/login", authController.Login)

	userInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	userReqBody, _ := json.Marshal(userInput)
	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq.Header.Set("Content-Type", "application/json")

	userRes, _ := userServer.Test(userReq, -1)
	var userBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&userBody)

	assert.Equal(t, http.StatusOK, userRes.StatusCode)
	assert.Equal(t, "success", userBody["status"])

	// ! check if cookie is set with token
	cookie := userRes.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	// ! now create org with token
	userServer.Use(isAuth.IsAdminCheck)
	userServer.Post("/api/org", orgController.CreateOrg)

	orgReq, _ := http.NewRequest(http.MethodPost, "/api/org", nil)
	orgReq.Header.Set("Content-Type", "application/json")
	orgReq.Header.Set("Cookie", cookie)

	orgRes, _ := userServer.Test(orgReq, -1)
	var orgBody map[string]interface{}
	json.NewDecoder(orgRes.Body).Decode(&orgBody)

	assert.Equal(t, http.StatusBadRequest, orgRes.StatusCode)
	assert.Equal(t, "fail", orgBody["status"])
}

func Test_CreateUserWithoutAdminLogin(t *testing.T) {
	adminServer := fiber.New()
	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Post("/api/auth/signup", authController.SignUp)

	userInput := UserSignUp{
		Username: "Rohit-2",
		Password: "pass-1",
		OrgId:    "1", // ! for now we'll use orgId 1 - which is created by admin and will always be present
	}
	userReqBody, _ := json.Marshal(userInput)
	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(userReqBody))
	userReq.Header.Set("Content-Type", "application/json")

	userRes, _ := adminServer.Test(userReq, -1)
	var userBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&userBody)

	assert.Equal(t, http.StatusUnauthorized, userRes.StatusCode)
	assert.Equal(t, "fail", userBody["status"])
}

func Test_CreateUserWithAdminLogin(t *testing.T) {
	adminServer := fiber.New()
	// ! first login as admin to get token
	adminServer.Post("/api/auth/login", authController.Login)

	adminInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	adminReqBody, _ := json.Marshal(adminInput)
	adminReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(adminReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)

	assert.Equal(t, http.StatusOK, adminRes.StatusCode)
	assert.Equal(t, "success", adminBody["status"])

	// ! check if cookie is set with token
	adminCookie := adminRes.Header.Get("Set-Cookie")
	assert.NotEmpty(t, adminCookie)

	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Post("/api/auth", authController.SignUp)

	userInput := UserSignUp{
		Username: "Rohit-2",
		Password: "pass-1",
		OrgId:    "1", // ! for now we'll use orgId 1 - which is created by admin and will always be present
	}

	userReqBody, _ := json.Marshal(userInput)
	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(userReqBody))
	userReq.Header.Set("Content-Type", "application/json")
	userReq.Header.Set("Cookie", adminCookie)

	userRes, _ := adminServer.Test(userReq, -1)
	var userBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&userBody)

	assert.Equal(t, http.StatusCreated, userRes.StatusCode)
	assert.Equal(t, "success", userBody["status"])
}

func Test_CreateUserWithNoBody(t *testing.T) {
	adminServer := fiber.New()
	// ! first login as admin to get token
	adminServer.Post("/api/auth/login", authController.Login)

	adminInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	adminReqBody, _ := json.Marshal(adminInput)
	adminReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(adminReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)

	assert.Equal(t, http.StatusOK, adminRes.StatusCode)
	assert.Equal(t, "success", adminBody["status"])

	// ! check if cookie is set with token
	adminCookie := adminRes.Header.Get("Set-Cookie")
	assert.NotEmpty(t, adminCookie)

	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Post("/api/auth", authController.SignUp)

	userReq, _ := http.NewRequest(http.MethodPost, "/api/auth", nil)
	userReq.Header.Set("Content-Type", "application/json")
	userReq.Header.Set("Cookie", adminCookie)

	userRes, _ := adminServer.Test(userReq, -1)
	var userBody map[string]interface{}
	json.NewDecoder(userRes.Body).Decode(&userBody)

	assert.Equal(t, http.StatusBadRequest, userRes.StatusCode)
	assert.Equal(t, "fail", userBody["status"])
}

func Test_GetUsersInfoWithoutToken(t *testing.T) {
	userServer := fiber.New()
	userServer.Post("/api/auth/user", authController.GetUser)

	userReq1, _ := http.NewRequest(http.MethodPost, "/api/auth/user", nil)
	userReq1.Header.Set("Content-Type", "application/json")

	userRes1, _ := userServer.Test(userReq1, -1)
	var userBody1 map[string]interface{}
	json.NewDecoder(userRes1.Body).Decode(&userBody1)

	// ! should fail since no or bad token is provided
	assert.Equal(t, http.StatusUnauthorized, userRes1.StatusCode)
	assert.Equal(t, "fail", userBody1["status"])
}

func Test_GetUserInfoWithToken(t *testing.T) {
	userServer := fiber.New()
	userInput := UserLogin{
		Username: "Rohit-2",
		Password: "pass-1",
	}
	userReqBody, _ := json.Marshal(userInput)

	userServer.Post("/api/auth/login", authController.Login)
	userReq2, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq2.Header.Set("Content-Type", "application/json")

	userRes2, _ := userServer.Test(userReq2, -1)
	var userBody2 map[string]interface{}
	json.NewDecoder(userRes2.Body).Decode(&userBody2)

	assert.Equal(t, http.StatusOK, userRes2.StatusCode)
	assert.Equal(t, "success", userBody2["status"])

	// ! check if cookie is set with token
	cookie := userRes2.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	// ! now get user with token
	userServer.Post("/api/auth/user", authController.GetUser)
	userReq3, _ := http.NewRequest(http.MethodPost, "/api/auth/user", nil)
	userReq3.Header.Set("Content-Type", "application/json")
	userReq3.Header.Set("Cookie", cookie)

	userRes3, _ := userServer.Test(userReq3, -1)
	var userBody3 map[string]interface{}
	json.NewDecoder(userRes3.Body).Decode(&userBody3)

	assert.Equal(t, http.StatusOK, userRes3.StatusCode)
	assert.Equal(t, "success", userBody3["status"])
	assert.Equal(t, "Rohit-2", userBody3["data"].(map[string]interface{})["username"])
	assert.Equal(t, float64(1), (userBody3["data"].(map[string]interface{})["org_id"]))
	assert.Equal(t, false, userBody3["data"].(map[string]interface{})["is_admin"])
}

func Test_LogoutUser(t *testing.T) {
	userServer := fiber.New()
	userInput := UserLogin{
		Username: "Rohit-2",
		Password: "pass-1",
	}
	userReqBody, _ := json.Marshal(userInput)

	userServer.Post("/api/auth/login", authController.Login)
	userReq2, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq2.Header.Set("Content-Type", "application/json")

	userRes2, _ := userServer.Test(userReq2, -1)
	var userBody2 map[string]interface{}
	json.NewDecoder(userRes2.Body).Decode(&userBody2)
	cookie := userRes2.Header.Get("Set-Cookie")

	// ! now logout user
	userServer.Post("/api/auth/logout", authController.Logout)
	userReq3, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	userReq3.Header.Set("Content-Type", "application/json")
	userReq3.Header.Set("Cookie", cookie)

	userRes3, _ := userServer.Test(userReq3, -1)
	var userBody3 map[string]interface{}
	json.NewDecoder(userRes3.Body).Decode(&userBody3)

	assert.Equal(t, http.StatusOK, userRes3.StatusCode)
	assert.Equal(t, "success", userBody3["status"])
}

func Test_RefreshTokenAfterLogin(t *testing.T) {
	userServer := fiber.New()
	userInput := UserLogin{
		Username: "Rohit-2",
		Password: "pass-1",
	}
	userReqBody, _ := json.Marshal(userInput)

	userServer.Post("/api/auth/login", authController.Login)
	userReq2, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq2.Header.Set("Content-Type", "application/json")

	userRes2, _ := userServer.Test(userReq2, -1)
	var userBody2 map[string]interface{}
	json.NewDecoder(userRes2.Body).Decode(&userBody2)
	cookie := userRes2.Header.Get("Set-Cookie")

	// ! now refresh token
	userServer.Post("/api/auth/refresh", authController.RefreshToken)
	userReq3, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	userReq3.Header.Set("Content-Type", "application/json")
	userReq3.Header.Set("Cookie", cookie)

	userRes3, _ := userServer.Test(userReq3, -1)
	var userBody3 map[string]interface{}
	json.NewDecoder(userRes3.Body).Decode(&userBody3)

	assert.Equal(t, http.StatusOK, userRes3.StatusCode)
	assert.Equal(t, "success", userBody3["status"])
}

func Test_GetOrgPeopleWithoutToken(t *testing.T) {
	userServer := fiber.New()
	userServer.Get("/api/org", orgController.GetOrgPeople)

	userReq1, _ := http.NewRequest(http.MethodGet, "/api/org", nil)
	userReq1.Header.Set("Content-Type", "application/json")

	userRes1, _ := userServer.Test(userReq1, -1)
	var userBody1 map[string]interface{}
	json.NewDecoder(userRes1.Body).Decode(&userBody1)

	// ! should fail since no or bad token is provided
	assert.Equal(t, http.StatusUnauthorized, userRes1.StatusCode)
	assert.Equal(t, "fail", userBody1["status"])
}

func Test_GetOrgPeopleWithToken(t *testing.T) {
	userServer := fiber.New()
	userInput := UserLogin{
		Username: "Rohit-2",
		Password: "pass-1",
	}
	userReqBody, _ := json.Marshal(userInput)

	userServer.Post("/api/auth/login", authController.Login)
	userReq2, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(userReqBody))
	userReq2.Header.Set("Content-Type", "application/json")

	userRes2, _ := userServer.Test(userReq2, -1)
	var userBody2 map[string]interface{}
	json.NewDecoder(userRes2.Body).Decode(&userBody2)
	cookie := userRes2.Header.Get("Set-Cookie")

	// ! now get org people with token
	userServer.Get("/api/org", orgController.GetOrgPeople)
	userReq3, _ := http.NewRequest(http.MethodGet, "/api/org", nil)
	userReq3.Header.Set("Content-Type", "application/json")
	userReq3.Header.Set("Cookie", cookie)

	userRes3, _ := userServer.Test(userReq3, -1)
	var userBody3 map[string]interface{}
	json.NewDecoder(userRes3.Body).Decode(&userBody3)

	assert.Equal(t, http.StatusOK, userRes3.StatusCode)
	assert.Equal(t, "success", userBody3["status"])

	// ! check if the data is correct
	assert.Equal(t, "Rohit-1", userBody3["data"].([]interface{})[0].(map[string]interface{})["username"])
	assert.Equal(t, true, userBody3["data"].([]interface{})[0].(map[string]interface{})["is_admin"])
	assert.Equal(t, "Rohit-2", userBody3["data"].([]interface{})[1].(map[string]interface{})["username"])
	assert.Equal(t, false, userBody3["data"].([]interface{})[1].(map[string]interface{})["is_admin"])
}

func Test_DeleteUserWithoutAdminToken(t *testing.T) {
	adminServer := fiber.New()
	adminServer.Post("/api/auth/login", authController.Login)

	adminInput := UserLogin{
		Username: "Rohit-2",
		Password: "pass-1",
	}

	adminReqBody, _ := json.Marshal(adminInput)
	adminReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(adminReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)
	cookie := adminRes.Header.Get("Set-Cookie")

	// ! try to delete but fails - as Rohit-2 is not admin
	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Delete("/api/auth", authController.DeleteUser)

	deleteUserInput := UserDelete{
		Username: "Rohit-2",
	}
	deleteUserReqBody, _ := json.Marshal(deleteUserInput)

	adminReq2, _ := http.NewRequest(http.MethodDelete, "/api/auth", bytes.NewBuffer(deleteUserReqBody))
	adminReq2.Header.Set("Content-Type", "application/json")
	adminReq2.Header.Set("Cookie", cookie)

	adminRes2, _ := adminServer.Test(adminReq2, -1)
	var adminBody2 map[string]interface{}
	json.NewDecoder(adminRes2.Body).Decode(&adminBody2)

	assert.Equal(t, http.StatusUnauthorized, adminRes2.StatusCode)
	assert.Equal(t, "fail", adminBody2["status"])
}

func Test_DeleteUserWithAdminToken(t *testing.T) {
	adminServer := fiber.New()
	adminServer.Post("/api/auth/login", authController.Login)

	adminInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	adminReqBody, _ := json.Marshal(adminInput)
	adminReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(adminReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)
	cookie := adminRes.Header.Get("Set-Cookie")

	// ! now delete user
	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Delete("/api/auth", authController.DeleteUser)

	deleteUserInput := UserDelete{
		Username: "Rohit-2",
	}
	deleteUserReqBody, _ := json.Marshal(deleteUserInput)

	adminReq2, _ := http.NewRequest(http.MethodDelete, "/api/auth", bytes.NewBuffer(deleteUserReqBody))
	adminReq2.Header.Set("Content-Type", "application/json")
	adminReq2.Header.Set("Cookie", cookie)

	adminRes2, _ := adminServer.Test(adminReq2, -1)
	var adminBody2 map[string]interface{}
	json.NewDecoder(adminRes2.Body).Decode(&adminBody2)

	assert.Equal(t, http.StatusOK, adminRes2.StatusCode)
	assert.Equal(t, "success", adminBody2["status"])
}

func Test_DeleteOrgWithoutAdminToken(t *testing.T) {
	adminServer := fiber.New()
	// ! try to delete but fails - as no user is logged in
	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Delete("/api/auth", orgController.DeleteOrg)

	deleteOrgInput := OrgDelete{
		Name: "Org-2",
	}
	deleteOrgReqBody, _ := json.Marshal(deleteOrgInput)

	adminReq, _ := http.NewRequest(http.MethodDelete, "/api/auth", bytes.NewBuffer(deleteOrgReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)

	assert.Equal(t, http.StatusUnauthorized, adminRes.StatusCode)
	assert.Equal(t, "fail", adminBody["status"])
}

func Test_DeleteOrgWithAdminToken(t *testing.T) {
	adminServer := fiber.New()
	adminServer.Post("/api/auth/login", authController.Login)

	adminInput := UserLogin{
		Username: "Rohit-1",
		Password: "pass-1",
	}

	adminReqBody, _ := json.Marshal(adminInput)
	adminReq, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(adminReqBody))
	adminReq.Header.Set("Content-Type", "application/json")

	adminRes, _ := adminServer.Test(adminReq, -1)
	var adminBody map[string]interface{}
	json.NewDecoder(adminRes.Body).Decode(&adminBody)
	cookie := adminRes.Header.Get("Set-Cookie")

	// ! now delete org
	adminServer.Use(isAuth.IsAdminCheck)
	adminServer.Delete("/api/org", orgController.DeleteOrg)

	deleteOrgInput := OrgDelete{
		Name: "Org-2",
	}
	deleteOrgReqBody, _ := json.Marshal(deleteOrgInput)

	adminReq2, _ := http.NewRequest(http.MethodDelete, "/api/org", bytes.NewBuffer(deleteOrgReqBody))
	adminReq2.Header.Set("Content-Type", "application/json")
	adminReq2.Header.Set("Cookie", cookie)

	adminRes2, _ := adminServer.Test(adminReq2, -1)
	var adminBody2 map[string]interface{}
	json.NewDecoder(adminRes2.Body).Decode(&adminBody2)

	assert.Equal(t, http.StatusOK, adminRes2.StatusCode)
	assert.Equal(t, "success", adminBody2["status"])
}

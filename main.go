package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
)

const (
	JwtSignKey                 = "BZ0niKtToA4TwoNjP1na"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpiryDuration  = time.Hour * 24
	RefreshTokenExpiryDuration = time.Hour * 24 * 7
)

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error":"invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpiryDuration, RefreshTokenExpiryDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	_, err = userSvc.Register(uReq)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	writer.Write([]byte(`"message":"user created successfully!"`))
}

func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error":"invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	var loginRequest userservice.LoginRequest
	err = json.Unmarshal(data, &loginRequest)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpiryDuration, RefreshTokenExpiryDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	resp, loginErr := userSvc.Login(loginRequest)
	if loginErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, loginErr.Error())

		return
	}

	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, marshalErr.Error())

		return
	}

	writer.Write(data)
}

func userProfileHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(writer, `{"error":"invalid method"}`)

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpiryDuration, RefreshTokenExpiryDuration)

	auth := req.Header.Get("Authorization")
	tokenClaims, err := authSvc.ParseToken(auth)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"token is not valid"}`)

		return
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	resp, profileErr := userSvc.Profile(userservice.ProfileRequest{UserID: tokenClaims.UserID})
	if profileErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, profileErr.Error())

		return
	}

	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, marshalErr.Error())

		return
	}

	writer.Write(data)
}

func main() {
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/users/login", userLoginHandler)
	http.HandleFunc("/users/profile", userProfileHandler)

	fmt.Println("server is listening on port 8080...")

	http.ListenAndServe(":8080", nil)

}

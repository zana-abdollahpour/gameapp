package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gameapp/repository/mysql"
	"gameapp/service/userservice"
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

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

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

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

	_, err = userSvc.Login(loginRequest)
	if err != nil {
		fmt.Fprintf(writer, `{"error":"%s"}`, err.Error())

		return
	}

	writer.Write([]byte(`"message":"user credentials are ok!"`))
}

func main() {
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/users/login", userLoginHandler)

	fmt.Println("server is listening on port 8080...")

	http.ListenAndServe(":8080", nil)

}

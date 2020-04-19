package spec_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DaggerJackfast/gst/e2e/test_utils"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
)

var _ = Describe("Auth", func() {
	var url string
	contentType := "application/json"

	BeforeEach(func() {
		err := Loader.Load()
		if err != nil {
			log.Fatalf("Not load test fixtures: %s", err)
		}

	})

	AfterEach(func() {
		tables := []string{"sessions", "user_profile_tokens", "users"}
		Cleaner.Clean(tables...)
	})

	Context("User registration", func() {
		BeforeEach(func() {
			url = fmt.Sprintf("%s/auth/register", Server.URL)
		})
		It("User can register", func() {
			user := map[string]string{
				"username": "test_user",
				"email":    "test@test.com",
				"password": "test_user_password",
			}
			data, err := json.Marshal(user)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("User register failed with json decode error", func() {
			response, err := Client.Post(url, contentType, bytes.NewBuffer([]byte("")))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "EOF",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User register failed with unique email error", func() {
			user := map[string]string{
				"username": "first_user",
				"email":    "first_user@test.test",
				"password": "first_user_password",
			}
			data, err := json.Marshal(user)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}

			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": `pq: duplicate key value violates unique constraint "users_email_key"`,
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
	})

	Context("User login", func() {
		BeforeEach(func() {
			url = fmt.Sprintf("%s/auth/login", Server.URL)
		})
		It("User login successfully", func() {
			login := map[string]string{
				"email":       "first_user@test.test",
				"password":    "first_user_password",
				"fingerprint": "testfingerprint",
			}
			data, err := json.Marshal(login)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}

			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("User login failed with json decode error", func() {
			response, err := Client.Post(url, contentType, bytes.NewBuffer([]byte("")))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "EOF",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User login failed with json validate error", func() {
			login := map[string]string{
				"email":       "first_user.test.test",
				"password":    "",
				"fingerprint": "testfingerprint",
			}
			data, err := json.Marshal(login)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "Key: 'EmailPasswordFingerprint.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User login failed with json validate error", func() {
			login := map[string]string{
				"email":       "first_user@test.test",
				"password":    "wrong password",
				"fingerprint": "testfingerprint",
			}
			data, err := json.Marshal(login)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusForbidden))
			expectedBody := map[string]string{
				"error": "Password is incorrect",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
	})

	Context("User forgot password", func() {
		BeforeEach(func(){
			url = fmt.Sprintf("%s/auth/forgot-password", Server.URL)
		})
		It("User recovery password successfully", func() {
			userEmail := map[string]string{
				"email": "first_user@test.test",
			}
			data, err := json.Marshal(userEmail)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("User recovery password failed with json decode error", func(){
			response, err := Client.Post(url, contentType, bytes.NewBuffer([]byte("")))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "EOF",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User login failed with json validate error", func() {
			userEmail := map[string]string{
				"email": "first_user.test.test",
			}
			data, err := json.Marshal(userEmail)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "Key: 'UserEmail.Email' Error:Field validation for 'Email' failed on the 'email' tag",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User login failed with email not found", func() {
			userEmail := map[string]string{
				"email": "unknown@test.test",
			}
			data, err := json.Marshal(userEmail)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "sql: no rows in result set",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})

	})
	Context("User reset password", func(){
		BeforeEach(func(){
			url = fmt.Sprintf("%s/auth/reset-password", Server.URL)
		})
		It("User reset password successfully", func(){
			emailPassword := map[string]string{
				"email": "first_user@test.test",
				"password": "first_user_password",
				"token": "fc5c078c673b9799e9a9b6a95531d77c",
			}
			data, err := json.Marshal(emailPassword)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			expectedBody := map[string]string{
				"status": "success",
				"message": "Your password successfully changed.",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User reset password failed with json decode error", func(){
			response, err := Client.Post(url, contentType, bytes.NewBuffer([]byte("")))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "EOF",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User reset failed with json validate error", func() {
			emailPassword := map[string]string{
				"email": "first_user.test.test",
				"password": "",
				"token": "",
			}
			data, err := json.Marshal(emailPassword)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "Key: 'EmailPasswordToken.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'EmailPasswordToken.Token' Error:Field validation for 'Token' failed on the 'required' tag",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User reset failed with email not found error", func() {
			emailPassword := map[string]string{
				"email": "wrong@test.test",
				"password": "first_user_password",
				"token": "fc5c078c673b9799e9a9b6a95531d77c",
			}
			data, err := json.Marshal(emailPassword)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "sql: no rows in result set",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
		It("User reset failed with token validate error", func() {
			emailPassword := map[string]string{
				"email": "first_user@test.test",
				"password": "first_user_password",
				"token": "wrong_token",
			}
			data, err := json.Marshal(emailPassword)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			expectedBody := map[string]string{
				"error": "The token is wrong.",
			}
			Expect(test_utils.GetResponseBodyJson(response)).To(Equal(expectedBody))
		})
	})

})

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
			url := fmt.Sprintf("%s/auth/register", Server.URL)
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("User can't register with exists email", func() {
			user := map[string]string{
				"username": "first_user",
				"email":    "first_user@test.test",
				"password": "first_user_password",
			}
			data, err := json.Marshal(user)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			url := fmt.Sprintf("%s/auth/register", Server.URL)
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
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
		It("User login failed with json validate error", func(){
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
		It("User login failed with json validate error", func(){
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
		It("User recovery password successfully", func() {
			userEmail := map[string]string{
				"email": "first_user@test.test",
			}
			data, err := json.Marshal(userEmail)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			url := fmt.Sprintf("%s/auth/forgot-password", Server.URL)
			response, err := Client.Post(url, contentType, bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})

})

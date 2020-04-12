package spec_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
)

var _ = Describe("Auth", func() {
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
	Context("User registration", func(){
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
			response, err := Client.Post(url, "application/json", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("User can't register with exists email", func(){
			user := map[string]string{
				"username": "first_user",
				"email": "first_user@test.test",
				"password": "first_user_password",
			}
			data, err := json.Marshal(user)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			url := fmt.Sprintf("%s/auth/register", Server.URL)
			response, err := Client.Post(url, "application/json", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		})
	})
	Context("User login", func(){
		It("User login successfully", func(){
			login:=map[string]string{
				"email":"first_user@test.test",
				"password": "first_user_password",
				"fingerprint": "testfingerprint",
			}
			data, err := json.Marshal(login)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			url := fmt.Sprintf("%s/auth/login", Server.URL)
			response, err := Client.Post(url, "application/json", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})
	Context("User forgot password", func(){
		It("User recovery password successfully", func(){
			userEmail := map[string]string{
				"email": "first_user@test.test",
			}
			data, err := json.Marshal(userEmail)
			if err != nil {
				log.Printf("json marshal error: %s", err)
			}
			url := fmt.Sprintf("%s/auth/forgot-password", Server.URL)
			response, err := Client.Post(url, "application/json", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Println("response data", bodyString)
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})

})

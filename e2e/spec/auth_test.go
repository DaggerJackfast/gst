package spec_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
)

var _ = Describe("Auth", func() {
	AfterEach(func() {
		tables := []string{"sessions", "user_profile_tokens", "users"}
		Cleaner.Clean(tables...)
	})
	It("User can register", func() {
		user := map[string]string{
			"username": "test_user",
			"email":    "test@test.com",
			"password": "qweqweqwe",
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
})

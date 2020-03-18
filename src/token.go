package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateTokenPair(userId uint64) (map[string]string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Second * time.Duration(ExpiredInAccessToken)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["exp"] = time.Now().Add(time.Second * time.Duration(ExpiredInRefreshToken)).Unix()
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected singing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	authHeader := strings.Split(bearerToken, " ")
	lengthOfAuthHeader := 2
	if len(authHeader) == lengthOfAuthHeader && authHeader[0] == "Bearer" {
		return authHeader[1]
	}
	return ""
}

func ExtractTokenId(r *http.Request) (uint64, error) {
	tokenString := ExtractToken(r)
	userId, err := ExtractId(tokenString)
	if err != nil {
		return 0, nil
	}
	return userId, nil
}

func ExtractId(tokenString string)(uint64, error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected singing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint64(id), nil
	}
	return 0, nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}

func GenerateToken(n int) (string, error){
	b:=make([]byte, n)
	if _, err := rand.Read(b); err != nil{
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetIp(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != ""{
		return forwarded
	}
	return r.RemoteAddr
}

func GetUserAgent(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	return userAgent
}

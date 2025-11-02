package api

import (
	"context"
	"errors"
	"log"
	"main/internal/evn"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClientRequest struct {
	Ip                string
	RequestSentPerMin int64
	LastSeenAt        time.Time
}

func (api *ApiService) HandleRateLimiter(h http.Handler) http.Handler {

	var clientRequests = make(map[string]*ClientRequest)
	var mutex sync.Mutex

	//remove ip every 5 minutes if LastSeenAt > 5 munite
	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for ip, client := range clientRequests {

				difference := client.LastSeenAt.Sub(time.Now())

				if difference > 5*time.Minute {
					delete(clientRequests, ip)
				}
				if time.Since(client.LastSeenAt) > 5*time.Minute {

				}
			}
			mutex.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mutex.Lock()
		defer mutex.Unlock()

		ip := r.RemoteAddr
		client, ok := clientRequests[ip]
		maxRequest := api.config.RateLimitConfig.MaxRequestPerMin

		if ok {
			if client.RequestSentPerMin > 0 {
				//client has not reached limit
				clientRequests[ip] = &ClientRequest{
					Ip:                ip,
					RequestSentPerMin: client.RequestSentPerMin - 1, //remove to request pull
					LastSeenAt:        time.Now(),
				}
				log.Println("curent request: " + (strconv.Itoa(int(client.RequestSentPerMin - 1))))
				h.ServeHTTP(w, r)
			} else {
				clientRequests[ip] = &ClientRequest{
					LastSeenAt: time.Now(),
				}
				tooManyRequest(w, r, errors.New("to many request limit reached"))
				return
			}
		} else {
			clientRequests[ip] = &ClientRequest{
				Ip:                ip,
				RequestSentPerMin: maxRequest - 1, // -1 is to count this current request
				LastSeenAt:        time.Now(),
			}
			h.ServeHTTP(w, r)
		}

		log.Default().Println("Limiter", "method", r.Method, "path", r.URL.Path, "ip: "+ip+"requests "+strconv.Itoa(int(0)))

	})
}

func HandleJWTAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeaderString := r.Header.Get("Authorization")

		if authHeaderString == "" {
			err := errors.New("authorization header required")
			unauthorized(w, r, err)
			return
		}

		tokenString := authHeaderString[7:]

		var secret_words string = "A request for a long text message: Search results showIf this is your intent, please clarify the context and what you want the text to be about."

		secret := evn.GetString(secret_words, "JWT_SERECT")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil {

			err := errors.New("invalid token")
			unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		exp := claims["exp"]
		username := claims["username"]

		float_exp := exp.(float64)

		t := time.Unix(int64(float_exp), 0)
		date_now := time.Now()

		if t.Before(date_now) {
			err := errors.New("token expired")
			unauthorized(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", username)
		h.ServeHTTP(w, r.WithContext(ctx))
	})

}

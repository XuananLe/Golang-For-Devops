package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type UserRequest struct{
	count int
	lastReset time.Time
}


type RateLimiter struct {
	threshold int
	mutex sync.Mutex
	duration time.Duration
	requests map[string]*UserRequest
}

var rl = NewRateLimiter(5, 5 * time.Second);


func NewRateLimiter(thresHold int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		threshold: thresHold,
		mutex: sync.Mutex{},
		duration: duration,
		requests: make(map[string]*UserRequest, 1000),
	}
}

func getIp(remoteAddr string) string   {
	hostIP, _, _ := net.SplitHostPort(remoteAddr)
	return hostIP
}


func (r *RateLimiter) isAllowed(ipAddr string) bool {
	r.mutex.Lock();
	defer r.mutex.Unlock();
	now := time.Now()
	
	if req, exists := r.requests[ipAddr]; exists {
		if now.Sub(req.lastReset) >= r.duration {
			req.count = 1;
			req.lastReset = now;
			return true;
		} else if req.count >= r.threshold {
			req.lastReset = now;
			return false
		}
		req.count++
        return true
	}
	r.requests[ipAddr] = &UserRequest{
        count:     1,
        lastReset: now,
    }
	return true;
}

func LoggingMap(data map[interface{}]int) {
	formatted := "{\n"
	for key, value := range data {
		formatted += fmt.Sprintf("  %v: %d\n", key, value)
	}
	formatted += "}"

	log.Println("Ratelimit:", formatted)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RateLimitMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.isAllowed(getIp(r.RemoteAddr)) {
			http.Error(w, "Too Many Requests (Rate Limit Exceeded)", http.StatusTooManyRequests)
			log.Printf("Blocked %s due to rate limit", r.RemoteAddr)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World! " + r.RemoteAddr)
}

func (r *RateLimiter) ResetRateLimit() {
	for {
		time.Sleep(r.duration)
		now := time.Now()
		r.mutex.Lock()
		for ip, req := range r.requests {
			if now.Sub(req.lastReset) >= r.duration {
				delete(r.requests, ip) // Chỉ xóa IP đã hết hạn
			}
		}
		r.mutex.Unlock()
		log.Println("Ratelimit map cleaned up expired entries.")
	}
}

func main() {
	go rl.ResetRateLimit();

	mux := http.NewServeMux()

	mux.Handle("/", LoggingMiddleware(RateLimitMiddleWare(http.HandlerFunc(HelloHandler))))

	port := ":8080"
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

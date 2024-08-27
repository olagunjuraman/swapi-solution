package utils

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func Respond(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func GetIp(request *http.Request) string {
	var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

	var ip string

	if xff := request.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")

		if i == -1 {
			i = len(xff)
		}

		ip = xff[:i]

	} else if xrip := request.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}

	if ip == "" {
		ipAddress, _, err := net.SplitHostPort(request.RemoteAddr)

		if err != nil {
			return ""
		}

		return strings.TrimSpace(ipAddress)
	}

	return ip
}

func ResourceId(url string) (int, error) {
	if !strings.Contains(url, "http") {
		return 0, errors.New("invalid URL")
	}

	s := strings.Split(url, "/")

	id, err := strconv.Atoi(s[len(s)-2])
	if err != nil {
		return 0, err
	}

	return id, err
}

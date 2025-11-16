package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <subscription_url> [hwid]")
		os.Exit(1)
	}

	subURL := os.Args[1]

	var hwid string
	if len(os.Args) >= 3 {
		hwid = os.Args[2]
	} else {
		hwid = generateRandomHWID()
	}

	fmt.Println("Subscription URL:", subURL)
	fmt.Println("Using x-hwid:", hwid)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", subURL, nil)
	if err != nil {
		fmt.Println("create request error:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "v2rayN/6.34 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("x-hwid", hwid)
	req.Header.Set("x-device-os", "Windows")
	req.Header.Set("x-ver-os", "10")
	req.Header.Set("x-device-model", "PC-Desktop")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http request error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("server returned:", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		os.Exit(1)
	}

	raw := strings.TrimSpace(string(body))
	if raw == "" {
		fmt.Println("empty response")
		return
	}

	decoded, ok := tryDecodeBase64(raw)
	if !ok {
		fmt.Println("Response is not valid base64, raw content:")
		fmt.Println(raw)
		return
	}

	fmt.Println("\nDecoded subscription:")
	fmt.Println(decoded)
}

func tryDecodeBase64(s string) (string, bool) {
	s = strings.TrimSpace(s)

	if out, err := decodeWithPadding(s, base64.StdEncoding); err == nil {
		return out, true
	}

	if out, err := decodeWithPadding(s, base64.URLEncoding); err == nil {
		return out, true
	}

	return "", false
}

func decodeWithPadding(s string, enc *base64.Encoding) (string, error) {
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	b, err := enc.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func generateRandomHWID() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "device-fallback"
	}
	return "device-" + hex.EncodeToString(buf)
}

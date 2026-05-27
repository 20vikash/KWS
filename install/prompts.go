package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// prompt displays a prompt with an optional default value and returns the user's input.
func prompt(label string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("  %s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("  %s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}

// promptRequired displays a prompt that requires a non-empty answer.
func promptRequired(label string) string {
	for {
		fmt.Printf("  %s: ", label)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("    ⚠ This field is required.")
	}
}

// promptPassword displays a prompt and auto-generates a random password if left empty.
func promptPassword(label string) string {
	fmt.Printf("  %s (press Enter to auto-generate): ", label)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		return input
	}
	pw := generateRandomPassword(24)
	fmt.Printf("    → Generated: %s\n", pw)
	return pw
}

// validateIP checks if the given string is a valid IP address.
func validateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// validateDomain performs a basic domain name check.
func validateDomain(domain string) bool {
	if domain == "" {
		return false
	}
	// Must contain at least one dot and no spaces
	if !strings.Contains(domain, ".") || strings.Contains(domain, " ") {
		return false
	}
	return true
}

// generateRandomPassword returns a random hex string of the given byte length.
func generateRandomPassword(byteLen int) string {
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback — this should never happen
		return "changeme_please_" + fmt.Sprintf("%d", os.Getpid())
	}
	return hex.EncodeToString(b)[:byteLen]
}

// collectDomainAndIP gathers the domain and public IP from the user.
func collectDomainAndIP(cfg *KWSConfig) {
	printSection("Domain & Network")

	for {
		cfg.Domain = promptRequired("Domain name (e.g., kwscloud.in)")
		if validateDomain(cfg.Domain) {
			break
		}
		fmt.Println("    ⚠ Invalid domain name. Must contain at least one dot and no spaces.")
	}

	for {
		cfg.PublicIP = promptRequired("Public IP address of this server")
		if validateIP(cfg.PublicIP) {
			break
		}
		fmt.Println("    ⚠ Invalid IP address.")
	}
}

// collectSSL derives SSL cert paths from the domain name.
func collectSSL(cfg *KWSConfig) {
	printSection("SSL Certificates")

	defaultCert := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", cfg.Domain)
	defaultKey := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", cfg.Domain)
	fmt.Println("  The certificate is wildcard by default (*.your-domain).")
	fmt.Println("  One cert/key pair covers all subdomains.")

	cfg.SSL.CertPath = prompt("SSL certificate path", defaultCert)
	cfg.SSL.KeyPath = prompt("SSL private key path", defaultKey)
}

// collectGmail gathers Gmail SMTP credentials.
func collectGmail(envCfg *EnvConfig) {
	printSection("Gmail SMTP (for email verification)")

	envCfg.GmailAddress = prompt("Gmail address", "")
	envCfg.GmailAppPassword = prompt("Gmail app password", "")

	if envCfg.GmailAddress == "" || envCfg.GmailAppPassword == "" {
		fmt.Println("    ⚠ Gmail not configured. Email verification will not work.")
	}
}

// collectPasswords generates or collects passwords for all services.
func collectPasswords(envCfg *EnvConfig) {
	printSection("Service Passwords")

	envCfg.DBPassword = promptPassword("PostgreSQL (main) password")
	envCfg.RedisPassword = promptPassword("Redis password")
	envCfg.MQPassword = promptPassword("RabbitMQ password")
	envCfg.PGServicePassword = promptPassword("PostgreSQL (user service) password")
}

// collectInstanceLimits lets users override default limits.
func collectInstanceLimits(cfg *KWSConfig) {
	printSection("Instance Limits")

	fmt.Println("  Press Enter to keep the defaults.")

	val := prompt("Instance memory limit", cfg.Instance.MemoryLimit)
	cfg.Instance.MemoryLimit = val
}

func printSection(title string) {
	fmt.Println()
	fmt.Printf("─── %s ───\n", title)
}

func collectAnsibleInfo(kwsCfg *KWSConfig) *AnsibleConfig {
	printSection("Ansible Server Provisioning")

	fmt.Println("  KWS can use Ansible to automatically set up your server.")
	fmt.Println("  Provide the SSH details of the target server, or press Enter to skip.")

	ans := &AnsibleConfig{}

	ans.TargetIP = prompt("Server IP address", kwsCfg.PublicIP)
	if ans.TargetIP == "" {
		fmt.Println("  ⚠ Skipping Ansible setup. You can run 'ansible-playbook -i ansible/inventory ansible/playbook.yaml' manually later.")
		return nil
	}

	ans.SSHUser = prompt("SSH user", "root")
	ans.SSHKeyPath = prompt("SSH private key path", "~/.ssh/id_rsa")

	return ans
}

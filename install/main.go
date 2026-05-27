package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	printBanner()

	exePath, err := os.Executable()
	if err != nil {
		exePath, _ = os.Getwd()
	}
	projectRoot := detectProjectRoot(exePath)

	fmt.Printf("Project root: %s\n", projectRoot)

	kwsCfg := DefaultKWSConfig()
	envCfg := DefaultEnvConfig()

	collectDomainAndIP(kwsCfg)
	collectSSL(kwsCfg)
	collectGmail(envCfg)
	collectPasswords(envCfg)
	collectInstanceLimits(kwsCfg)

	envCfg.ServicesSubnet = kwsCfg.Network.ServicesSubnet
	envCfg.ServicesGateway = kwsCfg.Network.ServicesGateway
	envCfg.PGServiceIP = kwsCfg.Services.PostgresIP
	envCfg.AdminerIP = kwsCfg.Services.AdminerIP
	envCfg.AttachServices = buildAttachServicesVar(kwsCfg)

	printSection("WireGuard Key Generation")
	wgPrivKey, wgPubKey, err := generateWireGuardKeys()
	if err != nil {
		fmt.Printf("  ⚠ Failed to generate WireGuard keys: %v\n", err)
		fmt.Println("  You will need to set WG_PRIVATE_KEY manually in .env")
		envCfg.WGPrivateKey = ""
		wgPrivKey = ""
		wgPubKey = ""
	} else {
		envCfg.WGPrivateKey = wgPrivKey
		envCfg.WGPublicKey = wgPubKey
		fmt.Printf("  ✓ Private key: (saved to .env)\n")
		fmt.Printf("  ✓ Public key:  %s\n", wgPubKey)
		fmt.Println("  ℹ Share the public key with clients for WireGuard peer configuration.")
	}

	ansCfg := collectAnsibleInfo(kwsCfg)

	printSection("Writing Configuration Files")

	cfgPath := filepath.Join(projectRoot, "kws_config.yaml")
	if err := writeKWSConfig(kwsCfg, cfgPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ %s\n", cfgPath)

	envPath := filepath.Join(projectRoot, ".env")
	if err := writeEnvFile(envCfg, envPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ %s\n", envPath)

	srcEnvPath := filepath.Join(projectRoot, "src", ".env")
	if err := writeEnvFile(envCfg, srcEnvPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
	} else {
		fmt.Printf("  ✓ %s\n", srcEnvPath)
	}

	nginxDir := filepath.Join(projectRoot, "nginx", "conf.d")
	os.MkdirAll(nginxDir, 0755)
	mainConfPath := filepath.Join(nginxDir, "main.conf")
	if err := generateNginxMainConf(kwsCfg, mainConfPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
	} else {
		fmt.Printf("  ✓ %s\n", mainConfPath)
	}

	defaultConfPath := filepath.Join(nginxDir, "00-default.conf")
	if err := generateNginxDefaultConf(kwsCfg, defaultConfPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
	} else {
		fmt.Printf("  ✓ %s\n", defaultConfPath)
	}

	dnsmasqPath := filepath.Join(projectRoot, "dnsmasq.conf")
	if err := generateDnsmasqConf(kwsCfg, dnsmasqPath); err != nil {
		fmt.Printf("  ✗ %v\n", err)
	} else {
		fmt.Printf("  ✓ %s\n", dnsmasqPath)
	}

	var invPath string
	if ansCfg != nil {
		ansCfg.WGPrivKey = wgPrivKey
		ansCfg.WGPubKey = wgPubKey
		invPath, err = generateAnsibleInventory(ansCfg, kwsCfg, projectRoot)
		if err != nil {
			fmt.Printf("  ✗ %v\n", err)
		} else {
			fmt.Printf("  ✓ %s\n", invPath)
		}
	}

	printSection("Setup Complete!")
	fmt.Println()
	fmt.Println("  Configuration files written successfully.")
	fmt.Println()

	if ansCfg != nil {
		fmt.Println("  Next: Provision the server with Ansible?")
		fmt.Printf("    ansible-playbook -i %s ansible/playbook.yaml\n", invPath)
		fmt.Println()
		run, err := confirm("  Run ansible-playbook now? (y/n)")
		if err != nil {
			fmt.Printf("    Skipping: %v\n", err)
		} else if run {
			runAnsiblePlaybook(projectRoot, invPath)
		} else {
			fmt.Println("    Skipped. Run the command above when ready.")
		}
		fmt.Println()
	}

	fmt.Println("  Remaining steps:")
	fmt.Println("    1. Ensure SSL certificates are in place for your domain")
	fmt.Println("    2. Run database migrations: make migrate_up")
	fmt.Println("    3. Start the platform:      make up")
	fmt.Println()
	fmt.Printf("  Domain:    %s\n", kwsCfg.Domain)
	fmt.Printf("  Public IP: %s\n", kwsCfg.PublicIP)
	fmt.Println()
}

func printBanner() {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║         KWS Platform Installer       ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  This tool will configure KWS for your server.")
	fmt.Println("  Press Enter to accept default values shown in [brackets].")
}

func detectProjectRoot(startPath string) string {
	dir := filepath.Dir(startPath)
	for i := 0; i < 5; i++ {
		if fileExists(filepath.Join(dir, "compose.yaml")) || fileExists(filepath.Join(dir, "Makefile")) {
			return dir
		}
		dir = filepath.Dir(dir)
	}

	cwd, _ := os.Getwd()
	if fileExists(filepath.Join(cwd, "compose.yaml")) {
		return cwd
	}
	parent := filepath.Dir(cwd)
	if fileExists(filepath.Join(parent, "compose.yaml")) {
		return parent
	}

	return cwd
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateWireGuardKeys() (privateKey, publicKey string, err error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate X25519 key: %w", err)
	}
	privateKey = base64.StdEncoding.EncodeToString(priv.Bytes())
	publicKey = base64.StdEncoding.EncodeToString(priv.PublicKey().Bytes())
	return privateKey, publicKey, nil
}

func buildAttachServicesVar(cfg *KWSConfig) string {
	var parts []string
	for _, ba := range cfg.Services.BridgeAttach {
		parts = append(parts, fmt.Sprintf("%s:%s:%s", ba.Container, ba.Bridge, ba.IPCIDR))
	}
	return strings.Join(parts, " ")
}

func confirm(msg string) (bool, error) {
	fmt.Printf("  %s [y/N]: ", msg)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}

func runAnsiblePlaybook(projectRoot, invPath string) {
	printSection("Running Ansible Playbook")

	playbook := filepath.Join(projectRoot, "ansible", "playbook.yaml")
	cmd := exec.Command("ansible-playbook", "-i", invPath, playbook)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("  Executing: ansible-playbook -i %s %s\n\n", invPath, playbook)

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n  ✗ Ansible playbook failed: %v\n", err)
		fmt.Println("  You can re-run it manually:")
		fmt.Printf("    ansible-playbook -i %s %s\n", invPath, playbook)
	} else {
		fmt.Println("\n  ✓ Server provisioning complete!")
	}
}
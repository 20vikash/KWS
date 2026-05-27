package nginx

import (
	"fmt"
	"html/template"
	"kws/kws/consts/config"
	"log"
	"os"
)

type Template struct {
	Domain     string
	IP         string
	Port       string
	BaseDomain string
	SSLCert    string
	SSLKey     string
}

// populateFromConfig fills in BaseDomain and SSL fields from the global config.
func (t *Template) populateFromConfig() {
	if t.BaseDomain == "" {
		t.BaseDomain = config.DOMAIN()
	}
	if t.SSLCert == "" {
		t.SSLCert = config.SSL_CERT_PATH()
	}
	if t.SSLKey == "" {
		t.SSLKey = config.SSL_KEY_PATH()
	}
}

const nginxTemplate = `
server {
    listen 80;
    server_name {{ .Domain }}.{{ .BaseDomain }};
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name {{ .Domain }}.{{ .BaseDomain }};

    ssl_certificate     {{ .SSLCert }};
    ssl_certificate_key {{ .SSLKey }};

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://{{ .IP }}:{{ .Port }};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_buffering off;
    }
}
`

const domainTemplate = `
server {
    listen 80;
    server_name {{ .Domain }};
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name {{ .Domain }};

    ssl_certificate     /etc/letsencrypt/live/{{ .Domain }}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/{{ .Domain }}/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://127.0.0.1:{{ .Port }};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # preserve client Host header
        proxy_set_header Host $http_host;

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_buffering off;
    }
}
`

func (t *Template) AddNewConf(templateType string) error {
	t.populateFromConfig()

	var finalTemplate string

	switch templateType {
	case config.INSTANCE_TEMPLATE:
		finalTemplate = nginxTemplate
	case config.DOMAIN_TEMPLATE:
		// For domain (tunnel) template, set the tunnel proxy port from config
		if t.Port == "" {
			t.Port = fmt.Sprintf("%d", config.TUNNEL_PROXY_PORT())
		}
		finalTemplate = domainTemplate
	}

	tmpl, err := template.New("nginx").Parse(finalTemplate)
	if err != nil {
		panic(err)
	}

	// Path to write the new Nginx config file
	filePath := fmt.Sprintf("%s/%s.conf", config.NGINX_CONF_DIR(), t.Domain)

	// Create and write to the file
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file nginx conf")
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, t); err != nil {
		log.Println("nginx templating failed")
		return err
	}

	return nil
}

func (t *Template) RemoveConf() error {
	// Path to the config file to be removed
	filePath := fmt.Sprintf("%s/%s.conf", config.NGINX_CONF_DIR(), t.Domain)

	// Attempt to remove the file
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file %s does not exist, nothing to remove.", filePath)
			return nil
		}
		log.Printf("Failed to remove config file %s: %v", filePath, err)
		return err
	}

	log.Printf("Config file %s successfully removed.", filePath)
	return nil
}


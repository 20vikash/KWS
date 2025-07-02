package nginx

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type Template struct {
	Domain string
	IP     string
	Port   string
}

const nginxTemplate = `
server {
    listen 80;
    server_name {{ .Domain }}.kwscloud.in;

    location / {
        proxy_pass http://{{ .IP }}:{{ .Port }};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
`

func (t *Template) AddNewConf() error {
	tmpl, err := template.New("nginx").Parse(nginxTemplate)
	if err != nil {
		panic(err)
	}

	// Path to write the new Nginx config file
	filePath := fmt.Sprintf("/app/nginx_conf/%s.conf", t.Domain)

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
	filePath := fmt.Sprintf("/app/nginx_conf/%s.conf", t.Domain)

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

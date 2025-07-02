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
	filePath := fmt.Sprintf("/etc/nginx/conf.d/%s.conf", t.Domain)

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

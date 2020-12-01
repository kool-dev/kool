package presets

// auto generated file

// GetTemplates get all templates
func GetTemplates() map[string]map[string]string {
	var templates = make(map[string]map[string]string)
	templates["cache"] = map[string]string{
		"memcached16.yml": `image: memcached:1.6-alpine
volumes:
  - cache:/data:delegated
networks:
  - kool_local
`,
		"redis60.yml": `image: redis:6-alpine
volumes:
  - cache:/data:delegated
networks:
  - kool_local
`,
	}
	templates["database"] = map[string]string{
		"mysql57.yml": `image: mysql:5.7
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
  MYSQL_DATABASE: "${DB_DATABASE:-database}"
  MYSQL_USER: "${DB_USERNAME:-user}"
  MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
  MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
volumes:
 - database:/var/lib/mysql:delegated
networks:
 - kool_local
`,
		"mysql80.yml": `image: mysql:8.0
command: --default-authentication-plugin=mysql_native_password
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  MYSQL_ROOT_PASSWORD: "${DB_PASSWORD:-rootpass}"
  MYSQL_DATABASE: "${DB_DATABASE:-database}"
  MYSQL_USER: "${DB_USERNAME:-user}"
  MYSQL_PASSWORD: "${DB_PASSWORD:-pass}"
  MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
volumes:
 - database:/var/lib/mysql:delegated
networks:
 - kool_local
`,
		"prostgresql130.yml": `image: postgres:13-alpine
ports:
  - "${KOOL_DATABASE_PORT:-3306}:3306"
environment:
  POSTGRES_DB: "${DB_DATABASE:-database}"
  POSTGRES_USER: "${DB_USERNAME:-user}"
  POSTGRES_PASSWORD: "${DB_PASSWORD:-pass}"
  POSTGRES_HOST_AUTH_METHOD: "trust"
volumes:
 - database:/var/lib/postgresql/data:delegated
networks:
 - kool_local
`,
	}
	return templates
}

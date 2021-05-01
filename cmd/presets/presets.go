package presets

// auto generated file

// GetAll get all presets
func GetAll() map[string]map[string]string {
	var presets = make(map[string]map[string]string)
	presets["adonis"] = map[string]string{}
	presets["golang-cli"] = map[string]string{
		"kool.yml": `scripts:
  # Helper for local development - compiling and installing locally
  dev:
    - kool run compile
    - kool run install

  # Runs go CLI with proper version for kool development
  go: kool docker --volume=gopath:/go --env='GOOS=$GOOS' golang:1.16.0 go

  # Compiling cli itself. In case you are on MacOS make sure to have your .env
  # file properly setting GOOS=darwin so you will be able to use the binary.
  compile: kool run go build -o my-cli
  install:
    - mv my-cli /usr/local/bin/my-cli
  fmt: kool run go fmt ./...
  lint: kool docker --volume=gopath:/go golangci/golangci-lint:v1.31.0 golangci-lint run -v
`,
	}
	presets["hugo"] = map[string]string{}
	presets["laravel"] = map[string]string{}
	presets["nestjs"] = map[string]string{}
	presets["nextjs"] = map[string]string{
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
      - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
	}
	presets["nodejs"] = map[string]string{
		"app.js": `const http = require("http");

const hostname = "0.0.0.0";
const port = 3000;

const server = http.createServer((req, res) => {
	res.statusCode = 200;
	res.setHeader("Content-Type", "text/plain");
	res.end("Hello World");
});

server.listen(port, hostname, () => {
	console.log("Server running at http://localhost:"+port+"/");
});
`,
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["node", "app.js"]
    ports:
      - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
	}
	presets["nuxtjs"] = map[string]string{
		"docker-compose.yml": `version: "3.7"
services:
  app:
    image: kooldev/node:14
    command: ["npm", "run", "dev"]
    ports:
     - "${KOOL_APP_PORT:-3000}:3000"
    environment:
      ASUSER: "${KOOL_ASUSER:-0}"
      UID: "${UID:-0}"
    volumes:
     - .:/app:delegated
    networks:
     - kool_local
     - kool_global

networks:
  kool_local:
  kool_global:
    external: true
    name: "${KOOL_GLOBAL_NETWORK:-kool_global}"
`,
	}
	presets["php"] = map[string]string{}
	presets["symfony"] = map[string]string{}
	presets["wordpress"] = map[string]string{}
	return presets
}

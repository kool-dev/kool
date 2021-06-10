# Generate PDFs

Easily add a `pdf` service container to your project that generates PDFs from any URL or HTML content.

## docker-compose.yml

```yaml
pdf:
  image: "kooldev/pdf:latest"
  expose:
    - 3000
  networks:
    - kool_local
```

### Full Example

```diff
version: "3.7"
services:
  app:
    image: kooldev/php:8.0-nginx
    ports:
      - ${KOOL_APP_PORT:-80}:80
    environment:
      ASUSER: ${KOOL_ASUSER:-0}
      UID: ${UID:-0}
    volumes:
      - .:/app:delegated
    networks:
      - kool_local
      - kool_global
  database:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    ports:
      - ${KOOL_DATABASE_PORT:-3306}:3306
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD-rootpass}
      MYSQL_DATABASE: ${DB_DATABASE-database}
      MYSQL_USER: ${DB_USERNAME-user}
      MYSQL_PASSWORD: ${DB_PASSWORD-pass}
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
    volumes:
      - database:/var/lib/mysql:delegated
    networks:
      - kool_local
    healthcheck:
      test:
        - CMD
        - mysqladmin
        - ping
+  pdf:
+    image: "kooldev/pdf:latest"
+    expose:
+      - 3000
+    networks:
+      - kool_local
volumes:
  database: null
networks:
  kool_local: null
  kool_global:
    external: true
    name: ${KOOL_GLOBAL_NETWORK:-kool_global}
```

## Usage

### From a URL:

Returns a rendered PDF from the provided URL (or JSON response with an error message and status).

Endpoint:
- `GET /from-url?url=`

Parameters:
- `url`: URL of the page you want to convert to a PDF

### From HTML content:

```
use GuzzleHttp\Client;

// `pdf` is the docker-compose service name
$pdf = (new Client())->post('http://pdf:3000/from-html', [
    'form_params' => [
        'html' => '<h1>This is my super kool HTML that I want to turn into an awesome PDF file!</h1><p>You get the idea of how this works <strong>:)</strong></p>',
        'options' => json_encode([
            'format' => 'A4',
            'printBackground' => false,
        ]),
    ],
])->getBody();

file_put_contents('path/to/my/super-kool.pdf', $pdf);
```

---

Read more at [https://github.com/kool-dev/pdf](https://github.com/kool-dev/pdf).

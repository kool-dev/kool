# Start a Spring Boot Project

Create a Spring Boot application using the Spring Initializr API:

```bash
kool create spring-boot my-spring-app
```

The wizard asks for the application coordinates, then fetches the current language, Java version, and dependency options from the Spring Initializr metadata API. This keeps the choices in sync with the provider instead of hardcoding them locally.

After the project is generated, database and cache services are optional. Available databases include PostgreSQL, MySQL, MariaDB, MongoDB, H2, and SQLite. H2 and SQLite use local file-backed data directories under `data/`; Redis and Memcached are available as cache services.

Run the generated application with:

```bash
kool run run
```

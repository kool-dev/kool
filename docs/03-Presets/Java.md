# Start a Java Project with Gradle

Create a Java application using Gradle and Docker:

```bash
kool create java my-java-app
```

The preset uses `gradle init` with Java 21, Kotlin build scripts, and JUnit Jupiter. It provides Docker-backed scripts for building, testing, and running the application.

During setup, database and cache services are optional. Database choices include PostgreSQL, MySQL, MariaDB, MongoDB, H2, and SQLite. H2 and SQLite create local file-backed data directories under `data/` instead of starting a database container. Redis and Memcached are available as cache services.

```bash
kool run build
kool run test
kool run run
```

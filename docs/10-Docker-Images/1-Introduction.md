When you start developing in containers, you quickly realize the official Docker images are built for deployment, and not for the special considerations (and nuances) of local development. One of the most common and recurring problems we see are permission issues with mapped volumes, due to host users being different from container users. Kool fixes this problem, and many others, by creating custom Docker images optimized for local development environments.

A few of the optimizations included in Kool's Docker images:
- UID mapping to host user to solve permission issues
- Alpine base images to remain small and up-to-date
- Configured with sane defaults - for development as well as production
- Environment variables to easily update the most common settings
- Battle-tested - a growing community has been using these images in production for quite a long time now!

## Kool-Optimized Docker Images

- PHP images: https://github.com/kool-dev/docker-php
- Nginx images: https://github.com/kool-dev/docker-nginx
- Node images: https://github.com/kool-dev/docker-node
- Java images: https://github.com/kool-dev/docker-java
- DevOps images: https://github.com/kool-dev/docker-toolkit

> Disclaimer: Kool Docker images follow our recommended best practices, which are aimed at making your life easier. However, you can use the **Kool CLI** with any Docker images you like - assuming you know what you're doing.

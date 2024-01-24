# Upgrade Guide

## Note on Versioning

**Kool** adheres to *Semantic Versioning (SemVer)* and is committed to maintaining backward compatibility, minimizing the potential for disruptive changes and multi-step upgrades.

## Upgrading to 3.x

The release of major version 3.x introduces a few breaking changes. Please review and adjust your projects accordingly.

### Changes to `kool.deploy.yml`: Now `kool.cloud.yml`

It is recommended to rename your configuration file for **Kool.dev Cloud** from `kool.deploy.yml` to `kool.cloud.yml`. Although the old naming convention remains functional, it has been deprecated and will be removed in future releases.

### Building Images for Deployment with `services.<service>.build`

Version 3.x introduces two significant changes:

- The YAML syntax for `services.<service>.build` in the `kool.cloud.yml` file must now align with the official Docker Compose reference for the `service.<service>.build` entry.
- Image building now occurs in your local environment—specifically, on the host where you execute `kool cloud deploy`. Therefore, ensure that the environment from which you run this command has a properly configured Docker-image build engine (that means Kool to be able to run `docker build` command).
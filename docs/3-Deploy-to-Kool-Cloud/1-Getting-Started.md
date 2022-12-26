The Kool Cloud supports a wide set of features designed to suite your dockerized web application deployment needs. It supports persisting folders across deployments, running daemons as extra containers, cron-like commands scheduling, hooks for before/after deployment, viewing running container logs, accessing the running container interactively, and much more.

The Deploy API was designed taking to heart the best Developer Experience for deploying containers to the cloud. By leveraging your existing local environment structure on `docker-compose.yml` and adding a sane and intuitive configuration layer that will even feel familiar from the first contact, the goal here is to provide a best-in-class offering for cloud platform engineering, offering you a platform to leverage Kubernetes and orchestrate your web applications in the cloud without all the hassle.

**kool cloud** is a the CLI suite of commands that allows you to configure, deploy, access and tail logs from the applications to Kool Cloud via the Deploy API.

## Deploy Docker Compose-based containerized apps in just a few simple steps

1. Sign up for Kool Cloud and get your access token.
	- Make sure your `.env` file contains two entries for configuring your deployment:
		- `KOOL_API_TOKEN` - the Deploy API access token you get from kool.dev web panel.
1. Configure your deploy with files straight in your application root folder
	- `kool.deploy.yml` - a "mirror" of your `docker-compose.yml` file, with extra pieces of data for your cloud deployment.
	- `Dockerfile` - usually you are going to need to build your app for deployment, in case you haven't yet.
	- Make sure you setup the necessary environment variables for your app to run in the cloud.
1. Deploy your application
	- Run `kool cloud deploy` - this will validate & deploy your applications.
	- Wait for it to finish and then access the provided deployment URL!
1. Doing more
	- View logs: `kool cloud logs` - you can check out the logs of your deployed containers.
	- Access cloud containers: `kool cloud exec` - you are able to execute commands - including interactive TTY sessions - within your cloud deployed containers.

---

Full documentation:

- [`kool.deploy.yml` Reference](/docs/3-Deploy-to-Kool-Cloud/2-kool.deploy.yml-Reference.md)

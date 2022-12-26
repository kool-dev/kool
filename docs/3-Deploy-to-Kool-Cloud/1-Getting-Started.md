The Kool Cloud supports a wide range of features designed to suit your needs for deploying dockerized web applications. It supports features such as persisting folders across deployments, running daemons as extra containers, scheduling commands like cron jobs, adding hooks for before and after deployment, viewing running container logs, accessing the running container interactively, and much more.

The Deploy API was designed with the best developer experience in mind for deploying containers to the cloud. By leveraging your existing local environment structure in `docker-compose.yml` and adding a sane and intuitive configuration layer that will feel familiar from the first contact, our goal is to provide a best-in-class offering for cloud platform engineering. This platform allows you to leverage Kubernetes and orchestrate your web applications in the cloud without all the hassle.

**kool cloud** is a CLI suite of commands that allows you to configure, deploy, access, and tail logs from the applications to Kool Cloud via the Deploy API.

## Deploy Docker Compose-based, containerized apps in just a few simple steps

1. Sign up for Kool Cloud and get your access token.
	- Make sure your `.env` file contains two entries for configuring your deployment:
		- `KOOL_API_TOKEN` - the Deploy API access token you get from the kool.dev web panel.
1. Configure your deployment with files straight in your application root folder
	- `kool.deploy.yml` - a "mirror" of your `docker-compose.yml` file, with extra pieces of data for your cloud deployment.
	- `Dockerfile` - usually, you are going to need to build your app for deployment if you haven't already.
	- Make sure you set up the necessary environment variables for your app to run in the cloud.
1. Deploy your application
	- Run `kool cloud deploy` - this will validate and deploy your application.
	- Wait for it to finish and then access the provided deployment URL!
1. Doing more
	- **View logs**
		- `kool cloud logs` - you can check the logs of your deployed containers.
	- **Access cloud containers**
		- `kool cloud exec` - you can execute commands, including interactive TTY sessions, within your cloud-deployed containers.

---

Full documentation:

- [`kool.deploy.yml` Reference](/docs/3-Deploy-to-Kool-Cloud/2-kool.deploy.yml-Reference.md)

# About Kool.dev Cloud

The [Kool.dev Cloud](https://kool.dev/cloud) brings the ease you always thought should exist when deploying web applications running in containers to the cloud.

For developers stepping into the world of cloud-native development, **Kool.dev** is a game-changer. It bridges the gap between the **complexity of container management** and the efficiency of cloud-native applications, empowering developers to focus more on creating value and less on managing infrastructure.

The Kool.dev Cloud API was designed with the best developer experience in mind for deploying containers to the cloud. By leveraging your existing local environment structure in `docker-compose.yml` and adding a sane and intuitive configuration layer that will feel familiar from the first sight, our goal is to provide a best-in-class offering for cloud platform engineering. This platform allows you to leverage Kubernetes and orchestrate your web applications in the cloud without all the hassle.

> **Enterprise**: You can use Kool.dev Cloud to deploy workloads to your own cloud vendor to keep things compliant - [contact us](mailto:contact@kool.dev) for the **"Bring your Own Cloud"** offer.

**Kool.dev Cloud** is the CLI suite of commands that allows you to configure, deploy, access, and tail logs from the applications to the cloud via the Kool.dev Cloud API.

# Features

## CPU and Memory Resources

Having your application in a container helps isolating its resources and allocate them more efficiently, while scaling up and down seamlessly.

Currently the sizing of containers go by a simple scale: `micro`, `small`, `medium`, `large`, `xlarge`, `xxlarge`. The `micro` size starts with **128 millicores** of CPU and **256Mbi** of memory limits. Each size then doubles the previous size limits.

## Persistent Storage

Kool.dev Cloud supports **persisting folders** across deployments. This is a key piece that facilitates moving legacy applications into containers.

## More that you may need

Supporting a wide range of features you will find no limits on any requirements your web application might have. Running **daemons** as extra containers, scheduling commands like **cron jobs**, adding **hooks to run before or after** deployment, **viewing logs** of running containers, accessing the running container **interactively**, and much more.

# Getting Started Steps

1. [Sign up for Kool.dev Cloud](https://kool.dev/register) and get your access token.
	- You can store your token in your `.env` file if you are using one:
		- `echo "KOOL_API_TOKEN=<my-token>" >> .env`
	- Or you can store your token in a real environment variable:
		- `export KOOL_API_TOKEN="<my token>"`
1. Configure your deployment with files directly in your application root folder. For that, you can use [`kool cloud setup`](/docs/commands-reference/kool-cloud-setup) to help guide you in creating the following files:
	- `kool.deploy.yml` - a "mirror" of your `docker-compose.yml` file, with extra pieces of data for customizing your cloud deployment.
	- `Dockerfile` - usually, you are going to need to build your app for deployment if you haven't already.
	- Make sure you set up the necessary [environment variables](/docs/kool-cloud/environment-variables) for your app to run in the cloud.
1. Deploy your application
	- Run `kool cloud deploy --domain=<your domain>` - this will validate and deploy your application.
	- Wait for it to finish and then access the provided deployment URL!
1. Doing more
	- **View logs**
		- `kool cloud logs` - you can check the logs of your deployed containers.
	- **Access running containers (like SSH-ing in)**
		- `kool cloud exec` - you can execute commands, including interactive TTY sessions, within your cloud-deployed containers. For example, `kool cloud exec app bash` to open a bash in my running container in the cloud.

---

Reference:

- [`kool.deploy.yml` Reference](/docs/kool-cloud/kool.deploy.yml-reference)

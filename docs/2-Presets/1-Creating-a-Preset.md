# Creating a new preset

`kool` has a builtin task-runner feature that allows us to automate steps/repetitive tasks - first of being all necessary setup steps for bootstrapping a new project. We use this feature to enable our presets, accomplishing two objectives:

- Keeping fast, clean and simple how to setup a new local development environment for popular frameworks and start coding.
- Have this project with all good practices and sane defaults for running in containers - for development and later deployment.

## Steps to create a preset

#### 1. Creating `presets/my-preset/config.yml`

Create a new folder under  `presets` and a `config.yml` file in it: `presets/my-preset/config.yml`.

The `config.yml` file is where we configure:

- Steps for creating a new project
- Steps for installing `kool` tailored local Docker environment to existing projects

Both of the two tasks described above are accomplished via a set of declarative `steps` and `actions` on what we can call Kool Automation language. Some of the actions include:

- `scripts` - running arbitrary shell script.
- `copy` - copying files from our preset of templates folder right into the local project.
- `merge` - merge YAML files - helpful for building `docker-compose.yml` or `kool.yml` dynamically.
- `recipe`: run a Recipe which is a group of steps/actions ready to reuse

You can find the full reference on the Kool Automation Langauge here TBD.

Check out some of our current presets as examples.

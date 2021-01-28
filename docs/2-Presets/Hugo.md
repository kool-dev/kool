### Using Hugo preset

#### Creating a new Hugo website

To make things easier we will use **kool** to install it for you.

```console
$ kool create hugo my-website

$ cd my-website
```
- **kool create** already executes **kool preset** internally so you can skip the command in the next step.

#### Adding kool to an existing Hugo website

Go to the project folder and run:

```console
$ cd my-website/
$ kool preset hugo
```

**kool preset** will create a few configuration files in order to enable you to configure / extend it. You don't need to execute it if you created the project with `kool create`.

### Using kool for Hugo development

- To start the container to serve your Hugo website:

```console
$ kool start
```

Then check out your site at `http://localhost`. If you wanna stop the container just run `kool stop`.

- To create some new content:

```console
$ kool run hugo new posts/my-super-post.md
```

---

Check your **kool.yml** to see what scripts you can run and add more.

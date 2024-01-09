All deployed environments will run under HTTPS. Kool.dev Cloud will automatically generate certificates for your environment the first time it is deployed using the Let'sEncrypt engine.

You always need to provide the environment domain when running a deploy.

```bash
kool cloud --token="<my token>" deploy --domain="my-app-domain.com"
```

You can provide those values via environment variables as well if that's easier for you:

- `KOOL_API_TOKEN` for the access token.
- `KOOL_DEPLOY_DOMAIN` for the domain environment you want to deploy to.

Important to notice: if you deploy to a new domain that doesn't currently exist in your Kool.dev Cloud panel, that is totally fine and will just create a very new environment for that domain.

### Test deployment domains

You are welcome to use a subdomain like `my-super-app.kool.cloud` on your staging or development environments. By using that, you will have HTTPS certificates up and running instantly for that environment after the first deploy.

### Production and custom domains

When you create an environment to be deployed using your own custom domain name, you will need to check out in the Kool.dev Cloud panel for that environment the instructions to where to point your A/CNAME for that domain.

HTTPS certificates will only be successfully generated once the DNS is correctly pointing your domain to Kool.dev Cloud.

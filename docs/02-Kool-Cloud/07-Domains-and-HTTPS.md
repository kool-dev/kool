# Managing domains and HTTPS

All deployed environments will run under HTTPS. Kool.dev Cloud will automatically generate certificates for your environment the first time it is deployed using the Let'sEncrypt engine (via K8S `certmanager`).

## Environment domain

You always need to provide the environment domain when running a deploy. The domain is going to be the key identifier of your environment when accessing it via the Kool CLI.

```bash
kool cloud --token="<my token>" deploy --domain="my-app-domain.com"
```

You can provide those values via environment variables as well if that's easier for you:

- `KOOL_API_TOKEN` for the access token.
- `KOOL_DEPLOY_DOMAIN` for the domain environment you want to deploy to.

Important to notice: if you deploy to a new domain that doesn't currently exist in your Kool.dev Cloud panel, that is totally fine and will just create a very new environment for that domain.

## Using multiple domains

Some applications may require serving more than just one domain. This can be easily achieve with Kool.dev Cloud.

When deploying via the CLI, you can specify any number of extra domains that should point to your application:

```bash
kool cloud --token="<my token>" deploy --domain="my-app-domain.com" --domain-extra="another-domain.com"
```

## Using a test deployment domain - `*.kool.cloud`

You are welcome to use a subdomain like `my-super-app.kool.cloud` on your staging or development environments. By using that, you will have HTTPS certificates up and running instantly for that environment after the first deploy.

## Custom domains - DNS records

When you create an environment to be deployed using your own custom domain name, you will need to check out in the Kool.dev Cloud panel for that environment the instructions to where to point your A/CNAME DNS records for that domain.

HTTPS certificates will only be successfully generated once the DNS is correctly pointing your domain to Kool.dev Cloud.

## HTTPS certificates

As stated above for all environments Kool.dev Cloud will trigger LetsEncrypt routines to generate a valid TLS certificate for your environment domain. All that is required for that to succeed is that the used domain points properly to Kool.dev Cloud IP addresses.

That is why it's important to notice that when deploying a custom domain (not a `*.kool.cloud`) your will only work properly after the DNS records point properly to Kool.dev Cloud and the LetsEncrypt process had time to perform the HTTP01 Acme challenge process.

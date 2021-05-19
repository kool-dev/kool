## kool share

Live share your local environment on the Internet using an HTTP tunnel

```
kool share
```

### Options

```
  -h, --help               help for share
      --port uint          The port from the target service that should be shared. If not provided, it will default to port 80.
      --service string     The name of the local service container you want to share. (default "app")
      --subdomain string   The subdomain used to generate your public https://subdomain.kool.live URL.
```

### Options inherited from parent commands

```
      --verbose   increases output verbosity
```

### Troubleshooting

**Problem:**

> Laravel cannot access http resource over https

**Answer:**

> In your App\Http\Middleware\TrustProxies.php make sure that the `$proxies` variable is as follows:

```php
    protected $proxies = '*';
```

-   Also make sure that if you are using laravel mix, use it accordingly.

This:

```html
<script src="{{ mix('js/app.js') }}" type="text/javascript"></script>
```

Not This:

```html
<script src="{{ asset('js/app.js') }}" type="text/javascript"></script>
```

Using it with asset might cause problems while using `kool share`.

### SEE ALSO

-   [kool](kool) - Cloud native environments made easy

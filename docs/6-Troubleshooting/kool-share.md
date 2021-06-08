### `kool share`

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
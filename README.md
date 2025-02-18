# ecommerce-go
2024s Unfinished freelance e-commerce website made with Go and TailwindCSS. This would be a real e-commerce, but the project was delayed undefinitely and the code severely changed.

## Flags and env used for the server start

<!-- Server port -->
#### Server address/port value flag
If empty, default value is :80.

```bash
-port=4000
```

<!-- Cloudinary connection string -->
#### Cloudinary connection string
Used for communicating with Cloudinary API, **cannot be empty, can be invalid keys**.

If empty, returns a connection error.

```bash
CLOUDINARY_URL="cloudinary://API_KEY:API_SECRET@CLOUD_NAME"
```

<!-- Render.com database connection string -->
#### Database connection string
Used for communicating with Postgresql database.

If empty, default value is:

```bash
DSN="host=localhost user=postgres password=root dbname=ecommerce_go_database port=5432 TimeZone=America/Sao_Paulo"
```

<!-- Domain -->
#### Domain string
Used to set the domain value. At least for now, it is needed for subdomains work.

If empty, default value is localhost.

```bash
DOMAIN="example.com"
```

<!-- Hash salt for hashids -->
#### Hash salt
Used for hashing and unhashing values.

```bash
HASH_SALT="abc123"
```

#### Default credentials for admin

```bash
DEFAULT_ADMIN_EMAIL="example"
DEFAULT_ADMIN_PASSWORD="example"
```

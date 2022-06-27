# Hello Traefik Plugin

This plugin return message, status code and headers with configuration and request query parameters.

Set default contents in configuration also override with request query parameters.

```sh
curl -i "http://localhost:8080/test?headers=Content-Type:application/json,X-Test:true&statusCode=202&message=Gunaydin"
```

```
HTTP/1.1 202 Accepted
Access-Control-Allow-Credentials: true
Content-Type: application/json
Referrer-Policy: no-referrer
X-Test: True
Date: Mon, 27 Jun 2022 10:32:55 GMT
Content-Length: 8

Gunaydin
```

## Configuration

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  plugins:
    hello:
      moduleName: github.com/worldline-go/traefik-plugin-hello
      version: v0.1.0
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: PathPrefix(`/test`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    my-plugin:
      plugin:
        hello:
          message: "I will ask you to write another 10 points there."
          statusCode: 200
          headers:
            Content-Type: "text/plain"
```

### Local Mode

Traefik also offers a developer mode that can be used for temporary testing of plugins not hosted on GitHub.
To use a plugin in local mode, the Traefik static configuration must define the module name (as is usual for Go packages) and a path to a [Go workspace](https://golang.org/doc/gopath_code.html#Workspaces), which can be the local GOPATH or any directory.

The plugins must be placed in `./plugins-local` directory,
which should be in the working directory of the process running the Traefik binary.
The source code of the plugin should be organized as follows:

```
./plugins-local/
    └── src
        └── github.com
            └── worldline-go
                └── traefik-plugin-hello
                    ├── .traefik.yml
                    ├── hello.go
                    ├── hello_test.go
                    ├── go.mod
                    ├── LICENSE
                    ├── Makefile
                    └── README.md
```

```yaml
# Static configuration

experimental:
  localPlugins:
    example:
      moduleName: github.com/worldline-go/traefik-plugin-hello
```

(In the above example, the `traefik-plugin-hello` plugin will be loaded from the path `./plugins-local/src/github.com/worldline-go/traefik-plugin-hello`.)

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    my-plugin:
      plugin:
        example:
          headers:
            Foo: Bar
```

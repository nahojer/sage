# sage

Sage provides a fast routing mechanism of HTTP requests to route values (typically HTTP handlers)
and is meant to be a building block of HTTP router/mux packages.

Parameterization of path segments is configurable, but the API is otherwise deliberately simple:

* no regex matching
* one route value per URL path and HTTP method pair
* prefix matching is supported, but there is no way to configure it

All of the documentation can be found on the [go.dev](https://pkg.go.dev/github.com/nahojer/sage?tab=doc) website.

Check out [httprouter](https://github.com/nahojer/httprouter). It uses sage under the hood to implement a HTTP router.

Is it Good? [Yes](https://news.ycombinator.com/item?id=3067434).

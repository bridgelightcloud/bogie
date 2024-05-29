package router

type Router struct {
	Routes map[string]Route
}

func NewRouter() Router {
	return Router{
		Routes: map[string]Route{},
	}
}

func (r *Router) AddRoute(route Route) {
	if _, ok := r.Routes[route.Signature]; ok {
		panic("Route already exists")
	}

	r.Routes[route.Signature] = route
}

package router

type Router struct {
	Routes map[string]Route
}

func NewRouter(routes []Route) Router {
	router := Router{
		Routes: map[string]Route{},
	}

	for _, route := range routes {
		router.addRoute(route)
	}

	return router
}

func (r *Router) addRoute(route Route) {
	if _, ok := r.Routes[route.Signature]; ok {
		panic("Route already exists")
	}

	r.Routes[route.Signature] = route
}

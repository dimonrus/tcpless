package tcpless

type (
	// Route definition
	Route string

	// Handler Procedure handler
	Handler func(client IClient)

	// Callback method
	handlerRegistry map[string]Handler

	// Route hook list
	routeHookRegistry map[string][]Handler
)

var (
	// Callback registry
	registry = make(handlerRegistry)

	// Registry for route hooks
	routeRegistry = make(routeHookRegistry)
)

// Handle hande route
func (r Route) Handle(route string, handler Handler) {
	registry[r.build(route)] = handler
	return
}

// Sub create sub route
func (r Route) Sub(route string) Route {
	return Route(r.build(route))
}

// build full route path
func (r Route) build(route string) string {
	if route == "" {
		return string(r)
	}
	return string(r) + "." + route
}

// Hook register route hook
func (r Route) Hook(handler Handler) Route {
	route := r.build("")
	routeRegistry[route] = append(routeRegistry[route], handler)
	return r
}

// Reg register new route
func (h Handler) Reg(route string, handler Handler) Handler {
	registry[route] = handler
	return h
}

// UnReg unregister route by name
func (h Handler) UnReg(route string) Handler {
	delete(registry, route)
	return h
}

// Route dedicate sub route
func (h Handler) Route(route string) Route {
	return Route(route)
}

// GetHooks Return hook list based on route
func (r routeHookRegistry) GetHooks(route string) []Handler {
	var result = make([]Handler, 0, 8)
	var j, k int
	l := len([]rune(route))
	for i := 0; i < l; i++ {
		if route[i] == '.' {
			if v, ok := routeRegistry[route[:i]]; ok {
				j += len(v)
				copy(result[k:j+1], v[:])
				k = j
			}
		}
	}
	return result[:k]
}

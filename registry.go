package tcpless

// Callback registry
var registry = make(handlerRegistry)

// Callback method
type handlerRegistry map[string]Handler

// Handler Procedure handler
type Handler func(client IClient)

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

package tcpless

import (
	"fmt"
	"testing"
)

func MyTestRoute(handler Handler) Handler {
	api := handler.Route("api")
	system := api.Sub("system")
	system.Handle("health", func(client IClient) {
		fmt.Println("health route")
	})
	v1 := api.Sub("v1")
	v1.Handle("user", func(client IClient) {
		fmt.Println("health route")
	})
	return handler
}

func TestHandler_Route(t *testing.T) {
	MyTestRoute(nil)
	if v, ok := registry["api.system.health"]; !ok {
		t.Fatal("route not exists api.system.health")
	} else {
		v(nil)
	}
}

func TestHooks(t *testing.T) {
	var handler Handler
	api := handler.Route("api")
	api.Hook(func(client IClient) {
		fmt.Println("Api Hook executed")
	})
	system := api.Sub("system")
	system.Handle("health", func(client IClient) {
		fmt.Println("health route")
	})
	system.Hook(func(client IClient) {
		fmt.Println("System Hook executed 1")
	})
	system.Handle("db", func(client IClient) {
		fmt.Println("db route")
	})
	system.Hook(func(client IClient) {
		fmt.Println("System Hook executed 2")
	})
	v1 := api.Sub("v1")
	v1.Hook(func(client IClient) {
		fmt.Println("API V1 Hook executed")
	})
	v1.Handle("user", func(client IClient) {
		fmt.Println("health route")
	})
	t.Log("--------API SYSTEM HOOKS")
	for _, hook := range routeRegistry.GetHooks("api.system.health") {
		hook(nil)
	}
	t.Log("--------API V1 HOOKS")
	for _, hook := range routeRegistry.GetHooks("api.v1.user") {
		hook(nil)
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkRouteHookRegistry_GetHooks
// BenchmarkRouteHookRegistry_GetHooks-8   	 4969099	       244.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHookRegistry_GetHooks(b *testing.B) {
	so.Do(func() {
		var handler Handler
		api := handler.Route("api")
		api.Hook(func(client IClient) {
			fmt.Println("Api Hook executed")
		})
		system := api.Sub("system")
		system.Handle("health", func(client IClient) {
			fmt.Println("health route")
		})
		system.Hook(func(client IClient) {
			fmt.Println("System Hook executed 1")
		})
		system.Handle("db", func(client IClient) {
			fmt.Println("db route")
		})
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hooks := routeRegistry.GetHooks("api.system.some")
		_ = hooks
	}
	b.ReportAllocs()
}

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkRoute_Handle
// BenchmarkRoute_Handle-8   	56593780	        21.35 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoute_Handle(b *testing.B) {
	var h Handler
	api := h.Route("api")
	system := api.Sub("system")
	system.Handle("health", func(client IClient) {
		fmt.Println("health route")
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		route := system.build("health")
		_ = route
	}
	b.ReportAllocs()
}

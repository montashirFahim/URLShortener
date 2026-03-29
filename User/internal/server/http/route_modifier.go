package http

// import (
//  "net/http"
// 	 github.com/getsentry/raven-go"
// )

type RouteModifier func(*Route)

// type handlerMiddleware func(http.Handler) http.Handler

// func addNewrelic(app newrelic.Application) RouteModifier {
// 	return func(r *Route) {
// 		if app == nil {
// 			return
// 		}
// 		_, nxtHandler := newrelic.WrapHandle(app, r.Pattern, r.Handler)
// 		r.Handler = nxtHandler
// 	}
// }

// func addSentry(cfg *cfg.Sentry) RouteModifier {
// 	return func(r *Route) {
// 		if cfg.Key == "" {
// 			return
// 		}

// 		r.Handler = raven.Recoverer(r.Handler)
// 	}
// }

// func addHubble(apm handlerMiddleware) RouteModifier {
// 	return func(r *Route) {
// 		r.Handler = apm(r.Handler)
// 	}
// }

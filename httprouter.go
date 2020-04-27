package adaptr

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const CtxRouteAuthorizedKey = ctxRouteAuthorizedType("routeAuthorized")

type ctxRouteAuthorizedType string

const CtxHttpRouterParamsKey = ctxHttpParamsCtxType("httpRouterParams")

type ctxHttpParamsCtxType string

const CtxRequestJsonStructKey = requestJsonStructType("reqJsonStruct")

type requestJsonStructType string

const CtxRequestBodyByteArrKey = requestBodyStringType("reqBodyString")

type requestBodyStringType string

const CtxRequestIdParamKey = requestIdParamType("reqIdParam")

type requestIdParamType string

const CtxTokenKey = ctxTokenKeyType("ctxTokenKey")

type ctxTokenKeyType string

const CtxTokenUserIdentKey = ctxTokenUserIdentType("ctxTokenUserIdentType")

type ctxTokenUserIdentType string

const CtxAuthorizationsKey = ctxTokenUserIdentType("ctxAuthorizationsType")

type ctxAuthorizationsType string

const CtxNamespaceKey = ctxTokenUserIdentType("ctxNamespaceType")

type ctxNamespaceType string

/*const CtxTokenAudienceKey = ctxTokenAudienceType("ctxTokenAudienceType")

type ctxTokenAudienceType string

const CtxWorkspaceIdentObjKey = ctxTokenAudienceType("ctxWorkspaceIdentObjType")

type ctxWorkspaceIdentObjType string*/

//var once = sync.Once{}

// Simple accepts the name of a function so you don't have to wrap it with http.HandlerFunc
// Example: r.GET("/", httprouterwrapper.Simple(controller.Index))
func compatibleHandlerFn(h http.HandlerFunc, httprParamsCtxKey interface{}) httprouter.Handle {
	return toHttpRouterHandle(http.Handler(h), httprParamsCtxKey)
}

// Compatible accepts a handler to make it compatible with http.HandlerFunc
// Example: r.GET("/", httprouterwrapper.Compatible(http.HandlerFunc(controller.Index)))
func compatibleHandler(h http.Handler, httprParamsCtxKey interface{}) httprouter.Handle {
	return toHttpRouterHandle(h, httprParamsCtxKey)
}
func toHttpRouterHandle(h http.Handler, httprParamsCtxKey interface{}) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if httprParamsCtxKey != nil {
			r = SetCtxValue(r, httprParamsCtxKey, p)
		}

		h.ServeHTTP(w, r)
	}
}

func HttprouterAdaptFn(f http.HandlerFunc, httprParamsCtxKey interface{}, adapters ...Adapter) httprouter.Handle {
	return HttprouterAdapt(http.HandlerFunc(f), httprParamsCtxKey, adapters...)
}
func HttprouterAdapt(h http.Handler, httprParamsCtxKey interface{}, adapters ...Adapter) httprouter.Handle {
	h = Adapt(h, adapters...)
	return compatibleHandler(h, httprParamsCtxKey)
}

func WrapHandleFuncAdapters(hFn http.HandlerFunc, adapters []Adapter, preAdaptrs []Adapter, postAdaptrs []Adapter) httprouter.Handle {
	//to beginning
	if preAdaptrs != nil {
		adapters = append(preAdaptrs, adapters...)
	}
	adapters = append(adapters, toAdapter(hFn))
	//to end
	if postAdaptrs != nil {
		adapters = append(adapters, postAdaptrs...)
	}

	return HttprouterAdaptFn(emptyHandlerFn, CtxHttpRouterParamsKey, adapters...)
}

func CreateOptionsRouterHandle(corsAdapter Adapter) httprouter.Handle {
	return WrapHandleFuncAdapters(emptyHandlerFn,
		[]Adapter{corsAdapter, AuthPermitAll(CtxRouteAuthorizedKey)},
		nil,
		nil,
	)
}

func emptyHandlerFn(w http.ResponseWriter, r *http.Request) {}

func toAdapter(handlerFunc http.HandlerFunc) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerFunc(w, r)
			h.ServeHTTP(w, r)
		})
	}
}

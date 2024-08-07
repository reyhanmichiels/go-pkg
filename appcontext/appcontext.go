package appcontext

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/header"
	"github.com/reyhanmichiels/go-pkg/language"
)

type contextKey string

const (
	acceptLanguage   contextKey = "AcceptLanguage"
	requestId        contextKey = "RequestId"
	serviceVersion   contextKey = "ServiceVersion"
	userAgent        contextKey = "UserAgent"
	userId           contextKey = "UserId"
	requestStartTime contextKey = "RequestStartTime"
	appResponseCode  contextKey = "AppResponseCode"
	deviceType       contextKey = "DeviceType"
	appErrorMessage  contextKey = "AppErrorMessage"
	cacheControl     contextKey = "CacheControl"
	eventName        contextKey = "EventName"
	eventDescription contextKey = "EventDescription"
	requestBody      contextKey = "RequestBody"
	requestURI       contextKey = "RequestURI"
	requestMethod    contextKey = "RequestMethod"
	requestQuery     contextKey = "RequestQuery"
	clientIP         contextKey = "ClientIP"
	responseHttpCode contextKey = "ResponseHttpCode"
	authToken        contextKey = "AuthToken"
	serviceName      contextKey = "ServiceName"
)

func SetAcceptLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, acceptLanguage, lang)
}

func GetAcceptLanguage(ctx context.Context) string {
	lang, ok := ctx.Value(acceptLanguage).(string)
	if !ok {
		// return english as the default language
		return language.English
	}

	return lang
}

func SetCacheControl(ctx context.Context, cache string) context.Context {
	return context.WithValue(ctx, cacheControl, cache)
}

func GetCacheControl(ctx context.Context) bool {
	c, ok := ctx.Value(cacheControl).(string)
	if !ok || c != header.CacheControlNoCache {
		// return false as the default vale
		return false
	} else {
		return true
	}
}

func SetRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestId, rid)
}

func GetRequestId(ctx context.Context) string {
	rid, ok := ctx.Value(requestId).(string)
	if !ok {
		return ""
	}
	return rid
}

func SetServiceVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, serviceVersion, version)
}

func GetServiceVersion(ctx context.Context) string {
	version, ok := ctx.Value(serviceVersion).(string)
	if !ok {
		return ""
	}
	return version
}

func SetUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, userAgent, ua)
}

func GetUserAgent(ctx context.Context) string {
	ua, ok := ctx.Value(userAgent).(string)
	if !ok {
		return ""
	}
	return ua
}

func SetUserId(ctx context.Context, ui int) context.Context {
	return context.WithValue(ctx, userId, ui)
}

func GetUserId(ctx context.Context) int {
	ui, ok := ctx.Value(userId).(int)
	if !ok {
		return 0
	}
	return ui
}

func SetRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTime, t)
}

func GetRequestStartTime(ctx context.Context) time.Time {
	t, _ := ctx.Value(requestStartTime).(time.Time)
	return t
}

func SetAppResponseCode(ctx context.Context, code codes.Code) context.Context {
	return context.WithValue(ctx, appResponseCode, code)
}

func GetAppResponseCode(ctx context.Context) codes.Code {
	code, _ := ctx.Value(appResponseCode).(codes.Code)
	return code
}

func GetAppResponseCodeInt(ctx context.Context) int {
	code, _ := ctx.Value(appResponseCode).(codes.Code)
	return int(code)
}

func SetDeviceType(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, deviceType, platform)
}

func GetDeviceType(ctx context.Context) string {
	platform, ok := ctx.Value(deviceType).(string)
	if !ok {
		// return web as the default device type
		return "web"
	}

	return platform
}

func SetAppErrorMessage(ctx context.Context, errMsg string) context.Context {
	return context.WithValue(ctx, appErrorMessage, errMsg)
}

func GetAppErrorMessage(ctx context.Context) string {
	errMsg, ok := ctx.Value(appErrorMessage).(string)
	if !ok {
		return ""
	}

	return errMsg
}

func SetEventName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, eventName, name)
}

func GetEventName(ctx context.Context) string {
	name, ok := ctx.Value(eventName).(string)
	if !ok {
		return ""
	}

	return name
}

func SetEventDescription(ctx context.Context, description string) context.Context {
	return context.WithValue(ctx, eventDescription, description)
}

func GetEventDescription(ctx context.Context) string {
	desc, ok := ctx.Value(eventDescription).(string)
	if !ok {
		return ""
	}

	return desc
}

func SetRequestBody(ctx context.Context, body any) context.Context {
	return context.WithValue(ctx, requestBody, body)
}

func GetRequestBody(ctx context.Context) any {
	return ctx.Value(requestBody)
}

func SetRequestURI(ctx context.Context, uri string) context.Context {
	return context.WithValue(ctx, requestURI, uri)
}

func GetRequestURI(ctx context.Context) string {
	uri, ok := ctx.Value(requestURI).(string)
	if !ok {
		return ""
	}

	return uri
}

func SetRequestQuery(ctx context.Context, query string) context.Context {
	return context.WithValue(ctx, requestQuery, query)
}

func GetRequestQuery(ctx context.Context) string {
	query, ok := ctx.Value(requestQuery).(string)
	if !ok {
		return ""
	}

	return query
}

func SetRequestMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, requestMethod, method)
}

func GetRequestMethod(ctx context.Context) string {
	method, ok := ctx.Value(requestMethod).(string)
	if !ok {
		return ""
	}

	return method
}

func SetRequestIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIP, ip)
}

func GetRequestIP(ctx context.Context) string {
	ip, ok := ctx.Value(clientIP).(string)
	if !ok {
		return ""
	}

	return ip
}

func SetResponseHttpCode(ctx context.Context, code int) context.Context {
	return context.WithValue(ctx, responseHttpCode, code)
}

func GetResponseHttpCode(ctx context.Context) int {
	code, _ := ctx.Value(responseHttpCode).(int)
	return code
}

func SetAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authToken, token)
}

func GetAuthToken(ctx context.Context) string {
	lang, ok := ctx.Value(authToken).(string)
	if !ok {
		return ""
	}

	return lang
}

func SetServiceName(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, serviceName, val)
}

func GetServiceName(ctx context.Context) string {
	val, ok := ctx.Value(serviceName).(string)
	if !ok {
		return ""
	}
	return val
}

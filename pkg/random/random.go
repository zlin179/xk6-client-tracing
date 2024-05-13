package random

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"net/http"
	"time"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

var (
	letters             = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	httpStatusesSuccess = []int64{200, 201, 202, 204}
	httpStatusesError   = []int64{400, 401, 403, 404, 405, 406, 408, 409, 410, 411, 412, 413, 414, 415, 417, 428, 427, 500, 501, 502}
	httpMethods         = []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut, http.MethodPatch}
	httpContentTypes    = []string{"application/json", "application/xml", "application/x-www-form-urlencoded", "text/plain", "text/html"}
	operations          = []string{"get", "list", "query", "search", "set", "add", "create", "update", "send", "remove", "delete"}
	serviceSuffix       = []string{"", "", "service", "backend", "api", "proxy", "engine"}
	dbNames             = []string{"redis", "mysql", "postgres", "memcached", "mongodb", "elasticsearch"}
	resources           = []string{
		"order", "payment", "customer", "product", "stock", "inventory",
		"shipping", "billing", "checkout", "cart", "search", "analytics"}

	rnd *rand.Rand
	rndMutex = &sync.Mutex{}
)

func init() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(int64(^uint64(0)>>1)))
	rnd = rand.New(rand.NewSource(seed.Int64()))
	rndMutex.Lock()
}

func SelectElement[T any](elements []T) T {
	rndMutex.Lock()
	defer rndMutex.Unlock()
	return elements[rnd.Intn(len(elements))]
}

func String(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = SelectElement(letters)
	}
	return string(s)
}

func K6String(n int) string {
	return "k6." + String(n)
}

func IntBetween(min, max int) int {
	rndMutex.Lock()
	n := rnd.Intn(max - min)
	rndMutex.Unlock()
	return min + n
}

func Duration(min, max time.Duration) time.Duration {
	rndMutex.Lock()
	n := rnd.Int63n(int64(max) - int64(min))
	rndMutex.Unlock()
	return min + time.Duration(n)
}

func IPAddr() string {
	return "192.168.1.1"
}

func Port() int {
	return IntBetween(8000, 9000)
}

func HTTPStatusSuccess() int64 {
	return SelectElement(httpStatusesSuccess)
}

func HTTPStatusErr() int64 {
	return SelectElement(httpStatusesError)
}

func HTTPMethod() string {
	return SelectElement(httpMethods)
}

func HTTPContentType() []any {
	return []any{SelectElement(httpContentTypes)}
}

func DBService() string {
	return SelectElement(dbNames)
}

func Service() string {
	resource := SelectElement(resources)
	return ServiceForResource(resource)
}

func ServiceForResource(resource string) string {
	name := resource
	suffix := SelectElement(serviceSuffix)
	if suffix != "" {
		name = name + "-" + suffix
	}
	return name
}

func Operation() string {
	resource := SelectElement(resources)
	return OperationForResource(resource)
}

func OperationForResource(resource string) string {
	op := SelectElement(operations)
	return op + "-" + resource
}

func TraceID() pcommon.TraceID {
	var b [16]byte
	rndMutex.Lock()
	_, _ = rnd.Read(b[:]) // always returns nil error
	rndMutex.Unlock()
	return b
}

func SpanID() pcommon.SpanID {
	var b [8]byte
	rndMutex.Lock()
	_, _ = rnd.Read(b[:]) // always returns nil error
	rndMutex.Unlock()
	return b
}

package main

import (
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // Service URLs
    userServiceURL := "http://user-service:8081"
    orderServiceURL := "http://order-service:8082"

    // Routes
    router.Any("/users/*path", reverseProxy(userServiceURL))
    router.Any("/orders/*path", reverseProxy(orderServiceURL))

    router.Run(":8080")
}

func reverseProxy(target string) gin.HandlerFunc {
    url, _ := url.Parse(target)
    proxy := httputil.NewSingleHostReverseProxy(url)

    return func(c *gin.Context) {
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}
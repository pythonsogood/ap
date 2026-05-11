package middleware

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	rateLimtierKeyFormat = "doctor:ratelimit:%s:%d"
)

func extractClientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)

	if !ok || p.Addr == nil {
		return ""
	}

	host, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return p.Addr.String()
	}

	return host
}

type RateLimiter struct {
	rdb       *redis.Client
	limit_rpm int64
}

func NewRateLimiter(rdb *redis.Client, limit_rpm uint) *RateLimiter {
	if limit_rpm == 0 {
		limit_rpm = 100
	}

	return &RateLimiter{rdb: rdb, limit_rpm: int64(limit_rpm)}
}

func (r *RateLimiter) GRPCUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if r == nil || r.rdb == nil {
			return handler(ctx, req)
		}

		clientIP := extractClientIP(ctx)
		if clientIP == "" {
			clientIP = "unknown"
		}

		now := time.Now().UTC()
		windowStart := now.Truncate(time.Minute).Unix()
		key := fmt.Sprintf(rateLimtierKeyFormat, clientIP, windowStart)

		count, err := r.rdb.Incr(ctx, key).Result()

		if err != nil {
			log.Println(err.Error())

			return handler(ctx, req)
		}

		if count == 1 {
			if err := r.rdb.Expire(ctx, key, time.Minute).Err(); err != nil {
				log.Println(err.Error())
			}
		}

		if count > r.limit_rpm {
			nextWindow := now.Truncate(time.Minute).Add(time.Minute)
			retryAfter := max(int64(nextWindow.Sub(now).Seconds()), 1)

			log.Printf("[rate-limit] blocked %s %s %d %d %ds\n", info.FullMethod, clientIP, count, r.limit_rpm, retryAfter)

			return nil, status.Error(codes.ResourceExhausted, fmt.Sprintf("rate limit exceeded: max %d requests/minute; retry after %d seconds", r.limit_rpm, retryAfter))
		}

		return handler(ctx, req)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/appointment/cmd/appointment/config"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/cache"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/database"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/event"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/handler"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/middleware"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/service"
	grpc_transport "github.com/pythonsogood/ap-assignment1/appointment/internal/transport/grpc"
	http_transport "github.com/pythonsogood/ap-assignment1/appointment/internal/transport/http"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func serve_http(server_addr string, appointment_service service.AppointmentService, event_publisher event.EventPublisher) (*gin.Engine, handler.AppointmentHTTPHandler, error) {
	router := gin.Default()

	appointment_handler := handler.NewAppointmentHTTPHandler(appointment_service, event_publisher)

	if err := http_transport.SetupAppointmentTransport(router, appointment_handler); err != nil {
		return router, appointment_handler, err
	}

	if err := router.Run(server_addr); err != nil {
		return router, appointment_handler, err
	}

	return router, appointment_handler, nil
}

func serve_grpc(server_addr string, appointment_service service.AppointmentService, event_publisher event.EventPublisher, rate_limiter *middleware.RateLimiter) (*net.Listener, *handler.AppointmentGRPCHandler, error) {
	lis, err := net.Listen("tcp", server_addr)

	if err != nil {
		return &lis, nil, err
	}

	appointment_handler := handler.NewAppointmentGRPCHandler(appointment_service, event_publisher)

	var s *grpc.Server

	if rate_limiter != nil {
		s = grpc.NewServer(grpc.UnaryInterceptor(rate_limiter.GRPCUnaryServerInterceptor()))
	} else {
		s = grpc.NewServer()
	}

	err = grpc_transport.SetupAppointmentdTransport(s, appointment_handler)

	if err != nil {
		return &lis, appointment_handler, err
	}

	if err := s.Serve(lis); err != nil {
		return &lis, appointment_handler, err
	}

	return &lis, appointment_handler, nil
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	var appointment_cache_repo cache.AppointmentCacheRepository
	var rate_limiter *middleware.RateLimiter

	switch conf.Cache.Type {
	case config.CacheTypeNone:
		log.Println("Caching disabled")
	case config.CacheTypeRedis:
		opts, err := redis.ParseURL(conf.Cache.Redis.Url)

		if err != nil {
			log.Println(err.Error())
			break
		}

		rdb := redis.NewClient(opts)

		appointment_cache_repo = cache.NewRedisAppointmentCacheRepository(rdb, time.Duration(conf.Cache.Ttl)*time.Second)

		rate_limiter = middleware.NewRateLimiter(rdb, conf.Server.RateLimitRpm)
	default:
		log.Println("Unsupported cache type!")
	}

	var appointment_db *sql.DB
	var appointment_repo repository.AppointmentRepository
	var event_publisher event.EventPublisher

	switch conf.Database.Type {
	case config.DatabaseTypeSQLite:
		appointment_db, err = database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

		if err != nil {
			panic(err.Error())
		}

		appointment_repo = repository.NewSQLiteAppointmentRepository(appointment_db)
	case config.DatabaseTypePostgres:
		m, err := migrate.New(
			"file://migrations",
			fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.Database.Postgres.User, conf.Database.Postgres.Password, conf.Database.Postgres.Host, conf.Database.Postgres.Port, conf.Database.Postgres.Db),
		)

		if err != nil {
			panic(err.Error())
		}

		if err := m.Up(); err != nil && err.Error() != "no change" {
			panic(err.Error())
		}

		appointment_db, err = database.PostgresConnectDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Database.Postgres.Host, conf.Database.Postgres.Port, conf.Database.Postgres.User, conf.Database.Postgres.Password, conf.Database.Postgres.Db))

		if err != nil {
			panic(err.Error())
		}

		appointment_repo = repository.NewPostgresAppointmentRepository(appointment_db)
	default:
		panic("Unsupported database type!")
	}

	switch conf.MessageBroker.Type {
	case config.MessageBrokerTypeNATS:
		nc, err := nats.Connect(conf.MessageBroker.Nats.ConnectionUrl)

		if err != nil {
			panic(err.Error())
		}

		event_publisher = event.NewNATSEventPublisher(nc)
	default:
		panic("Unsupported message broker type!")
	}

	server_addr := fmt.Sprintf(":%d", conf.Server.Port)

	// http_client := http.Client{
	// 	Timeout: time.Duration(conf.Services.Doctor.Timeout*time.Second),
	// }

	// doctor_service := service.NewHTTPDoctorService(conf.Services.Doctor.Address, &http_client)
	doctor_service := service.NewGRPCDoctorService(conf.Services.Doctor.Address, time.Duration(conf.Services.Doctor.Timeout)*time.Second)

	appointment_service := service.NewAppointmentService(appointment_repo, appointment_cache_repo, doctor_service)

	// _, _, err = serve_http(server_addr, appointment_service, event_publisher)

	// if err != nil {
	// 	panic(err.Error())
	// }

	_, _, err = serve_grpc(server_addr, appointment_service, event_publisher, rate_limiter)

	if err != nil {
		panic(err.Error())
	}
}

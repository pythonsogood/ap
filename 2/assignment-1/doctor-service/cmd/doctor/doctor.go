package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/doctor/cmd/doctor/config"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/cache"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/database"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/event"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/repository"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
	grpc_transport "github.com/pythonsogood/ap-assignment1/doctor/internal/transport/grpc"
	http_transport "github.com/pythonsogood/ap-assignment1/doctor/internal/transport/http"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func serve_http(server_addr string, doctor_service service.DoctorService, event_publisher event.EventPublisher) (*gin.Engine, handler.DoctorHTTPHandler, error) {
	router := gin.Default()

	doctor_handler := handler.NewDoctorHTTPHandler(doctor_service, event_publisher)

	if err := http_transport.SetupDoctorTransport(router, doctor_handler); err != nil {
		return router, doctor_handler, err
	}

	if err := router.Run(server_addr); err != nil {
		return router, doctor_handler, err
	}

	return router, doctor_handler, nil
}

func serve_grpc(server_addr string, doctor_service service.DoctorService, event_publisher event.EventPublisher) (*net.Listener, *handler.DoctorGRPCHandler, error) {
	lis, err := net.Listen("tcp", server_addr)

	if err != nil {
		return &lis, nil, err
	}

	doctor_handler := handler.NewDoctorGRPCHandler(doctor_service, event_publisher)

	s := grpc.NewServer()

	err = grpc_transport.SetupDoctorTransport(s, doctor_handler)

	if err != nil {
		return &lis, doctor_handler, err
	}

	if err := s.Serve(lis); err != nil {
		return &lis, doctor_handler, err
	}

	return &lis, doctor_handler, nil
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	var doctor_cache_repo cache.DoctorCacheRepository

	switch conf.Cache.Type {
	case config.CacheTypeRedis:
		opts, err := redis.ParseURL(conf.Cache.Redis.Url)

		if err != nil {
			log.Println(err.Error())
			break
		}

		rdb := redis.NewClient(opts)

		doctor_cache_repo = cache.NewRedisDoctorCacheRepository(rdb, time.Duration(conf.Cache.Ttl))
	default:
		log.Println("Unsupported cache type!")
	}

	var doctor_db *sql.DB
	var doctor_repo repository.DoctorRepository
	var event_publisher event.EventPublisher

	switch conf.Database.Type {
	case config.DatabaseTypeSQLite:
		doctor_db, err = database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

		if err != nil {
			panic(err.Error())
		}

		doctor_repo = repository.NewSQLiteDoctorRepository(doctor_db)
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

		doctor_db, err = database.PostgresConnectDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Database.Postgres.Host, conf.Database.Postgres.Port, conf.Database.Postgres.User, conf.Database.Postgres.Password, conf.Database.Postgres.Db))

		if err != nil {
			panic(err.Error())
		}

		doctor_repo = repository.NewPostgresDoctorRepository(doctor_db)
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

	doctor_service := service.NewDoctorService(doctor_repo, doctor_cache_repo)

	// _, _, err = serve_http(server_addr, doctor_service, event_publisher)

	// if err != nil {
	// 	panic(err.Error())
	// }

	_, _, err = serve_grpc(server_addr, doctor_service, event_publisher)

	if err != nil {
		panic(err.Error())
	}
}

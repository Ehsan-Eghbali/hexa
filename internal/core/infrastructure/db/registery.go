package db

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"hexagonal/internal/core/infrastructure/config"
	"hexagonal/internal/core/port/dbPO"
	"log"
)

type Registry struct {
	databases map[string]dbPO.Database
}

// NewDBRegistry creates a new DBRegistry and registers all dbPO adapters.
func NewDBRegistry(cfg *config.Config) (*Registry, error) {
	registry := &Registry{
		databases: make(map[string]dbPO.Database),
	}

	// Register all databases from config
	for _, dbConfig := range cfg.Databases {
		var db dbPO.Database
		switch dbConfig.Type {
		case "postgres":
			db = &PostgresAdapter{Config: &PostgresConfig{
				Host:     dbConfig.Host,
				Port:     dbConfig.Port,
				User:     dbConfig.User,
				Password: dbConfig.Password,
				Name:     dbConfig.Name,
			}}
		case "redis":
			db = &RedisAdapter{Config: &RedisConfig{
				Host:     dbConfig.Host,
				Port:     dbConfig.Port,
				Password: dbConfig.Password,
			}}
		default:
			return nil, fmt.Errorf("unsupported dbPO type: %s", dbConfig.Type)
		}

		// Connect to the dbPO
		if err := db.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", dbConfig.Type, err)
		}

		// Run migrations if necessary
		if err := db.RunMigrations(); err != nil {
			return nil, fmt.Errorf("failed to run migrations for %s: %w", dbConfig.Type, err)
		}

		registry.databases[dbConfig.Type] = db
	}
	return registry, nil
}

// GetDatabase retrieves a registered dbPO by type.
func (r *Registry) GetDatabase(dbType string) (dbPO.Database, error) {
	db, exists := r.databases[dbType]
	if !exists {
		return nil, fmt.Errorf("dbPO of type %s not registered", dbType)
	}
	return db, nil
}

// GetRedisClient retrieves the Redis client from the registry.
func (r *Registry) GetRedisClient() *redis.Client {
	db, err := r.GetDatabase("redis")
	if err != nil {
		log.Fatalf("failed to connect to redis: %s", err)
	}
	return db.(*RedisAdapter).Client
}

// GetRedisCache retrieves the Redis cache client from the registry.
func (r *Registry) GetRedisCache() *redis.Client {
	db, err := r.GetDatabase("redis-cache")
	if err != nil {
		log.Fatalf("failed to connect to redis-cache: %s", err)
	}
	return db.(*RedisAdapter).Client

}

// GetPostgres retrieves the Postgres client from the registry.
func (r *Registry) GetPostgres() *sql.DB {
	db, err := r.GetDatabase("postgres")
	if err != nil {
		log.Fatalf("failed to connect to postgres: %s", err)
	}
	return db.(*PostgresAdapter).DB
}

// CloseAll closes all registered databases.
func (r *Registry) CloseAll() error {
	for _, db := range r.databases {
		if err := db.Close(); err != nil {
			return fmt.Errorf("error closing dbPO: %w", err)
		}
	}
	return nil
}

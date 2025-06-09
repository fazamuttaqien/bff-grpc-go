package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Advice struct {
	UserId    string
	Advice    string
	CreatedAt time.Time
}

var _ = loadLocalEnv()
var (
	db   = GetEnv("POSTGRES_DB")
	user = GetEnv("POSTGRES_USER")
	pwd  = GetEnv("POSTGRES_PASSWORD")
	host = GetEnv("POSTGRES_HOST")
)

func NewClient(ctx context.Context) (*pgxpool.Pool, error) {
	url := "postgres://" + user + ":" + pwd + "@" + host + "/" + db
	client, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, errors.New("cannot connect to postgres instance")
	}
	return client, nil
}

func CreateOne(client *pgxpool.Pool, ctx context.Context, advice *Advice) error {
	query := `
		INSERT INTO advices(user_id, advice, created_at) 
		VALUES($1, $2, CURRENT_TIMESTAMP)
	`
	_, err := client.Exec(ctx, query, advice.UserId, advice.Advice)
	return err
}

func UpdateOne(client *pgxpool.Pool, ctx context.Context, advice *Advice) error {
	query := `
		UPDATE advices 
		SET advice=$1, created_at=CURRENT_TIMESTAMP 
		WHERE user_id=$2
	`
	_, err := client.Exec(ctx, query, advice.Advice, advice.UserId)
	return err
}

func FindOne(client *pgxpool.Pool, ctx context.Context, id string) (*Advice, error) {
	advice := Advice{UserId: id}

	query := `
		SELECT advice, created_at FROM advices
		WHERE user_id=$1
	`

	if err := client.QueryRow(ctx, query, id).Scan(&advice.Advice, &advice.CreatedAt); err != nil {
		return nil, err
	}
	return &advice, nil
}

func loadLocalEnv() any {
	if _, runningInContainer := os.LookupEnv("ADVICE_GRPC_SERVICE"); !runningInContainer {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Environment variable not found: ", key)
	}
	return value
}

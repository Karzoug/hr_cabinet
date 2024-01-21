package integration

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	mincont "github.com/Employee-s-file-cabinet/backend/integration/container/minio"
	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	pgcont "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	mhcont "github.com/Employee-s-file-cabinet/backend/integration/container/mailhog"
	"github.com/Employee-s-file-cabinet/backend/internal/app"
	"github.com/Employee-s-file-cabinet/backend/internal/config"
)

const (
	host   = "localhost"
	domain = "localhost"

	envType  = "development"
	logLevel = "debug"
	baseURL  = "/api/v1"

	recoveryFrom = "recovery@ecabinet.acceleratorpracticum.ru"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(fmt.Errorf("integration tests: couldn't change directory: %w", err))
	}
}

type UserTestSuite struct {
	suite.Suite
	ctx        context.Context
	containers struct {
		pg      *pgcont.PostgresContainer
		minio   *mincont.Container
		mailhog *mhcont.Container
	}
	port int
}

func TestUserRepoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests: skipping in short mode")
	}
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// postgresql
	pgContainer, err := pgcont.RunContainer(suite.ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		pgcont.WithDatabase("ecabinet"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	if err != nil {
		log.Fatal("postgresql: create container: ", err)
	}
	suite.containers.pg = pgContainer
	if err := suite.fillDatabase(); err != nil {
		log.Fatal("postgresql: fill database: ", err)
	}

	// minio
	minioContainter, err := mincont.RunContainer(suite.ctx)
	if err != nil {
		log.Fatal("minio: create container: ", err)
	}
	suite.containers.minio = minioContainter

	// mailhog
	mailhogContainer, err := mhcont.RunContainer(suite.ctx)
	if err != nil {
		log.Fatal("mailhog: create container: ", err)
	}
	suite.containers.mailhog = mailhogContainer

	// server
	suite.setEnvs()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to get config: %s", err.Error())
	}
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	}))
	go func() {
		if err := app.Run(suite.ctx, cfg, logger); err != nil {
			logger.Error("app stopped with error", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		}
	}()

	require.NoError(suite.T(), suite.checkHealth())
}

func (suite *UserTestSuite) TearDownSuite() {
	if err := suite.containers.pg.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	if err := suite.containers.minio.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating minio container: %s", err)
	}
	if err := suite.containers.mailhog.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mailhog container: %s", err)
	}
}

func (suite *UserTestSuite) serverURL() string {
	return fmt.Sprintf("http://%s:%d%s", host, suite.port, baseURL)
}

func (suite *UserTestSuite) setEnvs() {
	t := suite.T()

	connStr, err := suite.containers.pg.ConnectionString(suite.ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	t.Setenv("PG_DSN", connStr)

	_, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	t.Setenv("HTTP_TOKEN_SECRET_KEY", hex.EncodeToString(private))

	suite.port, err = getFreePort()
	if err != nil {
		log.Fatal(err)
	}
	t.Setenv("HTTP_PORT", fmt.Sprintf("%d", suite.port))

	t.Setenv("S3_HOST", suite.containers.minio.Host)
	t.Setenv("S3_PORT", fmt.Sprintf("%d", suite.containers.minio.Port))
	t.Setenv("S3_ACCESS_KEY_ID", suite.containers.minio.AccessKeyID)
	t.Setenv("S3_SECRET_ACCESS_KEY", suite.containers.minio.SecretAccessKey)
	t.Setenv("ENV_TYPE", envType)
	t.Setenv("LOG_LEVEL", logLevel)
	t.Setenv("MAIL_SMTP_HOST", suite.containers.mailhog.Host)
	t.Setenv("MAIL_SMTP_PORT", fmt.Sprintf("%d", suite.containers.mailhog.Port))
	t.Setenv("MAIL_FROM", recoveryFrom)
	t.Setenv("RECOVERY_DOMAIN", domain)
}

func (suite *UserTestSuite) checkHealth() error {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	resp, err := retryClient.Get(suite.serverURL() + "/health")
	if err != nil {
		return fmt.Errorf("failed to server check health: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to server check health: got status code: %s", resp.Status)
	}
	return nil
}

func (suite *UserTestSuite) fillDatabase() error {
	connStr, err := suite.containers.pg.ConnectionString(suite.ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	tx, err := db.BeginTx(suite.ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, err := os.ReadFile("integration/testdata/init.sql")
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(suite.ctx, string(query)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

var (
	pg       *PgRepository
	testUser = model.User{
		ID:           uuid.New(),
		Login:        "testLogin",
		Password:     []byte("testPassword"),
		RefreshToken: "",
	}
)

func SetupTestPostgres() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=profileuser",
		"POSTGRES_PASSWORD=profilepassword",
		"POSTGRES_DB=profiledb"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	err = RunMigrations(resource.GetPort("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	dbURL := fmt.Sprintf("postgres://profileuser:profilepassword@localhost:%s/profiledb", resource.GetPort("5432/tcp"))
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func RunMigrations(port string) error {
	cmd := exec.Command("flyway", "-url=jdbc:postgresql://localhost:"+port+"/profiledb", "-user=profileuser", "-password=profilepassword", "-locations=filesystem:../../migrations", "-connectRetries=10", "migrate")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	dbpool, cleanupPostgres, err := SetupTestPostgres()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
		cleanupPostgres()
		os.Exit(1)
	}
	pg = NewPgRepository(dbpool)
	exitVal := m.Run()
	cleanupPostgres()
	os.Exit(exitVal)
}

func TestSignUpUser(t *testing.T) {
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
}

func TestGetByLoginUser(t *testing.T) {
	testUser.ID = uuid.New()
	testUser.Login = "testLogin2"
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	password, id, err := pg.GetByLogin(context.Background(), testUser.Login)
	require.NoError(t, err)
	require.Equal(t, password, testUser.Password)
	require.NotEqual(t, id, uuid.Nil)
}

func TestAddRefreshToken(t *testing.T) {
	testUser.ID = uuid.New()
	testUser.Login = "testLogin3"
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	err = pg.AddRefreshToken(context.Background(), testUser.ID, testUser.RefreshToken)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	testUser.ID = uuid.New()
	testUser.Login = "testLogin4"
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	err = pg.DeleteAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetEmptyRefreshTokenByID(t *testing.T) {
	testUser.ID = uuid.New()
	testUser.Login = "testLogin6"
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	err = pg.AddRefreshToken(context.Background(), testUser.ID, testUser.RefreshToken)
	require.NoError(t, err)
	refreshToken, err := pg.GetRefreshTokenByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Empty(t, refreshToken)
}

func TestGetRefreshTokenByID(t *testing.T) {
	testUser.ID = uuid.New()
	testUser.Login = "testLogin5"
	err := pg.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	testUser.RefreshToken = "testRefreshToken"
	err = pg.AddRefreshToken(context.Background(), testUser.ID, testUser.RefreshToken)
	require.NoError(t, err)
	refreshToken, err := pg.GetRefreshTokenByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, refreshToken, testUser.RefreshToken, "testRefreshToken")
}

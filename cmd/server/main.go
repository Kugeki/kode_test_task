package main

import (
	"context"
	"github.com/Kugeki/kode_test_task/internal/adapters/pgstore"
	"github.com/Kugeki/kode_test_task/internal/adapters/yaspeller"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"github.com/Kugeki/kode_test_task/internal/ports/rest"
	"github.com/Kugeki/kode_test_task/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lmittmann/tint"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	spellClientTimeout = 10 * time.Second
	logLevel           = slog.LevelDebug

	restAddr         = ":8080"
	restReadTimeout  = 10 * time.Second
	restWriteTimeout = 10 * time.Second

	dbURL = ""

	jwtExpireDuration = 12 * time.Hour
	jwtSecretKey      = "jwt.secret.key//kode.test.task"

	shutdownTimeout = 15 * time.Second
)

func main() {
	log := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.StampMilli,
		AddSource:  true,
		Level:      logLevel,
	}))
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbURL = os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Error("POSTGRES_URL environment variable is not provided")
		return
	}

	spellClientOpts := []yaspeller.ClientOpt{
		yaspeller.WithTimeout(spellClientTimeout),
		yaspeller.WithLogger(log, slog.LevelInfo),
	}
	spellClient, err := yaspeller.NewClient(spellClientOpts...)
	if err != nil {
		log.Error("can't create yandex spell client", slog.Any("error", err))
		return
	}

	checkOpts := []yaspeller.CheckOpt{
		yaspeller.WithCheckFormat("plain"),
		yaspeller.WithCheckIgnoreURLs(),
		yaspeller.WithCheckIgnoreDigits(),
	}
	spellClient.AddCheckOptions(checkOpts...)

	store, err := pgstore.New(ctx, log, dbURL)
	if err != nil {
		log.Error("can't create postgres store", slog.Any("error", err))
		return
	}
	defer store.Close()

	authUsecase := usecases.NewAuthUC(store.User())
	noteUsecase := usecases.NewNoteUC(spellClient, store.Note())

	err = createUsers(ctx, authUsecase)
	if err != nil {
		log.Error("users has not created", slog.Any("error", err))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(rest.NewLogFormatter(log, slog.LevelInfo)))
	router.Use(middleware.Recoverer)

	authHandler := rest.NewAuthHandler(log, authUsecase, jwtSecretKey, jwtExpireDuration)
	authHandler.SetupRoutes(router)

	noteHandler := rest.NewNoteHandler(log, noteUsecase, jwtSecretKey)
	noteHandler.SetupRoutes(router)

	restOpts := []rest.Opt{
		rest.WithAddr(restAddr),
		rest.WithErrorLog(slog.NewLogLogger(log.Handler(), slog.LevelError)),
		rest.WithReadTimeout(restReadTimeout),
		rest.WithWriteTimeout(restWriteTimeout),
	}
	restServer, err := rest.NewServer(router, log, restOpts...)

	go func() {
		err := restServer.Run()
		if err != nil && err != http.ErrServerClosed {
			log.Error("rest server run error", slog.Any("error", err))
		}
	}()

	<-ctx.Done()

	log.Info("graceful shutdown is beginning...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	done := make(chan struct{}, 1)

	go func() {
		err := restServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Error("rest server shutdown error", slog.Any("error", err))
		}
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Error("graceful shutdown has failed", slog.Any("error", shutdownCtx.Err()))
	case <-done:
		log.Info("graceful shutdown has completed successfully")
	}
}

func createUsers(ctx context.Context, uc *usecases.AuthUC) error {
	err := uc.CreateUser(ctx, &domain.User{Name: "test1"}, "password1")
	if err != nil {
		return err
	}
	err = uc.CreateUser(ctx, &domain.User{Name: "test2"}, "password2")
	if err != nil {
		return err
	}
	err = uc.CreateUser(ctx, &domain.User{Name: "test3"}, "password3")
	if err != nil {
		return err
	}
	err = uc.CreateUser(ctx, &domain.User{Name: "test4"}, "password4")
	if err != nil {
		return err
	}

	return nil
}

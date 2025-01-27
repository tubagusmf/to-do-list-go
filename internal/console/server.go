package console

import (
	"log"
	"net/http"
	"sync"
	"to-do-list/db"
	"to-do-list/internal/config"
	"to-do-list/internal/repository"
	"to-do-list/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	handlerHttp "to-do-list/internal/delivery/http"
)

func init() {
	rootCmd.AddCommand(serverCMD)
}

var serverCMD = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP server",
	Long:  "Start the HTTP server to handle incoming requests for the to-do list application.",
	Run:   httpServer,
}

func httpServer(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	postgresDB := db.NewPostgres()
	sqlDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
	}
	defer sqlDB.Close()

	userRepo := repository.NewUserRepo(postgresDB)
	taskRepo := repository.NewTaskRepo(postgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)

	e := echo.New()

	handlerHttp.NewUserHandler(e, userUsecase)
	handlerHttp.NewTaskHandler(e, taskUsecase)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		errCh <- e.Start(":3000")
	}()

	go func() {
		defer wg.Done()
		<-errCh
	}()

	wg.Wait()

	if err := <-errCh; err != nil {
		if err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}
}

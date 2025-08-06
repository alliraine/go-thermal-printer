package bootstrap

import "fmt"

func Bootstrap() error {
	svc, err := initServices()
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	router, err := initRouter(svc)
	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	err = router.Run("127.0.0.1:8080")
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

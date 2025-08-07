package bootstrap

import (
	"fmt"
)

func Bootstrap() error {
	svc, err := initServices()
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	router, err := initRouter(svc)
	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	serverConfig := svc.configService.GetServerConfig()
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	err = router.Run(addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

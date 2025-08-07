package bootstrap

import (
	"fmt"

	"github.com/jonasclaes/go-thermal-printer/pkg/service"
)

type services struct {
	configService  *service.ConfigService
	printService   *service.PrintService
	printerService *service.PrinterService
}

func initServices() (svc *services, err error) {
	svc = &services{}

	svc.configService, err = service.NewConfigService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config service: %w", err)
	}

	svc.printService, err = service.NewPrintService(svc.configService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize print service: %w", err)
	}

	svc.printerService, err = service.NewPrinterService(svc.printService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize printer service: %w", err)
	}

	return svc, nil
}

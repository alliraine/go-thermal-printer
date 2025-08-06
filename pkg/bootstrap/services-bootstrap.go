package bootstrap

import (
	"fmt"

	"github.com/jonasclaes/go-thermal-printer/pkg/service"
)

type services struct {
	printService   *service.PrintService
	printerService *service.PrinterService
}

func initServices() (svc *services, err error) {
	svc = &services{}

	svc.printService, err = service.NewPrintService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize print service: %w", err)
	}

	svc.printerService, err = service.NewPrinterService(svc.printService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize printer service: %w", err)
	}

	return svc, nil
}

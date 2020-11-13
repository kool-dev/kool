package cmd

import (
	"kool-dev/kool/cmd/updater"

	"time"
)

// UpdateAwareService holds functions to implement the checker to see if theres a new version available
type UpdateAwareService struct {
	KoolService

	updater updater.Updater
}

// UpdateWrapper wraps the service with checker logic
func UpdateWrapper(service KoolService) *UpdateAwareService {
	return &UpdateAwareService{
		service,
		&updater.DefaultUpdater{RootCommand: rootCmd},
	}
}

// Execute runs the check logic and proxies to original service
func (u *UpdateAwareService) Execute(args []string) (err error) {
	if !u.KoolService.IsTerminal() {
		if err = u.KoolService.Execute(args); err != nil {
			return err
		}
		return
	}

	ch := make(chan bool)

	go u.updater.CheckForUpdates(u.updater.GetCurrentVersion(), ch)

	if err = u.KoolService.Execute(args); err != nil {
		return err
	}

	select {
	case update := <-ch:
		if update {
			defer u.KoolService.Warning("Theres a new Kool Version available! Run kool self-update to update!")
		}
	case <-time.After(300 * time.Millisecond):
		break
	}
	close(ch)

	return
}

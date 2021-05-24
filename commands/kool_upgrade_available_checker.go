package commands

import (
	"kool-dev/kool/services/updater"

	"time"
)

// UpdateAwareService holds functions to implement the checker to see if theres a new version available
type UpdateAwareService struct {
	KoolService

	updater updater.Updater
}

// CheckNewVersion wraps the service with checker logic
func CheckNewVersion(service KoolService, updater updater.Updater) *UpdateAwareService {
	return &UpdateAwareService{
		service,
		updater,
	}
}

// Execute runs the check logic and proxies to original service
func (u *UpdateAwareService) Execute(args []string) (err error) {
	if !u.KoolService.IsTerminal() {
		err = u.KoolService.Execute(args)
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
			defer u.KoolService.Warning("There's a new version available! Run kool self-update to upgrade!")
		}
	case <-time.After(1000 * time.Millisecond):
		break
	}

	return
}

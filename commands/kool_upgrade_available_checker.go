package commands

import (
	"kool-dev/kool/services/updater"

	"time"
)

// UpdateAwareService holds functions to implement the checker to see if theres a new version available
type UpdateAwareService struct {
	KoolService

	updater updater.Updater
	skip    bool
}

// CheckNewVersion wraps the service with checker logic
func CheckNewVersion(service KoolService, updater updater.Updater, skip bool) *UpdateAwareService {
	return &UpdateAwareService{
		service,
		updater,
		skip,
	}
}

// Execute runs the check logic and proxies to original service
func (u *UpdateAwareService) Execute(args []string) (err error) {
	if u.skip || !u.KoolService.Shell().IsTerminal() {
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
			defer u.KoolService.Shell().Warning("There's a new version available! Run kool self-update to upgrade!")
		}
	case <-time.After(time.Second):
		break
	}

	return
}

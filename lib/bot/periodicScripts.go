package bot

import (
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
)

type PeriodicScript struct {
	Name        string
	CronSpec    string
	Function    func()
	cronEntryId cron.EntryID
}

func RegisterPeriodicScript(script PeriodicScript) error {
	if periodicScripts == nil {
		periodicScripts = make(map[string]PeriodicScript)
	}
	if cronJobs == nil {
		cronJobs = cron.New(cron.WithSeconds())
	}

	cronJobs.Stop()

	// If script has previously been registerd, then remove the associated cron func
	if _, ok := periodicScripts[script.Name]; ok {
		cronJobs.Remove(periodicScripts[script.Name].cronEntryId)
		delete(periodicScripts, script.Name)
	}

	cronEntryId, err := cronJobs.AddFunc(script.CronSpec, script.Function)
	if err != nil {
		logger.Errorf("Unable to register periodic script %s", err)
	} else {
		script.cronEntryId = cronEntryId
		periodicScripts[script.Name] = script
	}

	if len(cronJobs.Entries()) > 0 {
		cronJobs.Start()
	} else {
		cronJobs.Stop()
	}

	return err
}

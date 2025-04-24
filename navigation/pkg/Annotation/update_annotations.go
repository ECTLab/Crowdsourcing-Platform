package Annotation


import (
    "context"
    log "github.com/sirupsen/logrus"
	"navigation/pkg/Clients/Redis"
    "time"
)

func UpdateOnlineReports(ctx context.Context) {
	go func() {
		log.Info("trying to update online reports")
		for {
			latestVersion, err := Redis.Crowdsourcing.GetOnlineReportsLatestVersion(ctx)
			if err != nil {
				log.WithError(err).Error("Couldn't get the latest Online Reports Version")
				time.Sleep(1 * time.Minute)
				continue
			}

			if Redis.CurrentOnlineReportsVersion == "" || latestVersion != Redis.CurrentOnlineReportsVersion {
				policeData, accidentData, err := Redis.Crowdsourcing.GetLatestOnlineReportsData(ctx)
				if err != nil || policeData == nil {
					log.WithError(err).Error("Couldn't get the latest Online Reports hash set")
					time.Sleep(1 * time.Minute)
					continue
				}

				Redis.PoliceData = policeData
				Redis.AccidentData = accidentData

				Redis.CurrentOnlineReportsVersion = latestVersion
				log.Infof("Online Reports is now updated to version: %s", latestVersion)
				time.Sleep(1 * time.Minute)
				log.Info("online reports updated")
			} else {
				time.Sleep(1 * time.Minute)
			}
		}
	}()
}

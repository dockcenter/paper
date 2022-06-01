package discover

import (
	"fmt"
	"golang.org/x/mod/semver"
)

func BuildCommand(promotion Promotion, environment string) string {
	if environment == "" {
		environment = semver.Canonical(promotion.Semver())[1:]
	}
	return fmt.Sprintf("drone build promote \"$DRONE_REPO\" \"$DRONE_BUILD_NUMBER\" %s --param=DOWNLOAD_URL=%s --param=DOCKER_TAGS=%s", environment, promotion.DownloadURL, promotion.DockerTags())
}

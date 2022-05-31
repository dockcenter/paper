package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/mod/semver"

	. "github.com/dockcenter/paper/internal/app/discover"
)

func main() {
	client := resty.New()
	const PROJECT string = "paper"
	const SUPPORTED_VERSION_GROUP int = 2
	const DOWNLOADS_KEY string = "application"

	// Parse environment variables
	event := os.Getenv("DRONE_BUILD_EVENT")
	branch := os.Getenv("DRONE_BRANCH")
	duration, err := time.ParseDuration(os.Getenv("DURATION"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Duration:", duration)

	// Get paper versions
	var project ProjectResponse
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", PROJECT)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(resp.Body(), &project)

	// When pushing to main, promote all supported versions' latest build
	var promotions []Promotion
	if event == "push" && branch == "main" {
		// Iterate all version groups
		var versions []string
		for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SUPPORTED_VERSION_GROUP:] {
			// Get all versions in version group
			var versionFamily VersionFamilyResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s", PROJECT, versionGroup)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			json.Unmarshal(resp.Body(), &versionFamily)

			versions = append(versions, versionFamily.Versions...)
		}

		// Get all builds of versions
		var semverMap = make(map[string]string)
		for i, version := range versions {
			// Get all builds for specific version
			var builds BuildsResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds", PROJECT, version)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			json.Unmarshal(resp.Body(), &builds)

			// Build promotion
			var promotion Promotion
			promotion.Version = builds.Version

			// Select latest build in each version
			build := builds.Builds[len(builds.Builds)-1]
			promotion.Build = build.Build
			promotion.DownloadURL = fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, promotion.Version, promotion.Build, build.Downloads[DOWNLOADS_KEY].Name)

			// Mark last version build as latest and Major
			if i == len(versions)-1 {
				promotion.Latest = true
				promotion.Major = true
			} else {
				promotion.Latest = false
				promotion.Major = false
			}

			// Compare semver and mark MajorMinor
			key := semver.MajorMinor(promotion.Semver())
			ver, ok := semverMap[key]
			if ok {
				if semver.Compare(promotion.Semver(), ver) > 0 {
					semverMap[key] = promotion.Semver()
				}
			} else {
				semverMap[key] = promotion.Semver()
			}

			promotions = append(promotions, promotion)
		}

		// mark MajorMinor based on semverMap
		for i, promotion := range promotions {
			key := semver.MajorMinor(promotion.Semver())
			promotions[i].MajorMinor = (semverMap[key] == promotion.Semver())
		}
	}

	for _, promotion := range promotions {
		fmt.Println(promotion.DockerTags())
	}
}

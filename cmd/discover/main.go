package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	. "github.com/dockcenter/paper/internal/app/discover"
	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	const PROJECT string = "paper"
	const SupportedVersionGroup int = 2
	const DownloadsKey string = "application"

	// Parse environment variables
	event := os.Getenv("DRONE_BUILD_EVENT")
	branch := os.Getenv("DRONE_BRANCH")
	duration, err := time.ParseDuration(os.Getenv("DURATION"))
	environment := os.Getenv("ENVIRONMENT")
	if err != nil {
		panic(err)
	}
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Duration:", duration)
	fmt.Println("Promotion environment:", environment)

	// Get paper versions
	var project ProjectResponse
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", PROJECT)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &project)
	if err != nil {
		panic(err)
	}

	// When pushing to main, promote all supported versions' latest build
	// Otherwise, get all builds in duration in supported vrsion groups
	var promotions []Promotion
	if event == "push" && branch == "main" {
		// Iterate all version groups
		var versions []string
		for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SupportedVersionGroup:] {
			// Get all versions in version group
			var versionFamily VersionFamilyResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s", PROJECT, versionGroup)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &versionFamily)
			if err != nil {
				panic(err)
			}

			versions = append(versions, versionFamily.Versions...)
		}

		// Get all builds of versions
		//var semverMap = make(map[string]string)
		for _, version := range versions {
			// Get all builds for specific version
			var builds BuildsResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds", PROJECT, version)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &builds)
			if err != nil {
				panic(err)
			}

			// Build promotion
			var promotion Promotion
			promotion.Version = builds.Version

			// Select latest build in each version
			build := builds.Builds[len(builds.Builds)-1]
			promotion.Build = build.Build
			promotion.DownloadURL = fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, promotion.Version, promotion.Build, build.Downloads[DownloadsKey].Name)
			promotions = append(promotions, promotion)
		}
	} else {
		// Get all builds for supported version groups
		//semverMap := make(map[string]string)
		for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SupportedVersionGroup:] {
			var versionFamilyBuilds VersionFamilyBuildsResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s/builds", PROJECT, versionGroup)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &versionFamilyBuilds)
			if err != nil {
				panic(err)
			}

			builds := versionFamilyBuilds.Builds
			for _, build := range builds {
				// Filter out builds that are longer than duration
				if time.Since(build.Time) > duration {
					continue
				}

				// Build promotion and append to promotions
				var promotion Promotion
				promotion.Version = build.Version
				promotion.Build = build.Build
				promotion.DownloadURL = fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, promotion.Version, promotion.Build, build.Downloads[DownloadsKey].Name)
				promotions = append(promotions, promotion)
			}
		}
	}

	MarkSemver(promotions)

	// Print tags to promote
	fmt.Println("\nTags to promote:")
	for _, promotion := range promotions {
		fmt.Println(promotion.DockerTags())
	}

	// Build promote commands and write to scripts/promote.sh
	cmd := "#!/bin/sh\n\n"
	for _, promotion := range promotions {
		cmd += BuildCommand(promotion, environment) + "\n"
	}

	// Write to scripts/promote.sh
	err = os.MkdirAll("scripts", 0700)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("scripts/promote.sh", []byte(cmd), 0700)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nShell script has been generated to scripts/promote.sh")
}

package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dockcenter/paper/internal/app/discover"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Parse environment variables
	event := os.Getenv("GITHUB_EVENT_NAME")
	branch := os.Getenv("GITHUB_REF_NAME")
	dryRunStr := os.Getenv("DRY_RUN")
	dryRun, err := strconv.ParseBool(dryRunStr)
	if err != nil {
		dryRun = false
	}
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Dry run:", dryRun)

	client := resty.New()

	// Get paper version groups
	var project ProjectResponse
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", Project)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &project)
	if err != nil {
		panic(err)
	}

	// Get paper builds
	var builds []VersionFamilyBuild
	for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SupportedVersionGroups:] {
		var versionFamilyBuilds VersionFamilyBuildsResponse
		url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s/builds", Project, versionGroup)
		resp, err := client.R().Get(url)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(resp.Body(), &versionFamilyBuilds)
		if err != nil {
			panic(err)
		}

		builds = append(builds, versionFamilyBuilds.Builds...)
	}

	// Get pushed tags
	var tags []string
	dockerTags := GetExistingTags(DockerRepository)
	for _, dockerTag := range dockerTags {
		tags = append(tags, dockerTag.Name)
	}

	var eventForImageInfo Event
	if event == "push" && branch == "main" {
		eventForImageInfo = Rebuild
	} else {
		eventForImageInfo = Cron
	}
	imageInfo := BuildImageInfo(builds, tags, eventForImageInfo)

	// Print tags to build
	fmt.Println("\nTags to promote:")
	for _, info := range imageInfo {
		fmt.Println(strings.Join(strings.Split(info.Tags, "\\n"), ","))
	}

	// Build workflow dispatch commands
	var commands []string
	commands = append(commands, "#!/bin/sh")
	for _, info := range imageInfo {
		commands = append(commands, BuildCommand(DockerBuildWorkflow, info))
	}
	command := strings.Join(commands, "\n")

	// Create scripts folder
	err = os.MkdirAll("scripts", 0700)
	if err != nil {
		panic(err)
	}

	if dryRun {
		// Create empty scripts/discover.sh
		err := os.WriteFile("scripts/dispatch.sh", []byte("#!/bin/sh\n"), 0700)
		if err != nil {
			panic(err)
		}

		// print scripts content
		fmt.Println("\nThis is a dry run, so we generate the following script but not write to scripts/dispatch.sh")
		fmt.Println("\n" + command)
	} else {
		// Write to scripts/dispatch.sh
		err = os.WriteFile("scripts/dispatch.sh", []byte(command), 0700)
		if err != nil {
			panic(err)
		}

		fmt.Println("\nShell script has been generated to scripts/dispatch.sh")
	}
}

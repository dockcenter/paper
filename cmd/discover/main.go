package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dockcenter/paper/internal/app/discover"
	"github.com/go-resty/resty/v2"
	"os"
	"os/exec"
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
		fmt.Println(strings.Join(strings.Split(info.Tags, "\n"), ","))
	}

	// Build workflow dispatch commands
	var commands []string
	for _, info := range imageInfo {
		commands = append(commands, BuildCommand(DockerBuildWorkflow, info))
	}
	command := strings.Join(commands, "\n")

	if dryRun {
		// print scripts content
		fmt.Println("\nThis is a dry run, so we generate the following script but don't run it")
		fmt.Println("\n" + command)
	} else {
		// Run command
		output, err := exec.Command("bash", "-c", command).CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			panic(err)
		}
		fmt.Println(output)
	}
}

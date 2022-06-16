package discover

import (
	"fmt"
	"github.com/dockcenter/paper/internal/pkg/utils/slices"
	"golang.org/x/mod/semver"
	"sort"
	"strings"
)

type Event int

const (
	Rebuild Event = iota
	Cron
)

type ImageInfo struct {
	DownloadURL string `json:"downloadURL"`
	Tags        string `json:"tags"`
}

func BuildImageInfo(builds []VersionFamilyBuild, existingTags []string, event Event) []ImageInfo {
	var notExistingTags []int
	sharedTags := make(map[string]int)

	// Mark builds
	for i, build := range builds {
		// Check if tag is existed
		if !slices.Contains(existingTags, GetUniqueTag(build.Version, build.Build)) {
			notExistingTags = append(notExistingTags, i)
		}

		// handle latest tag
		if i == len(builds)-1 {
			sharedTags["latest"] = i
		}

		// handle semver tags
		sharedTags[semver.Major("v" + build.Version)[1:]] = i
		sharedTags[semver.MajorMinor("v" + build.Version)[1:]] = i
		sharedTags[semver.Canonical("v" + build.Version)[1:]] = i
	}

	// Filter builds
	var returnedIndexes []int
	if event == Rebuild {
		for _, index := range sharedTags {
			if !slices.Contains(returnedIndexes, index) {
				returnedIndexes = append(returnedIndexes, index)
			}
		}
	} else {
		returnedIndexes = append(returnedIndexes, notExistingTags...)
	}

	sort.Ints(returnedIndexes)

	// Build image info
	var imageInfo []ImageInfo
	for _, index := range returnedIndexes {
		build := builds[index]

		// Build docker tags
		var tags []string
		tags = append(tags, GetUniqueTag(build.Version, build.Build))
		for k, v := range sharedTags {
			if v == index {
				tags = append(tags, k)
			}
		}

		info := ImageInfo{
			DownloadURL: fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", Project, build.Version, build.Build, build.Downloads[DownloadsKey].Name),
			Tags:        strings.Join(tags, "\\n"),
		}
		imageInfo = append(imageInfo, info)
	}

	return imageInfo
}

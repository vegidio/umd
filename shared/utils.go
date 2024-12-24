package shared

import (
	"crypto/sha1"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"path/filepath"
	"strconv"
	"strings"
)

var fs = afero.NewOsFs()

type Tag struct {
	Name string `json:"name"`
}

func IsOutdated(currentVersion string, repo string) bool {
	var tags []Tag
	var url = fmt.Sprintf("https://api.github.com/repos/%s/tags", repo)

	client := resty.New()
	resp, err := client.R().SetResult(&tags).Get(url)
	if err != nil || resp.StatusCode() != 200 {
		return true
	}

	return currentVersion != tags[0].Name
}

func CreateReport(directory string, downloads []Download) {
	filePath := filepath.Join(directory, "_report.md")
	file, err := fs.Create(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	// Filter the failed downloads
	failedDownloads := lo.Filter(downloads, func(d Download, _ int) bool {
		return !d.IsSuccess
	})

	fileContent := "# UMD - Download Report\n\n"
	fileContent += "## Failed Downloads\n\n"
	fileContent += fmt.Sprintf("- Total: %d\n", len(failedDownloads))

	for _, download := range failedDownloads {
		fileContent += fmt.Sprintf("### 🔗 Link: %s - ❌ **Failure**\n", download.Url)
		fileContent += "### 📝 Error:\n"
		fileContent += "```\n"
		fileContent += fmt.Sprintf("%s\n", download.Error)
		fileContent += "```\n"
		fileContent += "---\n"
	}

	if len(failedDownloads) > 0 {
		fileContent += createManualDownloadCommand(failedDownloads)
	}

	_, _ = file.WriteString(fileContent)
}

func createManualDownloadCommand(downloads []Download) string {
	fileContent := "\n## Retry Failed Downloads\n\n"
	fileContent += "You can retry the failed downloads by using either [aria2](https://aria2.github.io) (recommended) or [wget](https://www.gnu.org/software/wget):\n\n"
	fileContent += "### Aria2\n\n"
	fileContent += "```bash\n"

	downloadList := lo.Reduce(downloads, func(acc string, d Download, _ int) string {
		return acc + fmt.Sprintf(" %s", d.Url)
	}, "$ aria2c --file-allocation=none --auto-file-renaming=false --always-resume=true --conditional-get=true -c -s 1 -x 5 -j 5 -m 10 -Z")

	line := ""
	for _, part := range strings.Split(downloadList, " ") {
		if (len(line) + len(part)) >= 118 {
			fileContent += line + " \\\n"
			line = "   "
		}

		line += " " + part
	}

	fileContent += line + "\n"
	fileContent += "```\n"

	return fileContent
}

func RemoveDuplicates(downloads []Download, onDuplicateDeleted func(download Download)) (int, []Download) {
	numDeleted := 0
	remaining := make([]Download, 0)

	duplicates := lo.GroupBy(downloads, func(d Download) string {
		return d.Hash
	})

	for _, value := range duplicates {
		remaining = append(remaining, value[0])
		deleteList := value[1:]

		for _, deleteFile := range deleteList {
			numDeleted++
			_ = fs.Remove(deleteFile.FilePath)

			if onDuplicateDeleted != nil {
				onDuplicateDeleted(deleteFile)
			}
		}
	}

	return numDeleted, remaining
}

func CreateTimestamp(num int64) string {
	base36 := strconv.FormatInt(num, 36)
	return fmt.Sprintf("%06s", base36)
}

func CreateHashSuffix(str string) string {
	hash := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)[:4]
}

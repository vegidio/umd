package shared

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dromara/dongle"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/vegidio/go-sak/fetch"
	"github.com/zeebo/blake3"
)

var fs = afero.NewOsFs()

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

func RemoveDuplicates(downloads []Download, onDuplicateDeleted func(download Download)) (int, []Download) {
	numDeleted := 0
	remaining := make([]Download, 0)

	// Filter out downloads with empty hash before grouping
	validDownloads := lo.Filter(downloads, func(d Download, _ int) bool {
		return d.Hash != ""
	})

	duplicates := lo.GroupBy(validDownloads, func(d Download) string {
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
	hash := blake3.Sum256([]byte(str))
	return dongle.Encode.FromBytes(hash[:]).ByBase62().String()[:4]
}

func GetMediaType(filePath string) string {
	lowerExt := strings.TrimPrefix(filepath.Ext(filePath), ".")

	switch lowerExt {
	case "avif", "gif", "jpg", "jpeg", "png", "webp":
		return "image"
	case "gifv", "m4v", "mkv", "mov", "mp4", "webm":
		return "video"
	default:
		return "unkwn"
	}
}

// GetCookies retrieves cookies based on the specified type and location.
//
// The cookiesType parameter determines the source of cookies:
//   - "automatic": retrieves cookies from the browser.
//   - "manual": loads cookies from a file.
//
// # Parameters:
//   - cookiesType: the method to retrieve cookies ("automatic" or "manual")
//   - location: the domain associated with the cookies (when cookiesType is "automatic") or the path to the cookie file
//     (when cookiesType is "manual")
//
// # Returns:
//   - []fetch.Cookie: a slice of cookies retrieved from the specified source
//   - error: an error if cookies are required but none were found, or if file loading fails
//
// # Example:
//
//	cookies, err := GetCookies("automatic", "simpcity.cr")
//	if err != nil {
//	    log.Fatal(err)
//	}
func GetCookies(cookiesType, location string) ([]fetch.Cookie, error) {
	cookies := make([]fetch.Cookie, 0)

	switch cookiesType {
	case "automatic":
		cookies = fetch.GetBrowserCookies(location)
	case "manual":
		if co, err := fetch.GetFileCookies(location); err == nil {
			cookies = co
		}
	}

	if cookiesType != "disabled" && len(cookies) == 0 {
		return nil, fmt.Errorf("no cookies found")
	}

	return cookies, nil
}

// region - Private functions

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

// endregion

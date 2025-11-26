package charm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vegidio/umd"
)

func PrintSite(name string) {
	fmt.Printf("\n🌍 Website found: %s", green.Render(name))
}

func PrintType(name string) {
	fmt.Printf("; extractor type: %s\n", yellow.Render(name))
}

func PrintCachedResults(typeName string, name string, resp *umd.Response) {
	size := strconv.Itoa(len(resp.Media))
	fmt.Printf("💾 Using cached results for %s %s... %s %s Found\n",
		strings.ToLower(typeName), bold.Render(name), pink.Render("⣿"), bold.Render(size))
}

func PrintDeleted(fileName string) {
	fmt.Printf("[%s] Deleting file %s\n", red.Render("D"), bold.Render(fileName))
}

func PrintError(message string, a ...interface{}) {
	format := fmt.Sprintf(message, a...)
	fmt.Printf("🧨 %s\n", red.Render(format))
}

func PrintNewVersion(message, url string) {
	fmt.Printf("\n🌟 %s %s\n", yellow.Render(message), yellowUnderline.Render(url))
}

func PrintDone(message string) {
	fmt.Printf("\n🌟 %s\n", yellowBold.Render(message))
}

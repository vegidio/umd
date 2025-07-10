package main

import (
	"cli/internal/charm"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vegidio/shared"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/fetch"
	"path/filepath"
)

func startQuery(
	url string,
	directory string,
	parallel int,
	limit int,
	extensions []string,
	noCache bool,
	cookies []fetch.Cookie,
) error {
	mp := shared.NewMixPanel(uuid.New().String())
	fields := make(map[string]any)
	var resp *umd.Response
	var err error

	fields["interface"] = "cli"
	fields["limit"] = limit

	extractor, err := umd.New(nil).FindExtractor(url)
	if err != nil {
		return err
	}

	fields["extractor"] = extractor.Type().String()
	charm.PrintSite(extractor.Type().String())

	source, err := extractor.SourceType()
	if err != nil {
		return err
	}

	fields["source"] = source.Type()
	fields["name"] = source.Name()
	charm.PrintType(source.Type())

	fullDir := filepath.Join(directory, extractor.Type().String(), source.Name())
	cachePath := filepath.Join(fullDir, "_cache.gob")

	// Load any existing cache
	if !noCache {
		resp, _ = shared.LoadCache(cachePath)

		if resp != nil {
			charm.PrintCachedResults(source.Type(), source.Name(), resp)
		}
	}

	fields["cache"] = resp != nil
	mp.Track("Start Download", fields)

	// nil means that nothing was found in the cache
	if resp == nil {
		resp, _ = extractor.QueryMedia(limit, extensions, true)

		err = charm.StartSpinner(source.Type(), source.Name(), resp)
		if err != nil {
			return err
		}

		_ = shared.SaveCache(cachePath, resp)
	}

	clear(fields)
	fields["parallel"] = parallel
	fields["mediaFound"] = len(resp.Media)

	result := shared.DownloadAll(resp.Media, fullDir, parallel, cookies)
	responses, err := charm.StartProgress(result, len(resp.Media))
	if err != nil {
		return err
	}

	downloads := lo.Map(responses, func(r *fetch.Response, _ int) shared.Download { return shared.ResponseToDownload(r) })
	successes := lo.CountBy(downloads, func(d shared.Download) bool { return d.IsSuccess })
	failures := lo.CountBy(downloads, func(d shared.Download) bool { return !d.IsSuccess })
	fields["numSuccesses"] = successes
	fields["numFailures"] = failures

	isFirstDuplicate := true
	_, remaining := shared.RemoveDuplicates(downloads, func(download shared.Download) {
		if isFirstDuplicate {
			fmt.Println("\n🚮 Removing duplicated downloads...")
			isFirstDuplicate = false
		}

		fileName := filepath.Base(download.FilePath)
		charm.PrintDeleted(fileName)
	})

	mp.Track("End Download", fields)
	shared.CreateReport(fullDir, remaining)

	charm.PrintDone("Done!")
	return nil
}

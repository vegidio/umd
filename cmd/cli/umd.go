package main

import (
	"cli/internal/charm"
	"fmt"
	"path/filepath"
	"shared"

	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/go-sak/o11y"
	"github.com/vegidio/umd"
)

func startQuery(
	url string,
	directory string,
	parallel int,
	limit int,
	extensions []string,
	noCache bool,
	noTelemetry bool,
	cookies []fetch.Cookie,
) error {
	tel := o11y.NewTelemetry(shared.OtelEndpoint, "umd", shared.Version, shared.OtelEnvironment, !noTelemetry)
	defer tel.Close()

	fields := make(map[string]any)
	var resp *umd.Response
	var err error

	fields["interface"] = "cli"
	fields["limit"] = limit
	fields["url"] = url

	u := umd.New()

	if len(cookies) > 0 {
		cookieHeader := fetch.CookiesToHeader(cookies)
		metadata := umd.Metadata{
			umd.Coomer:   map[string]any{"cookie": cookieHeader},
			umd.SimpCity: map[string]any{"cookie": cookieHeader},
		}

		u = u.WithMetadata(metadata)
	}

	extractor, err := u.FindExtractor(url)
	if err != nil {
		tel.LogError("Extractor not found", fields, err)
		return err
	}

	fields["extractor"] = extractor.Type().String()
	charm.PrintSite(extractor.Type().String())

	source, err := extractor.SourceType()
	if err != nil {
		fmt.Printf("\n\n")
		tel.LogError("Source type not found", fields, err)
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
	tel.LogInfo("Start download", fields)

	// nil means that nothing was found in the cache
	if resp == nil {
		resp, _ = extractor.QueryMedia(limit, extensions, true)

		err = charm.StartSpinner(source.Type(), source.Name(), resp)
		if err != nil {
			tel.LogError("Error while querying media", fields, err)
			return err
		}

		_ = shared.SaveCache(cachePath, resp)
	}

	clear(fields)
	fields["parallel"] = parallel
	fields["media.found"] = len(resp.Media)

	result := shared.DownloadAll(resp.Media, fullDir, parallel)
	responses, err := charm.StartProgress(result, len(resp.Media))
	if err != nil {
		tel.LogError("Error while downloading media", fields, err)
		return err
	}

	downloads := lo.Map(responses, func(r *fetch.Response, _ int) shared.Download { return shared.ResponseToDownload(r) })
	successes := lo.CountBy(downloads, func(d shared.Download) bool { return d.IsSuccess })
	failures := lo.CountBy(downloads, func(d shared.Download) bool { return !d.IsSuccess })
	fields["downloads.success"] = successes
	fields["downloads.failure"] = failures

	isFirstDuplicate := true
	_, remaining := shared.RemoveDuplicates(downloads, func(download shared.Download) {
		if isFirstDuplicate {
			fmt.Println("\n🚮 Removing duplicated downloads...")
			isFirstDuplicate = false
		}

		fileName := filepath.Base(download.FilePath)
		charm.PrintDeleted(fileName)
	})

	tel.LogInfo("End download", fields)

	shared.CreateReport(fullDir, remaining)

	charm.PrintDone("Done!")
	return nil
}

# Usage

With the library properly [installed](installation.md), you need to initialize **UMD** by using `umd.New` passing any metadata or callback (if needed, otherwise just pass `nil`). Then, with the newly `Umd` object you can call the method `FindExtractor` passing that URL that you want to query.

**UMD** will automatically detect what site/content you're trying to fetch media information; if the URL is not supported, then it will return an error.

If everything goes well and **UMD** detects the URL returning a suitable extractor, you can use the methods below:

## QueryMedia()

```go linenums="1"
umd := umd.New(nil)
extractor, _ := umd.FindExtractor("https://www.reddit.com/user/atomicbrunette18")
resp, _ := extractor.QueryMedia(100, nil, true)
```

The method `QueryMedia` is an async function that returns immediately after the query is started. To get the result of the query, you need to use the `Response` object.

### Response object

The response represents the result of the query. It contains many fields, the most important are:

- `Url`: The URL that was queried.
- `Media`: An array of `Media` objects founds during the query.
- `Done`: Is a channel that emits a signal when the query is finished; `nil` if it completes successfully, otherwise it will contain an error.
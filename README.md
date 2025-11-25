# Universal Media Downloader (UMD)

<p align="center">
<img src="assets/icon.avif" width="300" alt="Universal Media Downloader (UMD)"/>
<br/>
<strong>UMD</strong> is an app to easily download media files hosted on popular websites.
<br/>
It supports the following sites:
<br/><br/>
<a href="https://bunkr.cr" target="_blank"><img src="https://img.shields.io/badge/Bunkr-7C3AED?&style=for-the-badge&logo=picardsurgeles&logoColor=white"/></a>
<a href="https://coomer.st" target="_blank"><img src="https://img.shields.io/badge/Coomer-1392F4?&style=for-the-badge&logo=onlyfans&logoColor=white"/></a>
<a href="https://cyberdrop.cr" target="_blank"><img src="https://img.shields.io/badge/Cyberdrop-FE41BE?&style=for-the-badge&logo=sunrise&logoColor=white"/></a>
<a href="https://erome.mcom" target="_blank"><img src="https://img.shields.io/badge/Erome-EE5F96?&style=for-the-badge&logo=githubsponsors&logoColor=white"/></a>
<a href="https://fapello.com" target="_blank"><img src="https://img.shields.io/badge/Fapello-FF647C?&style=for-the-badge&logo=favro&logoColor=white"/></a>
<a href="https://imaglr.com" target="_blank"><img src="https://img.shields.io/badge/Imaglr-0A3257?&style=for-the-badge&logo=quizlet&logoColor=white"/></a>
<br>
<a href="https://jpg6.su" target="_blank"><img src="https://img.shields.io/badge/JpgFish-E97B1D?&style=for-the-badge&logo=fishshell&logoColor=white"/></a>
<a href="https://kemono.cr" target="_blank"><img src="https://img.shields.io/badge/Kemono-E6712F?&style=for-the-badge&logo=keystone&logoColor=white"/></a>
<a href="https://reddit.com" target="_blank"><img src="https://img.shields.io/badge/Reddit-FF4500?&style=for-the-badge&logo=reddit&logoColor=white"/></a>
<a href="https://redgifs.com" target="_blank"><img src="https://img.shields.io/badge/RedGifs-D70002?&style=for-the-badge&logo=codeigniter&logoColor=white"/></a>
<a href="https://saint2.su" target="_blank"><img src="https://img.shields.io/badge/Saint-6F1212?&style=for-the-badge&logo=codecov&logoColor=white"/></a>
<a href="https://simpcity.cr" target="_blank"><img src="https://img.shields.io/badge/SimpCity-3CB8C6?&style=for-the-badge&logo=wikiversity&logoColor=white"/></a>
</p>

## 🖼️ Usage

There are two ways to use this app: using the GUI or the CLI.

The GUI is the easiest way to use the app, with an intuitive interface that allows you to download media files with just a few clicks. The CLI is more advanced and allows you to download media files in a more automated way.

Both versions are available for Windows, macOS, and Linux. Download the [latest release](https://github.com/vegidio/umd/releases) that matches your computer architecture and operating system and follow the instructions below:

### GUI

<p align="center">
<img src="assets/gui-screenshot.avif" width="80%" alt="UMD - GUI"/>
</p>

1. `Enter a URL` of the website where the media is hosted and then click on the `Query` button. If the URL belongs to a website supported by **UMD**, it will query the website and show the media files in the list below.
    - You can limit the number of files that will be queried by entering a number in the `Limit` field.
2. Select the media files that you want to download by clicking on the checkbox next to the file name, or click on the checkboxes to automatically select all images and/or videos.
3. Click on the directory field to select the directory where you want the files to be saved and then click on the `Download` button to start the download.

### CLI

<p align="center">
<img src="assets/cli-screenshot.avif" width="80%" alt="UMD - CLI"/>
</p>

Run the command below in the terminal:

```bash
$ umd-dl -d ~/Downloads/UMD <url>
```

Where:

- `-d` (optional): the directory where you want the files to be saved; default is the current directory.
- `-l` (optional): the maximum number of files to download; default is 99.999 files.
- `<url>` (mandatory): the URL of the website where the media is hosted (**Attention:** it's important that the URL is the last parameter).

For the full list of parameters, type `umd-dl --help` in the terminal.

## 💡 Features

### Cached results

When you query the same URL multiple times, the app caches the results to speed up the process, **as long as the files are saved in the same directory.**

You can adjust the app settings to bypass the cache and always fetch fresh results.

### Resume/skip previously downloaded files

When downloading files to the same directory, the app will automatically resume the files that were partially downloaded and skip those have already been fully downloaded.

This feature helps to avoid duplicate downloads, ensuring only new or previously failed files are downloaded.

### Remove duplicates

This application will automatically delete all files that are identical.

## 💣 Troubleshooting

### Download is not working or taking too long

Some websites, like Coomer/Kemono, use techniques to block tools from scraping their content. This includes slowing down the download speed or **even blocking the IP address of the user.**

Unfortunately, there’s nothing **UMD** can do to improve this, so please be patient and allow the app to complete its work.

### "App Is Damaged/Blocked..." (Windows & macOS only)

For a couple of years now, Microsoft and Apple have required developers to join their "Developer Program" to gain the pretentious status of an _identified developer_ 😛.

Translating to non-BS language, this means that if you’re not registered with them (i.e., paying the fee), you can’t freely distribute Windows or macOS software. Apps from unidentified developers will display a message saying the app is damaged or blocked and can’t be opened.

To bypass this, open the Terminal and run one of the commands below (depending on your operating system), replacing `<path-to-app>` with the correct path to where you’ve installed the app:

- Windows: `Unblock-File -Path <path-to-app>`
- macOS: `xattr -d com.apple.quarantine <path-to-app>`

## 🛠️ Build

### Dependencies

In order to build this project you will need the following dependencies installed in your computer:

- [Golang](https://go.dev/doc/install)
- [Task](https://taskfile.dev/installation/)

If you want to build the GUI you will also need:

- [Node.js](https://nodejs.org/en/download/)
- [PNPM](https://pnpm.io/installation)
- [Wails 2+](https://wails.io/docs/gettingstarted/installation)

### Compiling

With all the dependencies installed, in the project's root folder run the command:

```bash
$ task <interface> os=<operating-system> arch=<architecture>
```

Where:

- `<interface>`: can be `cli` or `gui`.
- `<operating-system>`: can be `windows`, `darwin` (macOS), or `linux`.
- `<architecture>`: can be `amd64` or `arm64`.

For example, if I wanted to build a GUI version of the app for Windows, on architecture AMD64, I would run the command:

```bash
$ task gui os=windows arch=amd64
```

## 📈 Telemetry

This app collects information about the data that you're downloading to help me track bugs and improve the general stability of the software.

**No identifiable information about you or your computer is tracked.** But if you still want to stop the telemetry, you can do that by adding the flag `--no-telemetry` in the CLI tool.

## 📝 License

**UMD** is released under the MIT License. See [LICENSE](LICENSE) for details.

## 👨🏾‍💻 Author

Vinicius Egidio ([vinicius.io](http://vinicius.io))


# Go-linkchecker

go-linkchecker is a simple command-line utility implemented in Go, it's designed to check a list of URLs (in Markdown format) for broken links. 

## Features

- Check multiple files
- Allows SSL errors
- Allows redirects
- Check for connection timeouts
- Configurable request delay
- Configurable connection timeout
- Skip saving results
- White list specific URLs

## Installation

This section gives a walkthrough on how to install and setup the go-linkchecker on your local system. Please follow these steps:

1. **Pre-requisites**
   - You need to have [Go](https://golang.org/) programming language installed on your system to run this utility. 

2. **Install the project**
   - Install the project using go command which generates an executable file.

        ```sh
        go install github.com/skrashevich/go-linkchecker@latest
        ```

After you've built the project, you can start using the utility by following the usage instructions. If you come across any errors during the installation, please check the pre-requisites again or create an issue on the GitHub page of the repository.


## Usage

Download or clone the repository, navigate to the project folder in your terminal and build the Go project with the `go build` command, which will create the executable file that you can run.

Here is a basic usage guide:

```sh
go-linkchecker -f file1.md,file2.md
```

In the above instance, go-linkchecker would check all URLs in these files `file1.md` and `file2.md`.

```sh
go-linkchecker -allow-ssl -allow-redirect -allow-timeout file1.md
```

In this instance, go-linkchecker would check URLs and allow SSL errors, redirects, and timeouts.

## Configuration

The utility has a set of configurable parameters listed below:

- `-f, --files` : Comma-separated markdown files to check
- `-a, --allow` : Status code errors to allow
- `-allow-ssl` : Allow SSL errors
- `-allow-redirect` : Allow redirected URLs
- `-allow-timeout` : Allow URLs that time out
- `-base-url` : Base URL to use for relative links
- `-d, --request-delay` : Set request delay in seconds
- `-t, --set-timeout` : Set connection timeout in seconds
- `-skip-save-results` : Skip saving results
- `-w, --white-list` : Comma separated URLs to white list


Please note this utility only works with markdown files and `http` or `https` URLs.
For more information about specific flags, please refer to the source code.


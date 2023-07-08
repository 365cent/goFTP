# goFTP
`goFTP` is a web-based application written in Go that serves as a file listing and download tool for remote FTP servers. It enables the listing of files and folders, direct download of files, as well as the ability to remotely download files from one FTP server to another. This tool is designed to be easy to set up and use, and it can be used as a one-time script or run as a background service.

## Screenshot
![Screenshot 2023-07-08 at 5.39.39 PM.png](Screenshot 2023-07-08 at 5.39.39 PM.png)

## Features
- List files and folders in a remote FTP server.
- Download files directly from the FTP server.
- Remotely download a file from a URL and upload it to the FTP server.
- Stream media files directly in the browser.
- The server also serves static files, which can be useful for CSS, JavaScript, or other client-side resources.

## Usage
### One-time Execution
To use `goFTP` as a one-time script, you should replace `ftpAddress`, `ftpUsername`, and `ftpPassword` variables in the `main.go` file with your own FTP server credentials and then execute the script using `go run main.go`.

### Run as a Service
For long-term usage, you can compile `goFTP` into a binary with `go build -o start main.go` and set it up as a service to run in the background. This can be achieved by using process managers like `supervisord`.

Here's an example of a `supervisord` configuration:
```
[program:ftp]
directory=/var/www/ftp
command=/var/www/ftp/start
autostart=true
autorestart=true
startretries=3
user=root
redirect_stderr=true
stdout_logfile=/var/www/ftp/stdout.log
stderr_logfile=/var/www/ftp/stderr.log
```
In this configuration, replace `/var/www/ftp` with the directory where you placed the `start` binary. This will run `goFTP` as a service, automatically starting it when the system boots or if it crashes.

## Dependencies
`goFTP` uses the following dependencies:
- The `net/http` package from the standard library to handle HTTP requests and responses.
- The `html/template` package from the standard library to render HTML templates.
- The `github.com/jlaffaye/ftp` package for FTP connections.

Please ensure you've installed these dependencies in your Go environment.

For more details, please review the source code and the inline comments.

## Logging
By default, `stdout` and `stderr` are logged to `stdout.log` and `stderr.log` respectively in the directory specified in the `supervisord` configuration. 

This provides an easy way to monitor the output and errors of `goFTP`. For example, you can use `tail -f /var/www/ftp/stdout.log` to follow the standard output in real-time.

## Note
Please note that the functionality of `goFTP` depends on the permissions and capabilities of the FTP server to which it connects. Make sure the provided FTP user has sufficient rights to list files, download files, and upload files.

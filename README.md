# Extended FS

This library allows accessing local and remote files and performing file operations, using the same API.

* Read a file from local file system: `stringContent, err := ReadFile("/var/www/html/index.html")`
* Read a file from the configured remote file system: `stringContent, err := ReadFile("sftp:///var/www/html/index.html")`
* Read a file from the in-binary file system (via pkger): `stringContent, err := ReadFile("pkger:///var/www/html/index.html")`

## Target

* Local: just set an absolute or relative file path, e.g. `/var/www/html`
* Remote: just add `ssh://` or `sftp://` protocol before the location, e.g. `ssh:///var/www/html` or `sftp:///var/www/html`
* Packaged (pkger): If you want to access files bundled with your binary using [pkger](https://github.com/markbates/pkger), just add `pkger://` protocol before the location, e.g. `pkger:///resources/file.jpg`
  * Note: _pkger_ only allows reading operations on files. Write and directory operations will not work!

## Examples

### Write local file

This will create a file `/var/www/html/index.html` with a HTML content on the local file system.

```golang
import "github.com/saitho/golang-extended-fs"

extended_fs.WriteFile("/var/www/html/index.html", "<h1>Hello World</h1>")
```

### Write remote file

This will create a remote file `/var/www/html/index.html` with a HTML content on the remote file system at `192.168.2.105`.
Local SSH keys are always loaded per default. Additionally, a custom key can be specified via `SshIdentify` setting.2

See [SFTP config](./sftp/config.go) for all available settings.

```golang
import (
  "github.com/saitho/golang-extended-fs"
  "github.com/saitho/golang-extended-fs/sftp"
)
sftp.Config.SshHost = "192.168.2.105"
extended_fs.WriteFile("ssh:///var/www/html/index.html", "<h1>Hello World</h1>")
```

```golang
import (
  "github.com/saitho/golang-extended-fs"
  "github.com/saitho/golang-extended-fs/sftp"
)
sftp.Config.SshHost = "192.168.2.105"
sftp.Config.LoadLocalSigners = false // do not load local SSH private keys
sftp.Config.SshIdentity = "/path/to/my/private_key.pem"
extended_fs.WriteFile("ssh:///var/www/html/index.html", "<h1>Hello World</h1>")
extended_fs.Chown("ssh:///var/www/html/index.html", 1001, 1001) // set user and group with id 1001 as owner
```

## Volume Toolkit

This is an utility to do various file based operations on a persistent volume. The toolkit is designed to work with the persistent volume mounted at `/data`.

### Commands

| Command   | Flags                                    | Description                                    |
| --------- | ---------------------------------------- | ---------------------------------------------- |
| `export`  | `--path <path>` `--http-method <method>` | Export a persistent volume to a file or URL.   |
| `import`  | `--path <path>` `--http-method <method>` | Import a persistent volume from a file or URL. |
| `size`    | None                                     | Measure the size of a persistent volume.       |
| `destroy` | None                                     | Destroy all data in the persistent volume.     |
| `file-op` | `--path <path>`                          | Perform file operations.                       |

#### File Related Commands

The `file-op` command supports the following subcommands:

| Subcommand | Flags                                       | Description                                       |
| ---------- | ------------------------------------------- | ------------------------------------------------- |
| `ls`       | `--path <path>`                             | List files in the specified directory.            |
| `cat`      | `--path <path>`                             | Display the contents of the specified file.       |
| `cp`       | `--path <src>` `<dest>`                     | Copy a file from source to destination.           |
| `mv`       | `--path <src>` `<dest>`                     | Move a file from source to destination.           |
| `rm`       | `--path <path>`                             | Remove the specified file or directory.           |
| `mkdir`    | `--path <path>`                             | Create a directory at the specified path.         |
| `chmod`    | `--path <path>` `<mode>`                    | Change the permissions of the specified file.     |
| `chown`    | `--path <path>` `--uid <uid>` `--gid <gid>` | Change the ownership of the specified file.       |
| `download` | `--path <dest>` `<url>`                     | Download a file from a URL to the specified path. |

**Note :** `cat` command will stream the binary data of the file to stdout.

### Output format

```json
{
  "success": true,
  "error": "",
  "data": "Data has been imported"
}
```

### Notes

- `/data` > This is the mounting point for the persistent volume
- `/app` > This is the mounting point for the toolkit. After export volume, the `backup.tar.gz` file will be created in this directory. And for import volume, the `backup.tar.gz` file should be placed in this directory.
- All the sizes will be in bytes format.
- File Metadata
  ```json
  {
    "name": "100MB.bin",
    "size": 104857600,
    "mode": "777",
    "mod_time": "2020-05-20T11:14:26Z",
    "uid": 999,
    "gid": 999,
    "is_dir": false
  }
  ```

### License

This project is licensed under Apache 2.0 License - see the LICENSE file for details.

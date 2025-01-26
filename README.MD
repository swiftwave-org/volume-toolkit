### Volume Toolkit

This is an utility to do various operation on docker volumes.

#### Commands

| Command   | Flags                | Description                                    |
| --------- | -------------------- | ---------------------------------------------- |
| `export`  | `path` `http-method` | Export a persistent volume to a file or URL.   |
| `import`  | `path` `http-method` | Import a persistent volume from a file or URL. |
| `size`    |                      | Measure the size of a persistent volume.       |
| `destroy` |                      | Destroy all data in the persistent volume.     |

### Output format

```json
{
  "success": true,
  "error": "",
  "data": "Data has been imported"
}
```

#### Notes

- `/data` > This is the mounting point for the persistent volume
- `/app` > This is the mounting point for the toolkit. After export volume, the `backup.tar.gz` file will be created in this directory. And for import volume, the `backup.tar.gz` file should be placed in this directory.
- All the sizes will be in bytes format.

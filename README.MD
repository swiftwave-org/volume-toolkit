> Note: This image is solely for internal use. So we are not accepting outside contributions except **Security Vulnerabilities** . 

### Volume Toolkit

**This image has only two purposes**
1. Export a persistent named volume to a *.tar.gz file.
2. Import a *.tar.gz file to a persistent named volume.

#### Mounting Points

| Mounting Point | Description |
| --- | --- |
| /data | This is the mounting point for the persistent volume. |
| /app | This is the mounting point for the toolkit. After export volume, the `backup.tar.gz` file will be created in this directory. And for import volume, the `backup.tar.gz` file should be placed in this directory. |

#### Usage

1. Export a persistent volume

  a. Mount and export
   ```bash
   docker run --rm -v <persistent-volume-name>:/data -v /path/to/backup:/app ghcr.io/swiftwave-org/volume-toolkit export
   ```
   > Note : It is recommended to use `tmp` for `/path/to/backup` as it will be deleted after the export process.

   b. Copy backup file
   ```bash
   cp /path/to/backup/backup.tar.gz <destination-path>
   ```
   c. Remove the path/to/backup
   ```bash
   rm -rf /path/to/backup
   ```

2. Import a persistent volume

   a. Create a tmp directory
   ```bash
   mkdir /path/to/backup
   ```
   b. Copy backup file
   ```bash
   cp <source-path>/backup.tar.gz /path/to/backup
   ```
   c. Wipe up persistent volume and import data from backup
   ```bash
   docker run --rm -v <persistent-volume-name>:/data -v /path/to/backup:/app ghcr.io/swiftwave-org/volume-toolkit import
   ```

3. Measure the size of a persistent volume

   a. Run the following command
   ```bash
   docker run --rm -v <persistent-volume-name>:/data  -v /path/to/backup:/app volume-toolkit size
   ```
   b. The size of the persistent volume will be written in `/path/to/backup/size.txt` file in bytes format.

# JPKG Format Standard Version 1

## File Header

File header is uncompressed and unencrypted. After writing the header, compression and encryption are then applied (data is compressed and then encrypted and then written).

1.  4 byte magic number "jpkg"
2.  8 byte version number (1)
3.  1 byte compression flag
4.  1 byte encryption flag
5.  2 byte FF
6. 16 byte padding

### COMPRESSION

|Value|   Description|
|-----|--------------|
|    0|No Compression|

### ENCRYPTION

|Value|   Description|
|-----|--------------|
|    0| No Encryption|

## Package Manifest

Compressed and encrypted package manifest.

1. 8 byte packaged at time (Unix Timestamp)
4. 8 byte file count (N files)
2. 8 byte metadata length (M bytes)
3. M byte metadata (encoded as a json string) 

This section's length is then padded to the next multiple of 16 bytes.

5. 16 byte padding
6.  N File Records (See next section)

### File Records

The absolute offset begins at 

1. 8 byte File Name Absolute Offset
2. 8 byte File Name Length (N bytes)
3. 8 byte File Metadata Absolute Offset
4. 8 byte File Metadata Length (M bytes)
5. 8 byte File Data Absolute Offset
6. 8 byte File Data Length (O bytes)

## Package Body

1. N byte File Name
2. M byte File Metadata (encoded as json string)
3. O byte File Data

## Package Footer

- SHA256 hash of file prior to compression and encryption
- Hash signed by private key (optional)

## Additional

### Padding

Padding length is calculated using the following formula, where M is the current offset / offset from the closest multiple of 16.


$$
16(\lfloor\frac{M}{16}\rfloor+1)-M\mod 16
$$
# JPKG Format Standard Version 1

## File Header

File header is uncompressed and unencrypted. After writing the header, compression and encryption are then applied (data is compressed and then encrypted and then written).

| Size (Bytes)|           Description|  Extra|
|-------------|----------------------|-------|
|            4|          Magic Number| "jpkg"|
|            8|        Version Number|      1|
|            1|      Compression Flag|      K|
|            1|       Encryption Flag|      E|
|            1|             Hash Flag|      H|
|            1| Crypto Signature Flag|      C|
|           16|               Padding|       |

### Compression Flag (K)

|Value|   Description|
|-----|--------------|
|    0|No Compression|

### Encryption Flag (E)

|Value|   Description|
|-----|--------------|
|    0| No Encryption|

### Hash Flag (H)

|Value|   Description|
|-----|--------------|
|    0|    No Hashing|

### Cryptographic Signature Flag (C)

|Value|   Description|
|-----|--------------|
|    0|  No Signature|

## Package Manifest

Compressed and encrypted package manifest.

| Size (Bytes)|                           Description|       Extra|
|-------------|--------------------------------------|------------|
|            8|               Unix Timestamp, seconds|            |
|            8|                            File Count|           M|
|            8|                   Package Name Length|           P|
|            8|               Package Metadata Length|           N|
|            P|                          Package Name|       ascii|
|            N|                      Package Metadata|json,  ascii|
|            -|        Padding to next multiple of 16|            |
|           16|                               Padding|            |

## Package Body

1.  M File Records (See below)
2. 16 byte padding

### File Records

| Size (Bytes)|                           Description|       Extra|
|-------------|--------------------------------------|------------|
|           16|                               UUID v4|        UUID|
|            8|                      Next File Offset|    Absolute|
|            8|                File Identifier Length|          FI|
|            8|                      File Path Length|          FP|
|            8|                  File Metadata Length|          FM|
|            8|             File Compressed Data Size|          CD|
|            8|           File Uncompressed Data Size|          UD|
|           FI|                       File Identifier|       ascii|
|           FP|                             File Path|       ascii|
|           FM|                         File Metadata| json, ascii|
|           CD|                  File Compressed Data|            |
|            -|        Padding to next multiple of 16|            |


## Package Footer

| Size (Bytes)|                           Description|          Extra|
|-------------|--------------------------------------|---------------|
|            -|                    Optional File Hash| Dependent on H|
|            -|      Optional Cryptographic Signature| Dependent on C|

## Additional

### Padding Algorithm

Padding length is calculated using the following formula, where M is the current offset / offset from the closest multiple of 16.

$$
16(\lfloor\frac{M}{16}\rfloor+1)-M\mod 16
$$

When padding cycle through the following bytes.

[0xDE, 0xAD, 0xBE, 0xEF]
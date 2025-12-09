#pragma endian big
#pragma magic [6A 70 6B 67] @ 0x00

import type.guid;
import type.magic;
import type.time;
import std.array;
import std.math;

fn PadTo16(){
    auto a = $/16;
    a = std::math::floor(a);
    auto b = ((16 * (a+1))-$)%16;
    return b;
};

struct Str<auto Size> {
    u8 chars[Size];
} [[sealed]];

struct Padding {
    u8 b[16];  
} [[sealed]];

struct Pad16 {
    u8 b[PadTo16()];
};

struct PackageHeader {
    type::Magic<"jpkg"> magic;
    u64 Version;
    u8 Compression, Encryption, Hash, Crypto;
    u128 Padding;
};

struct PackageManifest {
    type::time64_t PackageTime;
    u64 FileCount, PkgNameLength, PkgMetadataLength;
    Str<PkgNameLength> PkgName;
    Str<PkgMetadataLength> PkgMetadata;
    Pad16 _1;
    Padding _2;
};

struct FileRecord {
    u64 NextRecord;
    u64 IdentifierLength, PathLength;
    type::GUID uuid;
    u64 MetadataLength, CompressedSize, UncompressedSize;
    Str<IdentifierLength> Identifier;
    Str<PathLength> Path;
    Str<MetadataLength> Metadata;
    u8 CompressedData[CompressedSize];
    Pad16 _1;
};

struct Package{
    PackageHeader Header;
    PackageManifest Manifest;
    FileRecord FileRecords[Manifest.FileCount];
};

Package package @ 0x00;
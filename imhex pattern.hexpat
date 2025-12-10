import type.time;
import type.guid;
import std.string;

#pragma endian big

struct Header {
    u8 magicNumber[4];
    u64 Version;
    u8 Compressiom, Encryption, Hasher, Signature;
};

struct Manifest {
    type::time64_t packagedAt;
    u64 FileCount;
    std::string::SizedStringBase<u64, u32> PackageName;
    std::string::SizedStringBase<u64, u32> PackageMetadata;
};

struct FileRecord {
    std::string::SizedStringBase<u64, u32> Identifier;
    std::string::SizedStringBase<u64, u32> Path;
    type::GUID UUID;
    std::string::SizedStringBase<u64, u32> Metadata;
    u64 CompressedSize, UncompressedSize;
    u8 Data[CompressedSize];
};

struct Package {
    Header header;
    Manifest manifest;    
    FileRecord Records[manifest.FileCount];
};

Package package @ 0;
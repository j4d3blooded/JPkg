build-packer:
	go build -o jpkg.exe ./app/dir_pack_unpack/main.go

build-blog-packer:
	go build -o blog_pack.exe ./app/blog_pack/main.go
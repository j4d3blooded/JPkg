package jpkg

type fsTree struct {
	root *JPkgFSDirectory
}

type fsNode interface {
	isFsNode()
	GetPath() string
}

type JPkgFSDirectory struct {
	pkg      *JPkgDecoder
	parent   *JPkgFSDirectory
	children []*fsNode
	path     string
}

func (fd *JPkgFSDirectory) IsRoot() bool {
	return fd.parent == nil
}

func (fd *JPkgFSDirectory) GetChild(name string) fsNode {

}

type JPkgFSFile struct {
	pkg    *JPkgDecoder
	parent *JPkgFSDirectory
	path   string
}

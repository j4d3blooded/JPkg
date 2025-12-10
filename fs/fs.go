package jpkg_fs

import (
	"slices"
	"strings"
)

type JPkgFS struct {
	Root *JPkgFSDirectory
}

type JPkgFSNode interface {
	isNode()
	GetName() string
	GetParent() *JPkgFSDirectory
}

type JPkgFSDirectory struct {
	Parent   *JPkgFSDirectory
	Name     string
	Children []JPkgFSNode
}

func (j *JPkgFSDirectory) GetName() string {
	return j.Name
}

func (j *JPkgFSDirectory) GetParent() *JPkgFSDirectory {
	return j.Parent
}

func (j *JPkgFSDirectory) isNode() {}

type JPkgFSFile struct {
	Parent *JPkgFSDirectory
	Name   string
}

func (j *JPkgFSFile) GetName() string {
	return j.Name
}

func (j *JPkgFSFile) GetParent() *JPkgFSDirectory {
	return j.Parent
}

func (j *JPkgFSFile) isNode() {}

func GetFullPath(node JPkgFSNode) string {
	segments := []string{node.GetName()}
	a := node.GetParent()
	for {
		segments = append(segments, a.Name)
		if a == nil {
			break
		}
	}
	slices.Reverse(segments)
	return strings.Join(segments, "//")
}

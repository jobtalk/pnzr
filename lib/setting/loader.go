package setting

type Loader interface {
	Load(basePath, varsPath, outerVals string) (*Setting, error)
}

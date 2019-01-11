package gen

type includeFiles map[string]struct{}

func newIncludeFiles() includeFiles {
	return includeFiles(map[string]struct{}{})
}

func (inf includeFiles) addIncludeFile(includeFile string) {
	inf[includeFile] = struct{}{}
}

func (inf includeFiles) unifyIncludeFiles(m includeFiles) {
	for i := range m {
		inf[i] = struct{}{}
	}
}

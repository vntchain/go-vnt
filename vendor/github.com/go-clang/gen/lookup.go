package gen

type Lookup struct {
	lookupEnum        map[string]*Enum
	lookupNonTypedefs map[string]string
	lookupStruct      map[string]*Struct
}

func NewLookup() Lookup {
	return Lookup{
		lookupEnum:        map[string]*Enum{},
		lookupNonTypedefs: map[string]string{},
		lookupStruct: map[string]*Struct{
			"cxstring": &Struct{
				Name:  "cxstring",
				CName: "CXString",
			},
		},
	}
}

func (l *Lookup) RegisterEnum(e *Enum) {
	if _, ok := l.lookupEnum[e.Name]; !ok {
		l.lookupEnum[e.Name] = e
		l.lookupNonTypedefs["enum "+e.CName] = e.Name
		l.lookupEnum[e.CName] = e
	}
}

func (l *Lookup) HasEnum(n string) (*Enum, bool) {
	e, ok := l.lookupEnum[n]

	return e, ok
}

func (l *Lookup) RegisterStruct(s *Struct) {
	if _, ok := l.lookupStruct[s.Name]; !ok {
		l.lookupStruct[s.Name] = s
		l.lookupNonTypedefs["struct "+s.CName] = s.Name
		l.lookupStruct[s.CName] = s
	}
}

func (l *Lookup) HasStruct(n string) (*Struct, bool) {
	s, ok := l.lookupStruct[n]

	return s, ok
}

func (l *Lookup) RemoveStruct(s *Struct) {
	delete(l.lookupStruct, s.Name)
	delete(l.lookupNonTypedefs, "struct "+s.CName)
	delete(l.lookupStruct, s.CName)
}

func (l *Lookup) LookupNonTypedef(s string) (string, bool) {
	n, ok := l.lookupNonTypedefs[s]

	return n, ok
}

func (l *Lookup) IsEnumOrStruct(name string) bool {
	if _, ok := l.HasEnum(name); ok {
		return true
	} else if _, ok := l.HasStruct(name); ok {
		return true
	}

	return false
}

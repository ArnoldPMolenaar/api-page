package responses

type TypeLookupList struct {
	Types []string `json:"types"`
}

// SetTypeLookupList sets the list of type names.
func (tll *TypeLookupList) SetTypeLookupList(types []string) {
	tll.Types = make([]string, len(types))
	copy(tll.Types, types)
}

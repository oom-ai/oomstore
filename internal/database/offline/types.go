package offline

type GetFeatureValuesStreamOpt struct {
	DataTable    string
	EntityName   string
	FeatureNames []string
	Limit        *uint64
}

package types

type OneStoreOpt struct {
	Host      string
	Port      string
	User      string
	Pass      string
	Workspace string
}

type ListFeatureOpt struct {
	EntityName *string
	GroupName  *string
}

type UpdateFeatureOpt struct {
	FeatureName    string
	NewDescription string
}

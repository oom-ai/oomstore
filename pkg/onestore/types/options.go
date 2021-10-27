package types

type OneStoreOpt struct {
	Host      string
	Port      string
	User      string
	Pass      string
	Workspace string
}

func (opt *OneStoreOpt) ToOneStoreOptV2() *OneStoreOptV2 {
	if opt == nil {
		return nil
	}

	postgresOpt := PostgresDbOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Pass:     opt.Pass,
		Database: opt.Workspace,
	}
	return &OneStoreOptV2{
		MetaStoreOpt:    MetaStoreOpt{PostgresDbOpt: &postgresOpt, Backend: POSTGRES},
		OnlineStoreOpt:  OnlineStoreOpt{PostgresDbOpt: &postgresOpt, Backend: POSTGRES},
		OfflineStoreOpt: OfflineStoreOpt{PostgresDbOpt: &postgresOpt, Backend: POSTGRES},
	}
}

type CreateFeatureOpt struct {
	FeatureName string
	GroupName   string
	DBValueType string
	Description string
}

type ListFeatureOpt struct {
	EntityName *string
	GroupName  *string
}

type UpdateFeatureOpt struct {
	FeatureName    string
	NewDescription string
}

type CreateEntityOpt struct {
	Name        string
	Length      int
	Description string
}

type CreateFeatureGroupOpt struct {
	Name        string
	EntityName  string
	Description string
}

type ExportFeatureValuesOpt struct {
	GroupName     string
	GroupRevision *int64
	FeatureNames  []string
	Limit         *uint64
}

type ImportBatchFeaturesOpt struct {
	GroupName   string
	Description string
	DataSource  LocalFileDataSourceOpt
}

type LocalFileDataSourceOpt struct {
	FilePath  string
	Delimiter string
}

type GetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKey    string
}

type MultiGetOnlineFeatureValuesOpt struct {
	FeatureNames []string
	EntityKeys   []string
}

type GetHistoricalFeatureValuesOpt struct {
	FeatureNames []string
	EntityRows   []EntityRow
}

type UpdateEntityOpt struct {
	EntityName     string
	NewDescription string
}

type UpdateFeatureGroupOpt struct {
	GroupName   string
	Description *string
	Revision    *int64
	DataTable   *string
}

type MaterializeOpt struct {
	GroupName     string
	GroupRevision *int64
}

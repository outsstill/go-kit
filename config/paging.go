package config

type PagingConfig struct {
	PerPage         int    `mapstructure:"perpage" yaml:"perpage"`
	UrlQueryOrder   string `mapstructure:"url_query_order" yaml:"url_query_order"`
	UrlQuerySort    string `mapstructure:"url_query_sort" yaml:"url_query_sort"`
	UrlQueryPage    string `mapstructure:"url_query_page" yaml:"url_query_page"`
	UrlQueryPerPage string `mapstructure:"url_query_per_page" yaml:"url_query_per_page"`
}

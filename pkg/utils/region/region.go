package region

type Regions map[string]Description

type Description struct {
	RegionNameAbbreviation string `json:"region_name_abbreviation"`
	RegionName             string `json:"region_name"`
}

// GetRegionsRelationshipComparison 获取区域关系对照信息
func GetRegionsRelationshipComparison() Regions {
	return Regions{
		"af-south-1": Description{
			RegionName:             "Africa (Cape Town)",
			RegionNameAbbreviation: "",
		},
		"ap-east-1": Description{
			RegionName:             "Asia Pacific (Hong Kong)",
			RegionNameAbbreviation: "CN-HK",
		},
		"ap-northeast-1": Description{
			RegionName:             "Asia Pacific (Tokyo)",
			RegionNameAbbreviation: "JP-TK",
		},
		"ap-northeast-2": Description{
			RegionName:             "Asia Pacific (Seoul)",
			RegionNameAbbreviation: "KR-S",
		},
		"ap-northeast-3": Description{
			RegionName:             "Asia Pacific (Osaka)",
			RegionNameAbbreviation: "JP-OSA",
		},
		"ap-south-1": Description{
			RegionName:             "Asia Pacific (Mumbai)",
			RegionNameAbbreviation: "IN-BOM",
		},
		"ap-south-2": Description{
			RegionName:             "Asia Pacific (Hyderabad)",
			RegionNameAbbreviation: "",
		},
		"ap-southeast-1": Description{
			RegionName:             "Asia Pacific (Singapore)",
			RegionNameAbbreviation: "SGP",
		},
		"ap-southeast-2": Description{
			RegionName:             "Asia Pacific (Sydney)",
			RegionNameAbbreviation: "AUS-SA",
		},
		"ap-southeast-3": Description{
			RegionName:             "Asia Pacific (Jakarta)",
			RegionNameAbbreviation: "JKT",
		},
		"ca-central-1": Description{
			RegionName:             "Canada (Central)",
			RegionNameAbbreviation: "",
		},
		"eu-central-1": Description{
			RegionName:             "Europe (Frankfurt)",
			RegionNameAbbreviation: "",
		},
		"eu-central-2": Description{
			RegionName:             "Europe (Zurich)",
			RegionNameAbbreviation: "",
		},
		"eu-north-1": Description{
			RegionName:             "Europe (Stockholm)",
			RegionNameAbbreviation: "",
		},
		"eu-south-1": Description{
			RegionName:             "Europe (Milan)",
			RegionNameAbbreviation: "",
		},
		"eu-south-2": Description{
			RegionName:             "Europe (Spain)",
			RegionNameAbbreviation: "",
		},
		"eu-west-1": Description{
			RegionName:             "Europe (Ireland)",
			RegionNameAbbreviation: "",
		},
		"eu-west-2": Description{
			RegionName:             "Europe (London)",
			RegionNameAbbreviation: "UK-LND",
		},
		"eu-west-3": Description{
			RegionName:             "Europe (Paris)",
			RegionNameAbbreviation: "FR-PAR",
		},
		"me-central-1": Description{
			RegionName:             "Middle East (UAE)",
			RegionNameAbbreviation: "",
		},
		"me-south-1": Description{
			RegionName:             "Middle East (Bahrain)",
			RegionNameAbbreviation: "",
		},
		"sa-east-1": Description{
			RegionName:             "South America (Sao Paulo)",
			RegionNameAbbreviation: "",
		},
		"us-east-1": Description{
			RegionName:             "US East (N. Virginia)",
			RegionNameAbbreviation: "",
		},
		"us-east-2": Description{
			RegionName:             "US East (Ohio)",
			RegionNameAbbreviation: "",
		},
		"us-west-1": Description{
			RegionName:             "US West (N. California)",
			RegionNameAbbreviation: "",
		},
		"us-west-2": Description{
			RegionName:             "US West (Oregon)",
			RegionNameAbbreviation: "",
		},
	}
}

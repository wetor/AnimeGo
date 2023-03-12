package plugin

type FunctionOptions struct {
	// Name
	//  @Description: 函数名
	Name string
	// ParamsSchema
	//  @Description: 参数校验
	//  @Description: 必选name参数: "{name}", 可选name参数: "{name},optional"
	ParamsSchema []string
	// ResultSchema
	//  @Description: 返回值校验，同参数校验
	ResultSchema []string
	// SkipSchemaCheck
	//  @Description: 跳过参数校验
	SkipSchemaCheck bool
}

type VariableOptions struct {
	// Name
	//  @Description: 变量名
	Name string
	// Nullable
	//  @Description: 可空
	Nullable bool
}

type LoadOptions struct {
	File      string
	Functions []*FunctionOptions
	Variables []*VariableOptions
}

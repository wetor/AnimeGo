package plugin

type FuncSchemaOptions struct {
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
	// DefaultArgs
	//  @Description: 默认参数列表，会被同名的运行时参数覆盖
	DefaultArgs map[string]any
}

type VarSchemaOptions struct {
	// Name
	//  @Description: 变量名
	Name string
	// Nullable
	//  @Description: 可空
	Nullable bool
}

type LoadOptions struct {
	File string
	// Code
	//  @Description: 源码模式加载，优先源码载入
	Code *string
	// GlobalVars
	//  @Description: 全局变量，将覆盖脚本中同名全局变量，但不会覆盖程序内置全局变量
	GlobalVars map[string]any
	// FuncSchema
	//  @Description: 函数约束
	FuncSchema []*FuncSchemaOptions
	// VarSchema
	//  @Description: 全局变量约束
	VarSchema []*VarSchemaOptions
}

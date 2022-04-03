package constant

type BuiltinJobConstant string

var (
	COPY_JOB   = BuiltinJobConstant("copy")
	DELETE_JOB = BuiltinJobConstant("delete")
	RENAME_JOB = BuiltinJobConstant("rename")
	MOVE_JOB   = BuiltinJobConstant("move")
)

package trace

type Tracer interface {
	Trace(...interface{}) // 任意の型を任意の個数受け取れるTraceメソッドを持つ
}

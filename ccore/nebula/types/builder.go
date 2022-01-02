package types

type ValueBuilder interface {
	NVal(*NullType) ValueBuilder
	BVal(*bool) ValueBuilder
	IVal(*int64) ValueBuilder
	FVal(*float64) ValueBuilder
	SVal([]byte) ValueBuilder
	DVal(Date) ValueBuilder
	TVal(Time) ValueBuilder
	DtVal(DateTime) ValueBuilder
	VVal(Vertex) ValueBuilder
	EVal(Edge) ValueBuilder
	PVal(Path) ValueBuilder
	LVal(NList) ValueBuilder
	MVal(NMap) ValueBuilder
	UVal(NSet) ValueBuilder
	GVal(DataSet) ValueBuilder
	GgVal(Geography) ValueBuilder
	DuVal(Duration) ValueBuilder
	Emit() Value
}

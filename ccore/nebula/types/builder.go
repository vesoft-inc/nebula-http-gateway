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

type DateBuilder interface {
	Year(int16) DateBuilder
	Month(int8) DateBuilder
	Day(int8) DateBuilder
	Emit() Date
}

type TimeBuilder interface {
	Hour(int8) TimeBuilder
	Minute(int8) TimeBuilder
	Sec(int8) TimeBuilder
	Microsec(int32) TimeBuilder
	Emit() Time
}

type DateTimeBuilder interface {
	Year(int16) DateTimeBuilder
	Month(int8) DateTimeBuilder
	Day(int8) DateTimeBuilder
	Hour(int8) DateTimeBuilder
	Minute(int8) DateTimeBuilder
	Sec(int8) DateTimeBuilder
	Microsec(int32) DateTimeBuilder
	Emit() DateTime
}

type EdgeBuilder interface {
	Src(Value) EdgeBuilder
	Dst(Value) EdgeBuilder
	Type(EdgeType) EdgeBuilder
	Name([]byte) EdgeBuilder
	Ranking(EdgeRanking) EdgeBuilder
	Props(map[string]Value) EdgeBuilder
	Emit() Edge
}

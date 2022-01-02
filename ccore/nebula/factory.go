package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"sync"
)

type (
	Factory interface {
		//NewDataset() DataSet
		//NewRow() Row
		NewValue() types.Value
		//NewDate() Date
		//NewTime() Time
		//NewDateTime() DateTime
		//NewVertex() Vertex
		//NewEdge() Edge
		//NewPath() Path
		//NewNList() NList
		//NewNMap() NMap
		//NewNSet() NSet
		//NewGeography() Geography
		//NewTag() Tag
		//NewStep() Step
		//NewPoint() Point
		//NewLineString() LineString
		//NewPolygon() Polygon
		//NewCoordinate() Coordinate
		//NewPlanDescription() PlanDescription
		//NewPlanNodeDescription() PlanNodeDescription
		//NewPair() Pair
		//NewProfilingStats() ProfilingStats
		//NewPlanNodeBranchInfo() PlanNodeBranchInfo
	}

	defaultFactory struct {
		o              Options
		factory        types.FactoryDriver
		initDriverOnce sync.Once
	}
)

func NewFactory(opts ...Option) (Factory, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	o.complete()
	if err := o.validate(); err != nil {
		return nil, err
	}
	f := &defaultFactory{
		o:       o,
		factory: newDriverFactory(),
	}

	// init driver when create factory
	err := f.initDriver()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *defaultFactory) NewValue() types.Value {
	return f.factory.NewValue()
}

func (f *defaultFactory) initDriver() error {
	factory, err := types.GetFactoryDriver(f.o.version)
	if err != nil {
		return err
	}
	f.factory = factory
	return nil
}

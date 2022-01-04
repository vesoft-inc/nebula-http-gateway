package nebula

import (
	"sync"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	Factory interface {
		//NewDatasetBuilder() types.DataSetBuilder
		//NewRowBuilder() types.RowBuilder
		NewValueBuilder() types.ValueBuilder
		NewDateBuilder() types.DateBuilder
		NewTimeBuilder() types.TimeBuilder
		NewDateTimeBuilder() types.DateTimeBuilder
		//NewVertexBuilder() types.VertexBuilder
		NewEdgeBuilder() types.EdgeBuilder
		//NewPathBuilder() types.PathBuilder
		//NewNListBuilder() types.NListBuilder
		//NewNMapBuilder() types.NMapBuilder
		//NewNSetBuilder() types.NSetBuilder
		//NewGeographyBuilder() types.GeographyBuilder
		//NewTagBuilder() types.TagBuilder
		//NewStepBuilder() types.StepBuilder
		//NewPointBuilder() types.PointBuilder
		//NewLineStringBuilder() types.LineStringBuilder
		//NewPolygonBuilder() types.PolygonBuilder
		//NewCoordinateBuilder() types.CoordinateBuilder
		//NewPlanDescriptionBuilder() types.PlanDescriptionBuilder
		//NewPlanNodeDescriptionBuilder() types.PlanNodeDescriptionBuilder
		//NewPairBuilder() types.PairBuilder
		//NewProfilingStatsBuilder() types.ProfilingStatsBuilder
		//NewPlanNodeBranchInfoBuilder() types.PlanNodeBranchInfoBuilder
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

func (f *defaultFactory) NewValueBuilder() types.ValueBuilder {
	return f.factory.NewValueBuilder()
}

func (f *defaultFactory) NewDateBuilder() types.DateBuilder {
	return f.factory.NewDateBuilder()
}

func (f *defaultFactory) NewTimeBuilder() types.TimeBuilder {
	return f.factory.NewTimeBuilder()
}

func (f *defaultFactory) NewDateTimeBuilder() types.DateTimeBuilder {
	return f.factory.NewDateTimeBuilder()
}

func (f *defaultFactory) NewEdgeBuilder() types.EdgeBuilder {
	return f.factory.NewEdgeBuilder()
}

func (f *defaultFactory) initDriver() error {
	factory, err := types.GetFactoryDriver(f.o.version)
	if err != nil {
		return err
	}
	f.factory = factory
	return nil
}

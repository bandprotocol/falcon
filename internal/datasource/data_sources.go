package datasource

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// Source defines an interface for retrieving data.
type Source interface {
	GetName() string
	GetData(ctx context.Context) (uint64, error)
}

// DataSources defines a list of source interface objects.
type DataSources struct {
	sources []Source
	log     *zap.Logger
}

// NewDataSources creates a new DataSources object.
func NewDataSources(sources []Source, log *zap.Logger) DataSources {
	return DataSources{
		sources: sources,
		log:     log,
	}
}

// GetDataResult defines the result of getting data from a specific source.
type GetDataResult struct {
	SourceName string
	Data       uint64
	Err        error
}

// NewGetDataResult creates a new GetDataResult object.
func NewGetDataResult(sourceName string, data uint64, err error) GetDataResult {
	return GetDataResult{
		SourceName: sourceName,
		Data:       data,
		Err:        err,
	}
}

// GetData gets data from multiple sources concurrently.
func (ds DataSources) GetData(ctx context.Context) (uint64, error) {
	if len(ds.sources) == 0 {
		ds.log.Error("no sources found")
		return 0, fmt.Errorf("no sources")
	}

	// get data from multiple sources concurrently
	resCh := make(chan GetDataResult, len(ds.sources))
	var wg sync.WaitGroup
	for _, s := range ds.sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			data, err := s.GetData(ctx)
			resCh <- NewGetDataResult(s.GetName(), data, err)
		}(s)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	results := make([]uint64, 0, len(ds.sources))
	for res := range resCh {
		if res.Err != nil {
			ds.log.Debug(
				fmt.Sprintf("Failed to get data: %s", res.Err),
				zap.String("data_source_name", res.SourceName),
			)
			continue
		}

		results = append(results, res.Data)
	}

	if len(results) == 0 {
		ds.log.Error("not retrieved any data from predefined sources")
		return 0, fmt.Errorf("no data")
	}

	// calculate the average of the data
	final, err := median(results)
	if err != nil {
		ds.log.Error("Failed to calculate the median", zap.Error(err))
		return 0, err
	}

	return final, nil
}

// Config defines interface for creating a new source object from the configuration.
type Config interface {
	Validate() error
	NewSource() (Source, error)
}

// DecodeDataSourceConfigHook decodes the data source configuration.
func DecodeDataSourceConfigHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if !to.Implements(reflect.TypeOf((*Config)(nil)).Elem()) {
		return data, nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot convert to map[string]interface{}")
	}

	sourceTypeStr, ok := dataMap["source_type"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot convert source_type to string")
	}

	switch ToSourceType(sourceTypeStr) {
	case SourceTypeFix:
		var conf FixSourceConfig
		if err := decodeConfig(data, &conf); err != nil {
			return nil, err
		}

		return conf, nil
	case SourceTypeWeb3Legacy:
		var conf Web3LegacySourceConfig
		if err := decodeConfig(data, &conf); err != nil {
			return nil, err
		}

		return conf, nil
	case SourceTypeWeb3EIP1559:
		var conf Web3EIP1559SourceConfig
		if err := decodeConfig(data, &conf); err != nil {
			return nil, err
		}

		return conf, nil
	default:
		return data, nil
	}
}

// decodeConfig decodes the input data into the output configuration.
func decodeConfig(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			DecodeSourceTypeHook,
		),
		Result: output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

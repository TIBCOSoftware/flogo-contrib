package tf

import (
	"fmt"

	tensorflow "github.com/TIBCOSoftware/flogo-contrib/activity/inference/tensorflow/tensorflow/core/example"
)

func Example(features map[string]interface{}) (*tensorflow.Example, error) {
	result := make(map[string]*tensorflow.Feature)
	for k, v := range features {
		switch t := v.(type) {
		case []byte:
			result[k] = toBytes(t)
		case [][]byte:
			result[k] = toBytesList(t)
		case string:
			result[k] = toBytes([]byte(t))
		case []string:
			b := make([][]byte, len(t))
			for i, s := range t {
				b[i] = []byte(s)
			}
			result[k] = toBytesList(b)
		case float32:
			result[k] = toFloat(t)
		case []float32:
			result[k] = toFloatList(t)
		case float64:
			// For now just convert float64 to float32
			result[k] = toFloat(float32(t))
		case []float64:
			// For now just convert float64 to float32
			f := make([]float32, len(t))
			for i, f64 := range t {
				f[i] = float32(f64)
			}
			result[k] = toFloatList(f)
		case int64:
			result[k] = toInt64(t)
		case []int64:
			result[k] = toInt64List(t)
		case int:
			result[k] = toInt64(int64(t))
		case []int:
			ints := make([]int64, len(t))
			for i, ii := range t {
				ints[i] = int64(ii)
			}
			result[k] = toInt64List(ints)
		case *tensorflow.Feature:
			result[k] = t
		default:
			return nil, fmt.Errorf("example: unsupported feature type %T: %q", t, t)
		}
	}
	return &tensorflow.Example{
		Features: &tensorflow.Features{result},
	}, nil
}

func toBytes(value []byte) *tensorflow.Feature {
	values := [][]byte{value}
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func toBytesList(values [][]byte) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_BytesList{&tensorflow.BytesList{values}}}
}

func toFloat(value float32) *tensorflow.Feature {
	values := []float32{value}
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func toFloatList(values []float32) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_FloatList{&tensorflow.FloatList{values}}}
}

func toInt64(value int64) *tensorflow.Feature {
	values := []int64{value}
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}

func toInt64List(values []int64) *tensorflow.Feature {
	return &tensorflow.Feature{&tensorflow.Feature_Int64List{&tensorflow.Int64List{values}}}
}

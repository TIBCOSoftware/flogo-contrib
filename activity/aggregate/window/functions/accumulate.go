package functions

func AggregateBlocksAccumulate(blocks []interface{}, start int, size int) interface{} {

	accum := make([]interface{}, 0, len(blocks))

	for i := 0; i < len(blocks); i++ {

		accum = append(accum, blocks[(start+i)%len(blocks)])
	}

	return accum
}

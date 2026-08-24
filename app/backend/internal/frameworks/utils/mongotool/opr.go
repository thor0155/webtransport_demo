package mongotool

type Operator = string

const (
	OperatorEq           Operator = "$eq"
	OperatorNe           Operator = "$ne"
	OperatorGt           Operator = "$gt"
	OperatorLt           Operator = "$lt"
	OperatorGte          Operator = "$gte"
	OperatorLte          Operator = "$lte"
	OperatorSet          Operator = "$set"
	OperatorAdd          Operator = "$add"
	OperatorIfNull       Operator = "$ifNull"
	OperatorConcatArrays Operator = "$concatArrays"
)

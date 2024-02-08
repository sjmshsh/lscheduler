package step

// SimpleStep .
type SimpleStep struct {
	name      string
	inputPtr  interface{}
	outputPtr interface{}
}

// NewSimpleStep .
func NewSimpleStep(name string, inputPtr, outputPtr interface{}) *SimpleStep {
	return &SimpleStep{
		name:      name,
		inputPtr:  inputPtr,
		outputPtr: outputPtr,
	}
}

// Name .
func (b SimpleStep) Name() string {
	return b.name
}

// InputPtr .
func (b SimpleStep) InputPtr() interface{} {
	return b.inputPtr
}

// OutputPtr .
func (b SimpleStep) OutputPtr() interface{} {
	return b.outputPtr
}

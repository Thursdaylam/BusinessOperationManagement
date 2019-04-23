package node

import (
	"testing"
)

func TestTraversalAccuracy(t *testing.T) {

	manual := Node{Value: 0.0, Left: nil, Right: nil, Prob: 0.712}

	manual.Left = &Node{Value: 1.0, Left: nil, Right: nil, Prob: 0}
	manual.Right = &Node{Value: 0.0, Left: nil, Right: nil, Prob: 0}

	value := EvaluateWorker(manual, 1000000)

	if value-0.712 > 0.712/100 {
		t.Errorf("Expected value to converge within 1 percent but got %f instead", value)
	}
}

package node

import (
	"math/rand"
)

// Node Has extra Prob float64 for decision tree
type Node struct {
	Value float64
	Left  *Node
	Right *Node
	Prob  float64
}

// EvaluateWorker takes in iterations to run multiple iterations
func EvaluateWorker(node Node, iterations int) float64 {
	c := make(chan float64, 8)
	token := make(chan bool, 8)
	sumChan := make(chan float64)

	go func(sumChan chan float64) {
		var sum float64
		for i := 0; i < iterations; i++ {
			pollValue := <-c
			sum = sum + pollValue
		}
		sumChan <- sum
	}(sumChan)

	for i := 0; i < iterations; i++ {
		token <- true
		go EvaluateWrapper(node, c, token)
	}

	//spawnWorkers(node, iterations, c)

	sum := <-sumChan
	average := sum / (float64(iterations))
	return average
}

// EvaluateWrapper wraps the Evaluate function to allow for channel passing
func EvaluateWrapper(node Node, c chan float64, token chan bool) {
	c <- Evaluate(node)
	<-token
}

// Evaluate simulates one walk of the decision tree and gives you a value
func Evaluate(node Node) float64 {
	if node.Left == nil || node.Right == nil {
		return node.Value
	}

	poll := rand.Float64()

	if poll < node.Prob {
		return Evaluate(*(node.Left)) + node.Value
	}

	return Evaluate(*(node.Right)) + node.Value

}

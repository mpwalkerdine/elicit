package calculator

// Calculator is a contrived example for specification purposes.
type Calculator struct{}

// Add is a + b.
func (c Calculator) Add(a, b int) int { return a + b }

// Sub is a - b.
func (c Calculator) Sub(a, b int) int { return a - b }

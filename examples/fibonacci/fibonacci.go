package fibonacci

// Sequence generates a Fibonacci slice
func Sequence(n, m int) []int {
	f := make([]int, m-n+1)

	a, b := 1, 1
	for i := 0; i <= m; i++ {
		if i > 2 {
			t := a + b
			a = b
			b = t
		}

		if n <= i && i <= m {
			f[i-n] = b
		}
	}

	return f
}

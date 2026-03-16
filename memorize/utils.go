package memorize

import (
	"math"

	"github.com/fatih/color"
)

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func EditDiff(a, b string) (string, int) {
	n, m := len(a), len(b)

	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j-1], dp[i-1][j], dp[i][j-1])
			}
		}
	}

	res := ""
	i, j := n, m
	remove := color.New(color.BgRed, color.FgWhite).SprintFunc()
	add := color.New(color.BgGreen, color.FgWhite).SprintFunc()

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			res = string(a[i-1]) + res
			i--
			j--
		} else if i > 0 && (j == 0 || dp[i][j] == dp[i-1][j]+1) {
			res = remove(string(a[i-1])) + res
			i--
		} else if j > 0 && (i == 0 || dp[i][j] == dp[i][j-1]+1) {
			res = add(string(b[j-1])) + res
			j--
		} else {
			res = remove(string(a[i-1])) + add(string(b[j-1])) + res
			i--
			j--
		}
	}

	return res, dp[n][m]
}

func QFromSimilarity(sim float64) int {

	q := int(math.Round(sim * 5))
	if q < 0 {
		q = 0
	} else if q > 5 {
		q = 5
	}
	return q
}

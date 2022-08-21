//给你一个整数 n ，统计并返回各位数字都不同的数字 x 的个数，其中 0 <= x < 10ⁿ 。
//
//
//
//
// 示例 1：
//
//
//输入：n = 2
//输出：91
//解释：答案应为除去 11、22、33、44、55、66、77、88、99 外，在 0 ≤ x < 100 范围内的所有数字。
//
//
// 示例 2：
//
//
//输入：n = 0
//输出：1
//
//
//
//
//
//
// 提示：
//
//
// 0 <= n <= 8
//
// Related Topics 数学 动态规划 回溯 👍 215 👎 0
package cn

//leetcode submit region begin(Prohibit modification and deletion)
func countNumbersWithUniqueDigits(n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return 10
	}
	ans, cur := 10, 9
	for i := 0; i < n-1; i++ {
		cur *= 9 - i
		ans += cur
	}
	return ans
}

//leetcode submit region end(Prohibit modification and deletion)

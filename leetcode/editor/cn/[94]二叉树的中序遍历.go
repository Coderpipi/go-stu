// 给定一个二叉树的根节点 root ，返回 它的 中序 遍历 。
//
// 示例 1：
//
// 输入：root = [1,null,2,3]
// 输出：[1,3,2]
//
// 示例 2：
//
// 输入：root = []
// 输出：[]
//
// 示例 3：
//
// 输入：root = [1]
// 输出：[1]
//
// 提示：
//
// 树中节点数目在范围 [0, 100] 内
// -100 <= Node.val <= 100
//
// 进阶: 递归算法很简单，你可以通过迭代算法完成吗？
//
// Related Topics 栈 树 深度优先搜索 二叉树 👍 2045 👎 0
package cn

// leetcode submit region begin(Prohibit modification and deletion)
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func inorderTraversal(root *TreeNode) []int {
	var (
		stk []*TreeNode
		res []int
	)
	for root != nil || len(stk) > 0 {
		for root != nil {
			stk = append(stk, root)
			root = root.Left
		}
		root = stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		res = append(res, root.Val)
		root = root.Right
	}

	return res
}

// leetcode submit region end(Prohibit modification and deletion)

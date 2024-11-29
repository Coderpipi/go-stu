// 给你一棵二叉树的根节点 root ，返回其节点值的 后序遍历 。
//
//
//
// 示例 1：
//
//
// 输入：root = [1,null,2,3]
// 输出：[3,2,1]
//
//
// 示例 2：
//
//
// 输入：root = []
// 输出：[]
//
//
// 示例 3：
//
//
// 输入：root = [1]
// 输出：[1]
//
//
//
//
// 提示：
//
//
// 树中节点的数目在范围 [0, 100] 内
// -100 <= Node.val <= 100
//
//
//
//
// 进阶：递归算法很简单，你可以通过迭代算法完成吗？
//
// Related Topics 栈 树 深度优先搜索 二叉树 👍 1163 👎 0

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
func postorderTraversal(root *TreeNode) []int {
	// 后序遍历, 左->右->根
	stack := []*TreeNode{}
	res := []int{}

	// 用于记录上一个访问的节点
	var prev *TreeNode
	for root != nil || len(stack) > 0 {

		// 遍历左子树, 并将左子树的节点入栈
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}

		// 当前root节点的左子树已经遍历完, 出栈当前节点
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 如果当前节点没有右子树, 则当前节点后序遍历完成, 直接将当前节点加入结果
		// 如果当前节点的右子节点等于上一次访问的节点, 则代表当前节点的右子树已经全部遍历完成, 也将当前节点加入结果, 完成后序遍历
		// 否则就将当前节点入栈, 继续遍历其右子树
		if root.Right == nil || root.Right == prev {
			res = append(res, root.Val)
			prev = root
			root = nil
		} else {
			stack = append(stack, root)
			root = root.Right
		}

	}
	return res
}

// leetcode submit region end(Prohibit modification and deletion)

package AddTwoNumbers

type ListNode struct {
	Val  int
	Next *ListNode
}

// 遍历相加
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	carryFlag := 0      // 进位标记
	this := &ListNode{} // 当前链表位置
	head := this        // 记录链表头位置

	for l1 != nil || l2 != nil || carryFlag == 1 {
		val := 0

		if l1 == nil && l2 == nil {
			// 都到头了，只处理进位
			val = carryFlag
		} else if l1 == nil {
			// l1到头了，只读l2的值
			val = l2.Val + carryFlag
			// 指针向后移动
			l2 = l2.Next
		} else if l2 == nil {
			// l2到头了，只读l1的值
			val = l1.Val + carryFlag
			// 指针向后移动
			l1 = l1.Next
		} else {
			// 相加
			val = l2.Val + l1.Val + carryFlag

			// 指针向后移动
			l1 = l1.Next
			l2 = l2.Next
		}

		if val >= 10 {
			// 进位
			carryFlag = 1
			val -= 10
		} else {
			// 进位重置
			carryFlag = 0
		}

		// 赋值结果
		this.Next = &ListNode{Val: val}
		// 指针向后移动
		this = this.Next
	}

	return head.Next
}

package ReverseLinkedList

// 遍历，主要是将当前指针指向前一个元素位置，所以要一个临时变量保存前一个元素的地址

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode // 当前链表位置的前一个元素地址，初始是nil

	for head != nil {
		next := head.Next // 存储下一个元素地址
		head.Next = prev  // 重置next指针到前一个元素
		prev = head       // 后移指针到当前元素位置
		head = next       // 后移指针到下一个元素位置
	}

	return prev
}

package AddTwoNumbers

import (
	"reflect"
	"testing"
)

func Test_addTwoNumbers(t *testing.T) {
	makeListNode := func(values []int) *ListNode {
		listNode := &ListNode{}
		head := listNode

		for _, v := range values {
			listNode.Next = &ListNode{Val: v}
			listNode = listNode.Next
		}

		return head.Next
	}

	type args struct {
		l1 *ListNode
		l2 *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "Example 1",
			args: args{
				l1: makeListNode([]int{2, 4, 3}),
				l2: makeListNode([]int{5, 6, 4}),
			},
			want: makeListNode([]int{7, 0, 8}),
		},
		{
			name: "Example 2",
			args: args{
				l1: makeListNode([]int{0}),
				l2: makeListNode([]int{0}),
			},
			want: makeListNode([]int{0}),
		},
		{
			name: "Example 3",
			args: args{
				l1: makeListNode([]int{9, 9, 9, 9, 9, 9, 9}),
				l2: makeListNode([]int{9, 9, 9, 9}),
			},
			want: makeListNode([]int{8, 9, 9, 9, 0, 0, 0, 1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addTwoNumbers(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addTwoNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

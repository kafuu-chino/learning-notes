package ReverseLinkedList

import (
	"reflect"
	"testing"
)

func Test_reverseList(t *testing.T) {
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
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "Example 1",
			args: args{
				head: makeListNode([]int{1, 2, 3, 4, 5}),
			},
			want: makeListNode([]int{5, 4, 3, 2, 1}),
		},
		{
			name: "Example 2",
			args: args{
				head: makeListNode([]int{1, 2}),
			},
			want: makeListNode([]int{2, 1}),
		},
		{
			name: "Example 3",
			args: args{
				head: makeListNode([]int{}),
			},
			want: makeListNode([]int{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverseList(tt.args.head); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reverseList() = %v, want %v", got, tt.want)
			}
		})
	}
}

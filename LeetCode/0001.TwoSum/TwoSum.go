package TwoSum

// 两数之和可以转化为，一个数字找另一个数字，另一个数字等于和减去当前数字。
// 即先计算出另一个数字去map中匹配，匹配不到则将要匹配的数字加入map，直到匹配成功。

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i, num := range nums {
		// 找到是否有匹配的数字
		if v, ok := m[num]; ok {
			return []int{v, i}
		}

		// 将要匹配的数字加入
		m[target-num] = i
	}

	return []int{}
}

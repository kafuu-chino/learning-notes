package LongestSubstringWithoutRepeatingCharacters

// 分析题目，首先想到暴力解法，遍历每个元素，分别计算长度。
// 思考下但是过程中是否会产生冗余计算，如"afbcbef"，得到"afbcb"时候，之后的子集"fbcb"等属于无效遍历，应该跳过到'c'继续遍历。
// 然后思考下，如"abcbef"从'c'开始，"cb"其实在之前的遍历中已经知晓，能否不要重复遍历，理论上一次遍历就知道所有元素的信息。
// 一次遍历，考虑使用map记录下元素信息（值和下标），匹配也很方便。

func lengthOfLongestSubstring(s string) int {
	// 下标的处理容易乱，首先需要理清楚，两个下标b-a，是(a,b]，所以计算的时候用原始下标即可。
	// 但是有个问题，从0开始遍历，那么下标为0的元素其实就不包括在内，所以采取从-1开始遍历的策略

	m := make(map[byte]int) // 记录元素信息和下标，防止无效遍历
	r := -1                 // 跳过无效遍历，记录从哪里开始遍历，需要包含第一个元素
	max := 0                // 最大长度

	for i := 0; i < len(s); i++ {
		if mv, ok := m[s[i]]; ok && r <= mv {
			// 找到是否有匹配重复元素，有的话跳到前一个匹配的元素，跳过无效遍历
			// 计算不包含当前元素，用原始数据即可
			r = mv
		}

		if l := i - r; l > max {
			// 判断长度是否超过最大值，超过的话替换
			max = l
		}

		// 记录值和下标
		m[s[i]] = i
	}

	return max
}
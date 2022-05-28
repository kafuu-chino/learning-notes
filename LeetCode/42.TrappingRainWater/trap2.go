package trap

// 双指针，左右各遍历到最高值
func trap2(height []int) int {
	leftHeight, leftIndex := 0, 0   // 左指针标记的到的高度和下标
	rightHeight, rightIndex := 0, 0 // 右指针标记的到的高度和下标

	holder := 0    // 总占位方块数
	totalRain := 0 // 总雨水数，占位方块最后返回一次性减掉，减少中间计算

	for i, j := 0, len(height)-1; i <= j; {
		if leftHeight == 0 && height[i] != 0 {
			// 标记当前，找到右侧第一个等于或者大于自身的方块
			leftHeight, leftIndex = height[i], i
		} else if i == leftIndex {
			// 左指针没动什么都不做
		} else {
			if height[i] >= leftHeight {
				// 结算一次
				totalRain += (i - leftIndex - 1) * leftHeight
				leftHeight, leftIndex = height[i], i
			} else {
				// 累加已经存在的方块，到时候减掉
				holder += height[i]
			}
		}

		if rightHeight == 0 && height[j] != 0 {
			// 标记当前，找到右侧第一个等于或者大于自身的方块
			rightHeight, rightIndex = height[j], j
		} else if j == rightIndex {
			// 右指针没动什么都不做
		} else {
			if height[j] >= rightHeight {
				// 结算一次
				totalRain += (rightIndex - j - 1) * rightHeight
				rightHeight, rightIndex = height[j], j
			} else {
				// 累加已经存在的方块，到时候减掉
				holder += height[j]
			}
		}

		if leftHeight <= rightHeight {
			i++
		} else {
			j--
		}
	}

	return totalRain - holder
}

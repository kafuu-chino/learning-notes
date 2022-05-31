package TrappingRainWater

// 初步思考很容易想到，找到最高点，往两边找到次高点为一个坑，矩形减去方块占位等于雨水数量，以次高点作为新的最高点继续找到次高点重复直到边界。
// 但是次高无法一次确认，需要遍历到头对比，所以转换一下思维。
// 从低到高遍历，寻找高于自身的方块就可以算作一个坑，遍历找到最高点。

// 先从左遍历找到最高值，再从右边遍历
func trap1(height []int) int {
	markHeight, markIndex := 0, 0 // 标记到的高度和下标
	holder := 0                   // 占位方块数
	totalRain := 0                // 总雨水数

	// 计算最高方块左边的雨水，顺便找出最高方块和下标
	for i, v := range height {
		if markHeight == 0 && v != 0 {
			// 标记当前，找到右侧第一个等于或者大于自身的方块
			markHeight, markIndex = v, i
		} else {
			if v >= markHeight {
				// 结算一次，雨水 = (跨度坐标*矮的高度) - 占位方块数
				totalRain += (i-markIndex-1)*markHeight - holder
				// 设置高度和下标到新的位置
				markHeight, markIndex = v, i
				holder = 0
			} else {
				// 累加已经存在的方块，结算时候减掉
				holder += v
			}
		}
	}

	maxIndex := markIndex        // 记录最大值
	markHeight, markIndex = 0, 0 // 清零
	holder = 0                   // 清零

	// 计算方块右边的雨水
	for i := len(height) - 1; i >= maxIndex; i-- {
		if markHeight == 0 && height[i] != 0 {
			// 标记当前，找到右侧第一个等于或者大于自身的方块
			markHeight, markIndex = height[i], i
		} else {
			if height[i] >= markHeight {
				// 结算一次，雨水 = (跨度坐标*矮的高度) - 占位方块数
				totalRain += (markIndex-i-1)*markHeight - holder
				// 设置高度和下标到新的位置
				markHeight, markIndex = height[i], i
				holder = 0
			} else {
				// 累加已经存在的方块，结算时候减掉
				holder += height[i]
			}
		}
	}

	return totalRain
}

// 双指针，左右同时遍历
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
				// 结算一次，雨水 = (跨度坐标*矮的高度) - 占位方块数
				totalRain += (i - leftIndex - 1) * leftHeight
				// 设置高度和下标到新的位置
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
				// 结算一次，雨水 = (跨度坐标*矮的高度) - 占位方块数
				totalRain += (rightIndex - j - 1) * rightHeight
				// 设置高度和下标到新的位置
				rightHeight, rightIndex = height[j], j
			} else {
				// 累加已经存在的方块，到时候减掉
				holder += height[j]
			}
		}

		// 根据左右高度选择下一步移动的指针
		if leftHeight <= rightHeight {
			i++
		} else {
			j--
		}
	}

	return totalRain - holder
}

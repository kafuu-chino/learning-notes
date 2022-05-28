package trap

// 向右寻找高于自己的格子视为一个坑，进行结算，直到最高的格子，再从右按一样的逻辑遍历到最高的格子

// 先从左遍历找到最高值，再从右边遍历
func trap1(height []int) int {
	markHeight, markIndex := 0, 0 // 标记到的高度和下标
	holder := 0                   // 占位方块数
	totalRain := 0                // 总雨水数

	// 计算最高方块左边的雨水，顺便找出最高方块和坐标
	for i, v := range height {
		if markHeight == 0 && v != 0 {
			// 标记当前，找到右侧第一个等于或者大于自身的方块
			markHeight, markIndex = v, i
		} else {
			if v >= markHeight {
				// 结算一次
				totalRain += (i-markIndex-1)*markHeight - holder
				markHeight, markIndex = v, i
				holder = 0
			} else {
				// 累加已经存在的方块，到时候减掉
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
				// 结算一次
				totalRain += (markIndex-i-1)*markHeight - holder
				markHeight, markIndex = height[i], i
				holder = 0
			} else {
				// 累加已经存在的方块，到时候减掉
				holder += height[i]
			}
		}
	}

	return totalRain
}

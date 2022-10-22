package main

import (
	"fmt"
	"strings"
)

// 两数之和

func twoSum(nums []int, target int) []int {
	var res = map[int][]int{}
	for k, v := range nums {
		res[v] = append(res[v], k)
	}
	rr := [][]int{}
	for k, v := range nums {
		value := target - v
		if vals, ok := res[value]; ok {
			// 重复值切片
			for _, val := range vals {
				if val != k {
					rr = append(rr, []int{k, val})
				}

			}
		}
	}
	fmt.Println("==", rr)
	return rr[0]
}
func twoSum2(nums []int, target int) []int {
	var m = make(map[int]int, len(nums))
	for index, val := range nums {
		if preIndex, ok := m[target-val]; ok {
			fmt.Println([]int{preIndex, index})
			return []int{preIndex, index}
		} else {
			m[val] = index
		}
	}
	return []int{}
}

// 罗马数字转整数
func romanToInt(s string) int {
	// data = map[string]int{
	// 	"I": 1,
	// 	"V": 5,
	// 	"X": 10,
	// 	"L": 50,
	// 	"C": 100,
	// 	"D": 500,
	// 	"M": 1000,
	// }
	res := strings.Split(s, "")
	for k := len(res) - 1; k >= 0; k-- {
		fmt.Println(k, res[k])
	}

	return 0
}

// 整数转二进制，并判断交叉数

// 最长公共前缀

// 合并两个有序链表

// 有效的括号 {[()]}[{}]{[}]

// 升序排列的数组原地去重（不申请新空间）

func main() {
	// twoSum2([]int{2, 5, 5, 11}, 10)
	romanToInt("LVIII")
}

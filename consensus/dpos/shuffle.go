package dpos

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"
)

const delegateInterval = 10 // 出块间隔10秒

// dpos乱序方法,同一个轮回里面传入不同的高度，排序列表是一样的
// delegateNumber:代理节点数量
// height:本次产块的块高
func Shuffle(height int64, delegateNumber int) []int {
	var truncDelegateList []int

	for i := 0; i < delegateNumber; i++ {
		truncDelegateList = append(truncDelegateList, i)
	}

	seed := math.Floor(float64(height / int64(delegateNumber)))
	//seed := strconv.FormatFloat(math.Floor(float64(height/101)), 'E', -1, 64)

	if height%int64(delegateNumber) > 0 {
		seed += 1
	}
	seedSource := strconv.FormatFloat(seed, 'E', -1, 64)
	var buf bytes.Buffer
	buf.WriteString(seedSource)
	hash := sha256.New()
	hash.Write(buf.Bytes())
	md := hash.Sum(nil)
	currentSend := hex.EncodeToString(md)

	delCount := len(truncDelegateList)
	for i := 0; i < delCount; i++ {
		for x := 0; x < 4 && i < delCount; i++ {
			newIndex := int(currentSend[x]) % delCount
			// 元素互换
			truncDelegateList[newIndex], truncDelegateList[i] = truncDelegateList[i], truncDelegateList[newIndex]
			x++
		}
	}
	return truncDelegateList

}

type ShuffleDel struct {
	//delegateIndex int
	workTime      int64
	address string
}

// lastRoundBlockHeight 上轮周期最后一个块高
// lastRoundBlockHeightTime 上轮周期最后一个块的ntp时间戳
// 本轮出块代理总数
func ShuffleNewRound(lastRoundBlockHeight int64, lastRoundBlockHeightTime int64, delegateNumber int,CurrentDposList []Delegate) []ShuffleDel {
	var newRoundList []ShuffleDel
	truncDelegateList := Shuffle(lastRoundBlockHeight+1, delegateNumber)
	for height := lastRoundBlockHeight + 1; height <= lastRoundBlockHeight+int64(delegateNumber); height++ {
		i := height % int64(delegateNumber)
		delegateIndex := truncDelegateList[i]
		workTime := lastRoundBlockHeightTime + (height-lastRoundBlockHeight)*delegateInterval*1000
		newRoundList = append(newRoundList, ShuffleDel{ workTime: workTime,address:CurrentDposList[delegateIndex].Address})
	}
	return newRoundList
}

// 获取当前产块节点在列表中的索引
func GetCurrentProduceNode(height int64, delegateNumber int) int {
	i := height % int64(delegateNumber)
	truncDelegateList := Shuffle(height, delegateNumber)
	return truncDelegateList[i]
}

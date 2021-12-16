package utils

import (
	"math/rand"
	"time"
)

func SimulateTransmission(tags, frameSize int) (colisions, emptySlots, tagsReaded int) {
	rand.Seed(time.Now().UnixNano())
	//fmt.Println(tags, frameSize)
	tagsReaded = 0
	colisions = 0
	emptySlots = 0
	slots := make([]int, frameSize)

	for i := 0; i < tags; i++ { // simulate each tag transmiting at random slot
		slotUsed := rand.Intn(frameSize)
		slots[slotUsed]++
		//fmt.Println(slotUsed)
	}
	//fmt.Println(slots)
	for i := 0; i < frameSize; i++ { // reading each slot
		if slots[i] == 0 {
			emptySlots++
		} else if slots[i] == 1 {
			tagsReaded++
		} else if slots[i] > 1 {
			colisions++
		}
	}
	return colisions, emptySlots, tagsReaded
}

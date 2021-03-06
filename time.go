// Copyright © 2018 J. Strobus White.
// This file is part of the blocktop blockchain development kit.
//
// Blocktop is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Blocktop is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with blocktop. If not, see <http://www.gnu.org/licenses/>.

package kernel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type KernelTime struct {
	blockFrequency float64
	blockInterval  time.Duration
	intervalLen    string
	cycleNumber    uint64
	startTime      int64
	cycleStartTime int64
}

var ktime *KernelTime

func initTime(blockFrequency float64) {
	ktime = &KernelTime{}
	ktime.setBlockFrequency(blockFrequency)
}

func (t *KernelTime) setBlockFrequency(rate float64) {
	interval := time.Duration(float64(time.Second) / rate)
	t.blockFrequency = rate
	t.blockInterval = interval
	t.intervalLen = strconv.FormatInt(int64(len(strconv.FormatInt(int64(interval), 10))), 10)
}

func (t *KernelTime) BlockFrequency() float64 {
	return t.blockFrequency
}

func (t *KernelTime) BlockInterval() time.Duration {
	return t.blockInterval
}

func (t *KernelTime) UpTime() time.Duration {
	return time.Duration(t.Nanos())
}

func (t *KernelTime) CycleNumber() uint64 {
	return t.cycleNumber
}

func (t *KernelTime) Nanos() int64 {
	return time.Now().UnixNano() - t.startTime
}

func (t *KernelTime) String() string {
	scycle := humanize.Comma(int64(t.CycleNumber()))
	cycleTime := (time.Now().UnixNano() - t.cycleStartTime) / 1000
	return fmt.Sprintf("%s.%s", scycle, leftPadZeroes(cycleTime, 6))
}

func (t *KernelTime) up() {
	t.startTime = time.Now().UnixNano()
}

func (t *KernelTime) startCycle() {
	now := time.Now().UnixNano()
	cycleTime := now - t.cycleStartTime
	metrics.setCycleTime(cycleTime)
	t.cycleNumber++
	t.cycleStartTime = now
}

func leftPadZeroes(val int64, overallLen int) string {
	v := fmt.Sprintf("%d", val)

	if len(v) > overallLen {
		v = v[0:overallLen]
	}

	s := fmt.Sprintf("%s%s", strings.Repeat("0", overallLen-len(v)), v)
	return s
}

// Twitter 的 Snowflake 算法 的实现
/*
SnowFlake 的结构如下（每部分用-分开）:
0 - 0000000000 0000000000 0000000000 0000000000 0 - 00000 - 00000 - 000000000000

- 最高位是符号位，正数是 0，负数是 1， id 一般是正数，因此最高位固定是 0；
- 41 位时间戳（毫秒级），注意，41 位时间戳不是存储当前时间的时间戳，而是存储时间戳的差值（当前时间戳 - 开始时间戳），这样能存的时间更长，开始时间一般指定为项目启动时间，由程序指定。可以使用 69 年：`(1<<41)/(1000*60*60*24*365)`；
- 10 位的机器相关位，可以部署在 1024 个节点，包括 5 位的 datacenterId（数据中心 ID） 和 5 位 workerId（工作机器 ID）；
- 12 位系列号，毫秒内的计数。12 位的计数顺序号支持每个节点每毫秒（同一机器，同一时间戳）产生 4096 个 ID 序号；

sonyflake 是 Sony 公司的一个开源项目，基本思路和 snowflake 差不多，不过位分配上稍有不同：
这里的时间只用了 39 个比特，但时间的位数变成了 10 ms, 所以理论上比 41 位表示的时间还要久(174年)
*/
package snowflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

// Init 需传入当前的机器 ID
func Init(machineID uint16) (err error) {
	sonyMachineID = machineID
	t, _ := time.Parse("2006-01-02", "2020-01-01")
	settings := sonyflake.Settings{
		StartTime: t,            // 如果不设置的话，默认是从 2014-09-01 00:00:00 +0000 UTC 开始
		MachineID: getMachineID, // 节点 id,可以由用户自定义函数，不定义的话，会默认将本机ip的低16位作为machineID
	}
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
}

// GenID 返回生成的 id 值
func GenID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sony flake not init")
		return
	}

	id, err = sonyFlake.NextID()
	return
}

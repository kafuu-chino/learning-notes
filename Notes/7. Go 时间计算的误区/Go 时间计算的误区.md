时间的计算有很多的坑，首先要打破常规思维。以下都是错误的

1.  1天有24小时
2.  1小时有60分钟
3.  1周有7天
4.  1年有365天
5.  [更多](https://infiniteundo.com/post/25326999628/falsehoods-programmers-believe-about-time)

如[uber-style](https://github.com/uber-go/guide/blob/master/style.md#use-time-to-handle-time)所述，避免用time.Duration，而使用time.Time去处理时间。

这里解释一下一天24小时的问题，根据不同的时区，有些时区于夏令时当日，一天只有23小时，而处于冬令时当日，一天却有25小时。请参考[夏令时wiki](https://zh.m.wikipedia.org/zh-hans/%E5%A4%8F%E6%97%B6%E5%88%B6)。

计算机和程序对于时间的计算，全部基于一个时间数据库，go的源码内部，目录go/lib/time下也存有一份用于时间计算，计算规则是准对于时区，时区大部分是基于城市，少部分使用公共时区。请参考[time-zones](https://www.iana.org/time-zones)。

结合起来，go对于夏令时的处理也是基于时区，当前城市或者地区时区是需要考虑夏令时的情况下，go也会自动处理夏令时，如当日是夏令时，凌晨2:00会跳到3:00，进行时间计算比如5:00-00:00也只有4小时。本质上不同时区或者夏令时的计算本质上是基于一个偏移量，当夏令时触发，本质上也是偏移量+1，注意计算规则改变但是时区是不变的。

这里介绍一下格林尼治标准时间GMT和协调世界时UTC，注意GMT是一个时区，而UTC是一个时间标准，两者都是不考虑夏令时的，所以如果需要固定偏移量，可以使用GMT+8等作为时区进行时间计算。更多参考[The Difference Between GMT and UTC](https://www.timeanddate.com/time/gmt-utc-time.html)。
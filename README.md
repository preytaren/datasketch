# DataSketch
dataSketch GO 实现，datasketch是一系列基数计数算法，基数计算（cardinality counting）指的是统计一批数据中的不重复元素的个数，常见于计算独立用户数（UV）、维度的独立取值数等等。
通常的基数计数使用集合，bitmap等数据结构，能够精确的计算出结果，但是需要占用较大的存储空间。而sketch系列算法基于概率与统计，内存占用友好，能够估算出基数计数的结果，误差小于1%。
## ThetaSketch
- 支持流式处理
- 内存友好，内存占用固定
- 支持交，并，差集运算
- [ThetaSketch](https://datasketches.github.io/docs/Theta/InverseEstimate.html)

## HyperLogLog Sketch
- 支持流式处理
- 支持交，并集运算
- [HyperLogLog](https://www.jianshu.com/p/41256ac5b03f)
---
概述: QOS联盟链合约示例
---

# 定义mapper名称
定义mapper名称 同一业务链中mapper不要重名
```go
const (
	OrderMapperName = "order"
)
```


# 定义要存储的数据结构体

根据业务逻辑定义要存储数据结构体,本示例中存储了订单的一些基础信息

```go
type Order struct {
	Id          string       //订单编号
	OrderName   string       //订单名称
	OrderFrom   string       //订单发起者
	OrderTo     string       //订单接收者
	OrderAmount types.BigInt //订单金额
	Status      int          //转账状态 0成功  1未成功
}
```

# 定义mapper对象

根据golang的语法定义  如果一个struct嵌套了另一个匿名结构体，那么这个结构可以直接访问匿名结构体的方法，从而实现继承
OrderMapper继承了BaseMapper

```go
type OrderMapper struct {
	*mapper.BaseMapper
}
```

同时还要实现IMapper中的Copy方法

```go
func (mapper *OrderMapper) Copy() mapper.IMapper {
	cpyMapper := &OrderMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}
```


# 定义存储逻辑
根据业务逻辑自定义存储规则 在本例中我们把订单交易信息存储起来 提供存储和获取以及删除的方法 

```go
func (mapper *OrderMapper) SaveOrder(key string,order *Order) {
	mapper.Set([]byte(key), order)
}

func (mapper *OrderMapper) GetOrder(key string) *Order {
	var order Order
	exist := mapper.Get([]byte(key), &order)
	if !exist {
		return nil
	}
	return &order
}


func (mapper *OrderMapper) DeleteOrder(key string) {
	mapper.Del([]byte(key))
}


```




[process.go 调用合约SDK 负责签名和广播交易](process.md) 
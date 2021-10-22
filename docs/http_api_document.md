## 产品管理

### 获取产品列表
- **说明:** 获取产品列表
- **URL:** `http://HOST:PORT/api/v1/product`
- **方法:** `GET`
- **请求报文:**
```json
{
  "page": {
    "current": 1,
    "size": 10
  }
}
```
- **请求参数说明：**

| 参数  | 必填  | 说明  |
| ------------ | ------------ | ------------ |
| page.current  | 选填  | 分页页码，从 1 开始  |
| page.size  | 选填  | 分页数量，默认为20 |

- **应答报文:**

```json
{
  "code": "code",
  "message": "OK",
  "data": {
    "list":[{
      "productId": "XXXX",
      "productType": 1,
      "productName": "产品"
    }],
    "page": {
      "total": 30,
      "size": 20,
      "current": 1,
      "pages": 2
    }
  }
}
```
- **应答参数说明：**

| 参数  | 必填  | 说明  |
| ------------ | ------------ | ------------ |
| list.productId  | 必填  | 产品ID  |
| list.productType  | 必填  | 产品类型 |
| list.productName  | 选填  | 产品名称  |
| page.total  | 必填  | 总条数  |
| page.size  | 必填  | 每页条数  |
| page.current  | 必填  | 当前页码  |
| page.pages  | 必填  | 总页数  |
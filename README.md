# tmq

## 2.27 

无符号加法溢出

**题目**

写出一个具有如下原型的函数

`int uadd_ok(unsigned x, unsigned y);`

当x + y没有发生溢出时返回1

**思路**

根据公式:

x+_w^uy =
\begin{cases}
x + y < 2^w  & \text{正常} \\
2^w < x + y < 2^{w+1} & \text{溢出}
\end{cases}
  
如果当发生溢出时, x + y > 2^w, 此时结果s = x + y - 2^w

假设 y < 2^w, 则 s < x + y - 2^w

于是函数原型如下

**答案**

```c
int uadd_ok(unsigned char x, unsigned char y) {
    return (unsigned char) x + y > x? 1 : 0;
}
```

## 2.28

无符号求反

**问题**

我们能用一个十六进制数字来表示长度 w = 4的位模式。对于这些数字的无符号解释，使用等式(2.12)填写下表，给出所示数字的无符号加法逆元的位表示（用16进制）

2.12 公式

```math
-_w^ux = 
\begin{cases}
x &\text{x=0}\\
2^w - x & \text{x>0}
\end{cases}
```


x |  |-_u^4x | |
---|---|---|---
十六进制|十进制|十进制|十六进制
0       |0      |0  | 0x00
5       |5      |11 | 0xb
8       |8      |8  | 0x8
D       |13     |3  | 0x3
F       |15     |1  | 0x1






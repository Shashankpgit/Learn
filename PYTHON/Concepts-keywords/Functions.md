Below is a **one-page Python Functions Revision Guide** designed for rapid recall. You should be able to glance through this and refresh **all core concepts in one pass**.

---

## 1. What Is a Function

A function is a reusable block of code that performs a specific task.

```python
def function_name(parameters):
    # function body
    return value
```

---

## 2. Defining & Calling Functions

```python
def greet():
    print("Hello")

greet()  # function call
```

---

## 3. Function Parameters & Arguments

### a) Positional Arguments

```python
def add(a, b):
    return a + b

add(2, 3)
```

### b) Keyword Arguments

```python
add(b=3, a=2)
```

### c) Default Arguments

```python
def power(base, exp=2):
    return base ** exp

power(3)
power(3, 3)
```

### d) Variable-Length Arguments

**`*args` (tuple)**

```python
def total(*args):
    return sum(args)
```

**`**kwargs` (dictionary)**

```python
def user_info(**kwargs):
    return kwargs
```

---

## 4. Return Statement

```python
def square(x):
    return x * x
```

Multiple returns:

```python
def calc(a, b):
    return a+b, a-b
```

---

## 5. Function Without Return (Returns `None`)

```python
def log(msg):
    print(msg)
```

---

## 6. Scope of Variables

### Local Scope

```python
def func():
    x = 10
```

### Global Scope

```python
x = 5
def func():
    global x
    x = 10
```

---

## 7. Lambda (Anonymous Functions)

Used for short, one-line logic.

```python
square = lambda x: x * x
```

With `map`, `filter`:

```python
map(lambda x: x*x, nums)
filter(lambda x: x > 0, nums)
```

---

## 8. Docstrings (Documentation)

```python
def add(a, b):
    """Returns sum of two numbers"""
    return a + b
```

Access:

```python
help(add)
```

---

## 9. Type Hints (Optional but Professional)

```python
def add(a: int, b: int) -> int:
    return a + b
```

---

## 10. Recursive Functions

Function calling itself.

```python
def factorial(n):
    if n == 1:
        return 1
    return n * factorial(n-1)
```

---

## 11. Higher-Order Functions

Functions that accept or return other functions.

```python
def apply(func, x):
    return func(x)
```

---

## 12. Built-in Function Patterns (Very Important)

```python
map(func, iterable)
filter(func, iterable)
reduce(func, iterable)  # from functools
```

---

## 13. Common Mistakes (Quick Check)

❌ Mutable default arguments

```python
def func(x=[]):  # WRONG
```

✅ Correct

```python
def func(x=None):
    if x is None:
        x = []
```

---

## 14. Practical Patterns (Interview / Real Code)

### Input Validation

```python
def divide(a, b):
    if b == 0:
        return None
    return a / b
```

### Function as Config

```python
def prod(x, y): return x*y
operation = prod
```

---

## 15. Quick Mental Checklist Before Writing a Function

* What are inputs?
* Should it return or just perform action?
* Default values?
* Edge cases?
* Reusability?

---

## 16. Ultra-Short Summary (Memory Hook)

* `def` → define
* `()` → call
* `return` → output
* `*args` → many values
* `**kwargs` → key-value
* `lambda` → short logic
* recursion → function calls itself
* scope → local vs global

---

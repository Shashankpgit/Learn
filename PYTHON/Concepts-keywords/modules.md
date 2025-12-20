Below is a **Python Modules – Quick Reference Sheet** designed for **one-pass revision**, with **clear concepts and minimal, correct examples**. This covers **what modules are**, **how to use them**, and **all common patterns you are expected to know**.

---

# Python Modules – Quick Sheet

---

## 1. What Is a Module

A **module** is a Python file (`.py`) containing variables, functions, or classes that can be reused.

Example:

```python
# math_utils.py
def add(a, b):
    return a + b
```

---

## 2. Importing Modules (MOST IMPORTANT)

### a) Import Entire Module

```python
import math
print(math.sqrt(16))
```

### b) Import Specific Items

```python
from math import sqrt, pi
print(sqrt(16))
```

### c) Import with Alias

```python
import numpy as np
```

### d) Import Everything (NOT recommended)

```python
from math import *
```

---

## 3. Built-in Modules (Core Python)

### math – Mathematical Operations

```python
import math
math.sqrt(25)
math.factorial(5)
math.pi
```

---

### random – Random Values

```python
import random
random.randint(1, 10)
random.choice([1, 2, 3])
random.shuffle(my_list)
```

---

### datetime – Date & Time

```python
from datetime import datetime
now = datetime.now()
now.strftime("%Y-%m-%d")
```

---

### os – Operating System

```python
import os
os.getcwd()
os.listdir(".")
os.mkdir("test")
```

---

### sys – Python Runtime

```python
import sys
sys.argv
sys.exit()
sys.path
```

---

### time – Time Operations

```python
import time
time.sleep(2)
time.time()
```

---

## 4. Utility & Productivity Modules

### itertools – Advanced Iteration

```python
import itertools
list(itertools.permutations([1,2,3], 2))
```

---

### functools – Functional Programming

```python
from functools import reduce
reduce(lambda a,b: a+b, [1,2,3,4])
```

---

### collections – Specialized Containers

```python
from collections import Counter
Counter("banana")
```

---

## 5. File & Data Handling Modules

### json – JSON Data

```python
import json
json.dumps({"a": 1})
json.loads('{"a": 1}')
```

---

### csv – CSV Files

```python
import csv
with open("data.csv") as f:
    reader = csv.reader(f)
```

---

### pickle – Object Serialization

```python
import pickle
pickle.dump(obj, file)
pickle.load(file)
```

---

## 6. Regular Expressions

### re – Pattern Matching

```python
import re
re.findall(r"\d+", "abc123xyz")
```

---

## 7. Networking & Web

### requests – HTTP Requests (External)

```python
import requests
r = requests.get("https://api.example.com")
r.json()
```

---

### urllib – Built-in HTTP

```python
from urllib.request import urlopen
urlopen("https://example.com")
```

---

## 8. Concurrency Modules

### threading – Threads

```python
import threading
t = threading.Thread(target=func)
t.start()
```

---

### multiprocessing – CPU-bound Tasks

```python
from multiprocessing import Process
p = Process(target=func)
p.start()
```

---

## 9. Custom (User-Defined) Modules

Directory:

```
project/
 ├── main.py
 └── utils.py
```

```python
# utils.py
def greet():
    return "Hi"
```

```python
# main.py
import utils
utils.greet()
```

---

## 10. Packages (Modules Folder)

```
mypkg/
 ├── __init__.py
 ├── file1.py
 └── file2.py
```

Import:

```python
from mypkg import file1
```

---

## 11. `__name__ == "__main__"` (VERY IMPORTANT)

```python
def main():
    print("Run directly")

if __name__ == "__main__":
    main()
```

Used to prevent auto-execution on import.

---

## 12. Virtual Environment & External Modules

Install:

```bash
pip install requests
```

Import:

```python
import requests
```

---

## 13. Standard Library vs Third-Party

| Type        | Example                      |
| ----------- | ---------------------------- |
| Standard    | `math`, `os`, `sys`          |
| Third-party | `numpy`, `pandas`, `fastapi` |
| Custom      | Your own `.py` files         |

---

## 14. Common Mistakes (Exam + Real Life)

❌ File name same as module

```python
math.py  # shadows math module
```

❌ Circular imports
Module A imports B, B imports A

❌ Using `import *`

---

## 15. One-Glance Memory Map

* `import module`
* `from module import item`
* built-in ≠ third-party
* module → file
* package → folder
* `__main__` → execution control
* `pip` → external modules

---

## 16. When You Write Production Code

* Prefer **explicit imports**
* Use **aliases for long names**
* Group imports:

  1. Standard
  2. Third-partys
  3. Local

---

Below is a **Python Built-in Functions – Quick Reference Sheet**, focused on **introspection, utility, and day-to-day usage**. This is structured so you can **scan once and recall everything**, with **precise examples** for each function.

---

# Python Built-in Functions – Quick Sheet

Python provides **70+ built-in functions**. You do **not** import them — they are always available.

---

## 1. Introspection & Debugging Functions

*(Used to explore objects, modules, and runtime behavior)*

---

### `dir()` – Discover Attributes

Lists all attributes and methods of an object.

```python
import math
dir(math)
```

Typical use:

* Exploring a module
* Inspecting objects interactively

---

### `help()` – Documentation

Displays official documentation.

```python
help(str)
help(math.sqrt)
```

---

### `type()` – Object Type

```python
type(10)
type([])
```

---

### `id()` – Memory Identity

```python
id(10)
```

Used for:

* Understanding object references
* Debugging mutability

---

### `callable()` – Is It Callable?

```python
callable(print)
callable(10)
```

---

## 2. Attribute & Reflection Functions

---

### `getattr()` – Get Attribute Safely

```python
getattr(obj, "name", None)
```

---

### `setattr()` – Set Attribute Dynamically

```python
setattr(obj, "age", 25)
```

---

### `hasattr()` – Check Attribute

```python
hasattr(obj, "salary")
```

---

### `delattr()` – Delete Attribute

```python
delattr(obj, "temp")
```

---

## 3. Data Type Conversion Functions

---

### Numeric Conversions

```python
int("10")
float("3.14")
complex(2, 3)
```

---

### Sequence Conversions

```python
list("abc")
tuple([1,2])
set([1,1,2])
```

---

### String Conversion

```python
str(100)
```

---

## 4. Math & Aggregation Functions

---

### `abs()`

```python
abs(-10)
```

---

### `round()`

```python
round(3.1415, 2)
```

---

### `sum()`

```python
sum([1,2,3])
```

---

### `min()` / `max()`

```python
min([1,2,3])
max([1,2,3])
```

---

### `pow()`

```python
pow(2, 3)
```

---

## 5. Iterables & Loop Helpers

---

### `len()`

```python
len("hello")
```

---

### `enumerate()` – Index + Value

```python
for i, v in enumerate(["a","b"]):
    print(i, v)
```

---

### `range()`

```python
range(1, 5)
```

---

### `zip()` – Combine Iterables

```python
zip([1,2], ["a","b"])
```

---

### `reversed()`

```python
list(reversed([1,2,3]))
```

---

### `sorted()`

```python
sorted([3,1,2])
```

---

## 6. Logical & Validation Functions

---

### `all()` – All True?

```python
all([True, True])
```

---

### `any()` – Any True?

```python
any([False, True])
```

---

### `isinstance()` – Type Check

```python
isinstance(10, int)
```

---

### `issubclass()`

```python
issubclass(bool, int)
```

---

## 7. Input / Output Functions

---

### `print()`

```python
print("Hello", end=" ")
```

---

### `input()`

```python
name = input("Enter name: ")
```

---

## 8. Functional Programming Helpers

---

### `map()`

```python
map(lambda x: x*x, [1,2,3])
```

---

### `filter()`

```python
filter(lambda x: x > 0, [-1, 2, 3])
```

---

### `reduce()` *(functools)*

```python
from functools import reduce
reduce(lambda a,b: a+b, [1,2,3])
```

---

## 9. Memory & Execution Control

---

### `globals()` / `locals()`

```python
globals()
locals()
```

---

### `eval()` – Evaluate Expression *(Use carefully)*

```python
eval("2 + 3")
```

---

### `exec()` – Execute Code

```python
exec("a = 10")
```

---

## 10. Object Creation & Low-Level Functions

---

### `object()`

```python
obj = object()
```

---

### `hash()`

```python
hash("abc")
```

---

### `repr()` – Developer Representation

```python
repr("abc")
```

---

## 11. File & Encoding Helpers

---

### `open()`

```python
open("file.txt", "r")
```

---

### `bytes()` / `bytearray()`

```python
bytes("hi", "utf-8")
```

---

### `ord()` / `chr()`

```python
ord("A")
chr(65)
```

---

## 12. Error Handling Functions

---

### `Exception`

```python
raise Exception("Error")
```

---

### `BaseException`

Root of all exceptions.

---

## 13. One-Glance Cheat Memory (VERY IMPORTANT)

| Purpose  | Function                   |
| -------- | -------------------------- |
| Explore  | `dir()`, `help()`          |
| Check    | `type()`, `isinstance()`   |
| Iterate  | `range()`, `enumerate()`   |
| Combine  | `zip()`                    |
| Validate | `all()`, `any()`           |
| Convert  | `int()`, `str()`, `list()` |
| Execute  | `eval()`, `exec()`         |

---

## 14. When to Use `dir()` Specifically

Use `dir()` when:

* You see a **new module**
* You don’t know **available methods**
* You are debugging interactively

Example:

```python
import datetime
dir(datetime)
```

---

## 15. Interview Tip

If asked:

> “How do you explore an unfamiliar Python object?”

Answer:

> “Using `dir()` to list attributes and `help()` to read documentation.”

---

### Recommended Next Topics

* Built-in functions **commonly asked in interviews**
* Advanced introspection (`inspect` module)
* Writing your own `__repr__`, `__str__`
* Python data model (`__dunder__` methods)

If you want, I can convert this into a **single printable PDF cheat sheet** or provide **hands-on micro-exercises** for mastery.

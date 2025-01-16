# NebulaGraph Python Client

The Python client for NebulaGraph 5.x.

## Installation

```bash
pip install nebulagraph-python # not yet published
```

from source

```bash
cd python
pip install -e .
```

## Get Started

1. We could easily connect and get a query result.

```python
from nebulagraph_python.client import NebulaClient

# Create client
client = NebulaClient(
    hosts=["127.0.0.1:9669"],
    username="root",
    password="NebulaGraph01"
)

query = """
    RETURN 1 AS a, 2 AS b
"""
# Execute query
result = client.execute(query)
```

2. Then we could inspect the result ourselves.

```python
# Print the result in table style
result.print()
```
```
╔═══╤═══╗
║   │   ║
║ a │ b ║
║   │   ║
╟───┼───╢
║   │   ║
║ 1 │ 2 ║
║   │   ║
╚═══╧═══╝

Summary
├── Rows: 1
└── Latency: 1450μs
```

```python
# Get one row
row = result.one()

# Get one value
cell = row["a"].cast_primitive()

# Print its value
print(cell, type(cell))
```
```
1 <class 'int'>
```


3. We could actually get primitive values from the result set.

```python
print(result.as_primitive_by_column())
print(list(result.as_primitive_by_row()))
```
```
{'a': [1], 'b': [2]}
[{'a': 1, 'b': 2}]
```

4. If needed we could also get a pandas dataframe from the result.

We need to install pandas first.

```bash
pip install pandas
```

Then we could get a pandas dataframe like this:

```python
result = client.execute(query)
df = result.as_pandas_df()
```
```
   a  b
0  1  2
```

## Console Tools

Run `ngcli --help` to get the help message. An example to connect to NebulaGraph is as follows:

```bash
ngcli -h 127.0.0.1:9669 -u root -p NebulaGraph01
```

## Dev and Debug

### Environment Variables

1. Set `NG_PYTHON_DEBUG` to `true` to enable debug mode.

```bash
export NG_PYTHON_DEBUG=true
```

2. Set `NG_PYTHON_LOG_LEVEL` to `DEBUG` to enable debug logging.

```bash
export NG_PYTHON_LOG_LEVEL=DEBUG
```

3. Set `NG_PYTHON_LOG_SINK` to `stdout` or `file` to specify the logging sink.

```bash
export NG_PYTHON_LOG_SINK=stdout
```

### Set environment variables in python code

```python
import os

os.environ["NG_PYTHON_DEBUG"] = "true"
os.environ["NG_PYTHON_LOG_LEVEL"] = "DEBUG"
os.environ["NG_PYTHON_LOG_SINK"] = "stdout"
```

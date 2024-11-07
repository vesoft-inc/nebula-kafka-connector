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
    USE movie
    MATCH p=(a:Movie{name:"Unpromised Land"})-[e:WithGenre]->(b:Genre) 
    RETURN p as path, e as edge_WithGenre, b as genre_node, a.name as movie_name, 3.14 as float_val, true as bool_val
    LIMIT 2
"""
# Execute query
result = client.execute(query)
```

2. Then we could inspect the result ourselves.

```python
# Get one row
result = client.execute(query)
row = result.next()

# Get column names
row.column_names

# Get one value
cell = row.col_values[0].cast()

# see a cell and its methods
cell

# see all methods of a cell
dir(cell)
```

3. We could actually get primitive values from a cell with `cast_primitive()` method.

```python
cell_primitive = cell.cast_primitive()

cell_primitive
```

And it looks like this:

```python
{'nodes': [{'id': 289226172909223937,
   'type': 'Movie',
   'labels': ['Movie'],
   'properties': {'name': 'Unpromised Land', 'id': 91}},
  {'id': 289483299716333572,
   'type': 'Genre',
   'labels': ['Genre'],
   'properties': {'name': 'Staged Documentary', 'id': 101}}],
 'edges': [{'src_id': 289226172909223937,
   'dst_id': 289483299716333572,
   'rank': 0,
   'type': 'WithGenre',
   'labels': ['WithGenre'],
   'properties': {},
   'direction': 'OUTGOING'}],
 'length': 2,
 'start_node': {'id': 289226172909223937,
  'type': 'Movie',
  'labels': ['Movie'],
  'properties': {'name': 'Unpromised Land', 'id': 91}},
 'end_node': {'id': 289483299716333572,
  'type': 'Genre',
  'labels': ['Genre'],
  'properties': {'name': 'Staged Documentary', 'id': 101}},
 'string_representation': '(289226172909223937@Movie:Movie{name:Unpromised Land,id:91})-[0@WithGenre:WithGenre{}]->(289483299716333572@Genre:Genre{name:Staged Documentary,id:101})'}

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

5. Also you could print the result directly.

```python
client.print_query_result(query) # or client.pq(query)
```

Result:

```
╔════════════════════════════════════════════════╤════════════════════════════════════════════════╤═════════════════════════════════════════════════╤═════════════════╤═══════════╤══════════╗
║                                                │                                                │                                                 │                 │           │          ║
║ path                                           │ edge_WithGenre                                 │ genre_node                                      │ movie_name      │ float_val │ bool_val ║
║                                                │                                                │                                                 │                 │           │          ║
╟────────────────────────────────────────────────┼────────────────────────────────────────────────┼─────────────────────────────────────────────────┼─────────────────┼───────────┼──────────╢
║                                                │                                                │                                                 │                 │           │          ║
║ (289226172909223937@Movie:Movie{name:Unpromis… │ (289226172909223937)-[0@WithGenre:['WithGenre… │ (289556378584875011@Genre:['Genre']{name:Soci,… │ Unpromised Land │ 3.14      │ True     ║
║ Land,id:91})-[0@WithGenre:WithGenre{}]->(2895… │                                                │                                                 │                 │           │          ║
║                                                │                                                │                                                 │                 │           │          ║
║                                                │                                                │                                                 │                 │           │          ║
║ (289226172909223937@Movie:Movie{name:Unpromis… │ (289226172909223937)-[0@WithGenre:['WithGenre… │ (289483299716333572@Genre:['Genre']{name:Staged │ Unpromised Land │ 3.14      │ True     ║
║ Land,id:91})-[0@WithGenre:WithGenre{}]->(2894… │                                                │ Documentary,id:101})                            │                 │           │          ║
║ Documentary,                                   │                                                │                                                 │                 │           │          ║
║                                                │                                                │                                                 │                 │           │          ║
╚════════════════════════════════════════════════╧════════════════════════════════════════════════╧═════════════════════════════════════════════════╧═════════════════╧═══════════╧══════════╝

Summary
├── Rows: 2
└── Latency: 3422μs
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

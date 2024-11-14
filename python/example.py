from nebulagraph_python.client import NebulaClient

# Create client
client = NebulaClient(
    hosts=["192.168.8.148:9669"],
    username="root",
    password="Nebula@123",
)
if client._session is None:
    raise RuntimeError("Failed to create session")
session = client._session

query = """
USE movie
MATCH p=(a:Movie{name:"Unpromised Land"})-[e:WithGenre]->(b:Genre) 
RETURN p as path, e as edge_WithGenre, b as genre_node, a.name as movie_name, 3.14 as float_val, true as bool_val
LIMIT 2
"""
# Execute query
result = session.execute(query)

# Print results
result.print()

# Convert to pandas DataFrame

result = session.execute(query)

df = result.as_pandas_df()
df.to_csv("query_result.csv", index=False)

# Use as a CLI tool

session.print_query_result(query)  # or client.pq(query)

# Get one row
result = session.execute(query)
row = result.next()

# Get column names
row.column_names

# Get one value
cell = row.col_values[0].cast()

# Cast to primitive
cell_primitive = row.col_values[0].cast_primitive()

######
# special value type example
######

query = """
RETURN local_datetime("2016-09-20T01:01:01", "%Y-%m-%dT%H:%M:%S") AS localdatetime,
       local_time("05:06:07.089", "%H:%M:%S") AS localtime,
       zoned_time("05:06:07.089 +08:00", "%H:%M:%S %Ez") AS zonetime,
       zoned_datetime("2016-09-20T01:01:01 +0800", "%Y-%m-%dT%H:%M:%S %z") AS zoneddatetime,
       date("Tue, 2016-09-20", "%a, %Y-%m-%d") AS d,
       RECORD {a: 1, b: true, c: "str literal"} AS record1,
       LIST [1, 2, 3, 4, 5] AS l,
       "str literal" AS str_literal
"""

session.print_query_result(query)  # or client.pq(query)


######
# embedding vector example
######

# FOR DDL, DML refer to ann.feature

# Query KNN

query = """
USE ann_test
MATCH (v:N1|N2)
ORDER BY vector_distance(vector<3, float>([1, 2, 3]), v.vec1) LIMIT 3
RETURN v.id as vid, v.vec1 as vec1
"""

session.pq(query)


######
# debug example
######

import os

os.environ["NG_PYTHON_DEBUG"] = "true"
os.environ["NG_PYTHON_LOG_LEVEL"] = "DEBUG"
os.environ["NG_PYTHON_LOG_SINK"] = "stdout"

client = NebulaClient(
    hosts=["192.168.8.148:9669"],
    username="root",
    password="Nebula@123",
)

query = """
RETURN 1
"""

session.pq(query)

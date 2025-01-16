import os

from nebulagraph_python.client import NebulaClient, NebulaPool, SessionConfig

os.environ["NG_PYTHON_DEBUG"] = "true"
os.environ["NG_PYTHON_LOG_LEVEL"] = "DEBUG"
os.environ["NG_PYTHON_LOG_SINK"] = "stdout"

# Create client
client = NebulaClient(
    hosts=["127.0.0.1:9119"],
    username="root",
    password="Nebula@123",
    session_config=SessionConfig(
        graph="movie",
        timezone="Asia/Shanghai",
        parameters={"a": "1", "b": "[1, 2, 3]"},
    ),
)

client.execute_py("RETURN $a, $b").print()
client.execute_py("SHOW CURRENT_SESSION").print()
client.execute_py("DESC GRAPH TYPE movie_type").print()


query = """
USE movie
MATCH p=(a:Movie{name:"Unpromised Land"})-[e:WithGenre]->(b:Genre) 
RETURN p as path, e as edge_WithGenre, b as genre_node, a.name as movie_name, 3.14 as float_val, true as bool_val
LIMIT 2
"""
# Execute query
result = client.execute(query)

# Print results
result.print()

# Convert to pandas DataFrame

result = client.execute(query)

df = result.as_pandas_df()
df.to_csv("query_result.csv", index=False)

# Use as a CLI tool

client.print_query_result(query)  # or client.pq(query)

# Get one row
result = client.execute(query)
row = result.one()

# Get column names
print(row.column_names)

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

client.print_query_result(query)  # or client.pq(query)


######
# embedding vector example
######

# FOR DDL, DML refer to ann.feature

# Query KNN

query = """
USE ann_test
MATCH (v:N1|N2)
ORDER BY vector_distance(vector<3, float>([1, 2, 3]), v.vec1) LIMIT 3
RETURN v, v.vec1 as vec1
"""

client.pq(query)

client.close()  # Explicitly close the client to release all resources

######
# pool example
######


pool = NebulaPool(
    hosts=["192.168.8.148:9669"],
    username="root",
    password="Nebula@123",
)

query = """
RETURN 1
"""

with pool.borrow() as client:
    client.execute(query).print()

######
# execute_py example
######

query = """
RETURN {{v1}} as v1, {{v2}} as v2, {{v3}} as v3
"""
args = {"v1": 1, "v2": "alice", "v3": [True, False, True]}

with pool.borrow() as client:
    res = client.execute_py(query, args)
    # get the first row in primitive type
    row = res.one().as_primitive()
    res.print()
    # assert the row is the same as the args, in python primitive type
    assert row == args
    # get the result in column-oriented primitive type
    print(res.as_primitive_by_column())
    # get the result in row-oriented primitive type
    print(list(res.as_primitive_by_row()))

pool.close()  # Explicitly close the pool to release all resources


package com.vesoft.nebula.connector.sink;

public class NebulaSinkConstant {
    public static String BATCH_INSERT_TEMPLATE = "INSERT %s `%s`(%s) VALUES %s";
    public static String VERTEX_VALUE_TEMPLATE = "%s: (%s)";
    public static String VERTEX_VALUE_TEMPLATE_WITH_POLICY = "%s(\"%s\"): (%s)";
    public static String ENDPOINT_TEMPLATE = "%s(\"%s\")";
    public static String EDGE_VALUE_WITHOUT_RANKING_TEMPLATE = "%s->%s: (%s)";
    public static String EDGE_VALUE_TEMPLATE = "%s->%s@%d: (%s)";
}

/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.util;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import java.util.List;
import java.util.Map;

public class NebulaUtils {

    public static String extractPropertyValue(Map<String, String> schema, Map.Entry<String,
            Object> entry) throws DataFormatException {
        String propName = entry.getKey();
        String value = String.valueOf(entry.getValue());
        switch (schema.get(propName)) {
            case "STRING":
                return mkString(value, "\"", "", "\"");
            case "DATE":
                return "date(\"" + value + "\")";
            case "LOCALDATETIME":
                return "localdatetime(\"" + value + "\")";
            case "LOCALTIME": {
                return "localtime(\"" + value + "\")";
            }
            case "DURATION": {
                return "duration(\"" + value + "\")";
            }
            default:
                return value;
        }
    }

    public static String extractIdValue(String dataType, String value) {
        switch (dataType) {
            case "STRING":
                return mkString(value, "\"", "", "\"");
            case "INT64":
                return value;
            default:
                throw new RuntimeException("data type " + dataType + " of edge source/target id " +
                        "is not supported.");
        }

    }

    public static String mkString(String value, String start, String sep, String end) {
        StringBuilder builder = new StringBuilder();
        boolean first = true;
        builder.append(start);
        for (char c : value.toCharArray()) {
            if (first) {
                builder.append(c);
                first = false;
            } else {
                builder.append(sep);
                builder.append(c);
            }
        }
        builder.append(end);
        return builder.toString();
    }

    public static String join(List<String> list, String sep) {
        StringBuilder builder = new StringBuilder();
        for (String value : list) {
            builder.append(value);
            builder.append(sep);
        }
        return builder.deleteCharAt(builder.length() - 1).toString();
    }
}

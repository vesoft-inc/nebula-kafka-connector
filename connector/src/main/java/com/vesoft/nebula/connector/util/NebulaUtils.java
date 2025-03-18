
package com.vesoft.nebula.connector.util;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import java.util.List;
import java.util.Map;

public class NebulaUtils {

    public static String extractPropertyValue(Map<String, String> schema,
                                              String propName,
                                              String value) {
        switch (schema.get(propName)) {
            case "STRING":
                return mkString(value, "\"", "", "\"");
            case "DATE":
                return "date(\"" + value + "\")";
            case "LOCAL DATETIME":
                return "local_datetime(\"" + value + "\")";
            case "LOCAL TIME": {
                return "local_time(\"" + value + "\")";
            }
            case "ZONED DATETIME":
                return "zoned_datetime(\"" + value + "\")";
            case "ZONED TIME":
                return "zoned_time(\"" + value + "\")";
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
                throw new RuntimeException("data type " + dataType + " of edge source/target id "
                        + "is not supported.");
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
        if (builder.length() > 0) {
            builder.deleteCharAt(builder.length() - 1);
        }
        return builder.toString();
    }

    public static boolean isNumeric(String str) {
        String newStr = str;
        if (str.startsWith("-")) {
            newStr = str.substring(1);
        }

        for (char c : newStr.toCharArray()) {
            if (!Character.isDigit(c)) {
                return false;
            }
        }
        return true;
    }

}

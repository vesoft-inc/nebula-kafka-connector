
package com.vesoft.nebula.connector.util;

import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class NebulaUtils {

    public static String extractPropertyValue(Map<String, String> schema,
                                              String propName,
                                              String value,
                                              String nullValue) {
        return extractValue(schema.get(propName), value, nullValue);
    }

    public static String extractValue(String dataType, String value, String nullValue) {
        if (value == null || value.equals(nullValue)) {
            return null;
        }
        if (dataType.equals("STRING")) {
            return mkString(StringEscapeUtil.escapeValue(value), "\"", "", "\"");
        }
        if (value.isEmpty()) {
            return null;
        }
        // process the list type
        if (dataType.startsWith("LIST")) {
            if (dataType.startsWith("LIST<STRING")) {
                StringBuilder sb = new StringBuilder();
                sb.append("LIST[");
                String  trimmedInput = value.replaceAll("^\\[|\\]$", "");
                Pattern pattern      = Pattern.compile("(['\"])((?:\\\\\\1|.)*?)\\1|([^,]+)");
                Matcher matcher      = pattern.matcher(trimmedInput);

                while (matcher.find()) {
                    if (matcher.group(1) != null) {
                        String ele = matcher.group(2)
                                .replace("\\" + matcher.group(1), matcher.group(1));
                        sb.append("\"")
                                .append(StringEscapeUtil.escapeValue(ele))
                                .append("\"").append(",");
                    } else {
                        sb.append("\"")
                                .append(StringEscapeUtil.escapeValue(matcher.group(3)))
                                .append("\"").append(",");
                    }
                }

                if (sb.length() > 5) {
                    sb.deleteCharAt(sb.length() - 1);
                }
                sb.append("]");
                return sb.toString();
            } else {
                return "LIST" + value;
            }
        }

        // process the vector type
        if (dataType.startsWith("VECTOR<")) {
            return dataType + "(" + value + ")";
        }
        if (dataType.startsWith("GEOGRAPHY")) {
            return "ST_GeogFromText(\"" + value + "\")";
        }
        // process other data type
        switch (dataType) {
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

    public static String mkString(String value, String start, String sep, String end) {
        StringBuilder builder = new StringBuilder();
        boolean       first   = true;
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

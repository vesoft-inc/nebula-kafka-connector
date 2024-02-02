/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.util;

import com.vesoft.nebula.connector.exceptions.DataFormatException;
import java.text.ParseException;
import java.text.SimpleDateFormat;


public class TimeUtils {
    static SimpleDateFormat dateFormat1 = new SimpleDateFormat("yyyy/MM/dd");
    static SimpleDateFormat dateFormat2 = new SimpleDateFormat("yyyy.MM.dd");
    static SimpleDateFormat dateFormat3 = new SimpleDateFormat("yyyy-MM-dd");
    static SimpleDateFormat timeFormat = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSSX");


    public static String convertDate(String originalDate) throws DataFormatException {
        if (originalDate == null) {
            return null;
        }
        String date = getSimpleDate(originalDate);
        if (date == null) {
            throw new DataFormatException(String.format("format of date %s is not valid for NebulaGraph.",
                    originalDate));
        } else {
            return "date(\"" + date + "\")";
        }
    }

    public static String convertLocalTime(String originalTime) throws DataFormatException {
        if (originalTime == null) {
            return null;
        }
        if (originalTime.split(":").length == 3) {
            return "localtime(\"" + originalTime + "\")";
        } else {
            throw new DataFormatException(String.format("format of time %s is not valid for NebulaGraph.",
                    originalTime));
        }
    }

    public static String convertLocalDatetime(String originalDateTime) throws DataFormatException {
        if (originalDateTime == null) {
            return null;
        }
        String date = null;
        String time = null;
        if (originalDateTime.contains("T")) {
            date = getSimpleDate(originalDateTime.split("T")[0]);
            time = originalDateTime.split("T")[1];
        }
        if (originalDateTime.contains(" ")) {
            date = getSimpleDate(originalDateTime.split(" ")[0]);
            time = originalDateTime.split(" ")[1];
        }

        if (date == null) {
            throw new DataFormatException(String.format("format of datetime %s is not valid for " +
                    "NebulaGraph.", originalDateTime));
        }

        String newDateTime = date + "T" + time;
        try {
            timeFormat.parse(newDateTime);
            return "localdatetime(\"" + newDateTime + "\")";
        } catch (ParseException e) {
            throw new DataFormatException(String.format("format of datetime %s is not valid for " +
                    "NebulaGraph.", originalDateTime));
        }
    }


    private static String getSimpleDate(String dateString) {
        String date = null;

        try {
            dateFormat1.setLenient(false);
            dateFormat1.parse(dateString);
            date = dateString.replace("/", "-");
        } catch (ParseException e) {
        }

        try {
            dateFormat2.setLenient(false);
            dateFormat2.parse(dateString);
            date = dateString.replace(".", "-");
        } catch (ParseException e) {
        }

        try {
            dateFormat3.setLenient(false);
            dateFormat3.parse(dateString);
            date = dateString;
        } catch (ParseException e) {
        }
        return date;
    }
}

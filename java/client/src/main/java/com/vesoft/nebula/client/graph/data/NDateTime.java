/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;


import com.vesoft.nebula.LocalDatetime;
import java.util.Objects;

public class NDateTime {
    private final LocalDatetime localDateTime;

    public NDateTime(LocalDatetime localDateTime) {
        this.localDateTime = localDateTime;
    }

    /**
     * @return utc datetime year
     */
    public short getYear() {
        return localDateTime.getYear();
    }

    /**
     * @return utc datetime month
     */
    public byte getMonth() {
        return localDateTime.getMonth();
    }

    /**
     * @return datetime day
     */
    public byte getDay() {
        return localDateTime.getDay();
    }

    /**
     * @return datetime hour
     */
    public byte getHour() {
        return localDateTime.getHour();
    }

    /**
     * @return datetime minute
     */
    public byte getMinute() {
        return localDateTime.getMinute();
    }

    /**
     * @return utc datetime second
     */
    public byte getSecond() {
        return localDateTime.getSec();
    }

    /**
     * @return utc datetime microsec
     */
    public int getMicrosec() {
        return localDateTime.getMicrosec();
    }

    @Override
    public String toString() {
        return String.format("%d-%02d-%02dT%02d:%02d:%02d.%06d",
                localDateTime.year, localDateTime.month, localDateTime.day,
                localDateTime.hour, localDateTime.minute, localDateTime.sec,
                localDateTime.microsec);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        NDateTime that = (NDateTime) o;
        return localDateTime.year == that.getYear()
                && localDateTime.month == that.getMonth()
                && localDateTime.day == that.getDay()
                && localDateTime.hour == that.getHour()
                && localDateTime.minute == that.getMinute()
                && localDateTime.sec == that.getSecond()
                && localDateTime.microsec == that.getMicrosec();
    }

    @Override
    public int hashCode() {
        return Objects.hash(localDateTime);
    }
}

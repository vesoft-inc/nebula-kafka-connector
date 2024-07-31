package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.client.graph.utils.ZoneOffsetUtil;
import com.vesoft.nebula.proto.common.LocalDatetime;
import com.vesoft.nebula.proto.common.ZonedDatetime;
import java.util.Objects;

public class NZonedDateTime {
    private final ZonedDatetime zonedDateTime;

    public NZonedDateTime(ZonedDatetime zonedDateTime) {
        this.zonedDateTime = zonedDateTime;
    }

    /**
     * @return utc datetime year
     */
    public int getYear() {
        return zonedDateTime.getYear();
    }

    /**
     * @return utc datetime month
     */
    public int getMonth() {
        return zonedDateTime.getMonth();
    }

    /**
     * @return datetime day
     */
    public int getDay() {
        return zonedDateTime.getDay();
    }

    /**
     * @return datetime hour
     */
    public int getHour() {
        return zonedDateTime.getHour();
    }

    /**
     * @return datetime minute
     */
    public int getMinute() {
        return zonedDateTime.getMinute();
    }

    /**
     * @return utc datetime second
     */
    public int getSecond() {
        return zonedDateTime.getSec();
    }

    /**
     * @return utc datetime microsec
     */
    public int getMicrosec() {
        return zonedDateTime.getMicrosec();
    }

    /**
     * @return zoned offset in seconds
     */
    public int getOffset() {
        return zonedDateTime.getOffset();
    }

    @Override
    public String toString() {
        return String.format("%d-%02d-%02dT%02d:%02d:%02d.%06d%s",
                zonedDateTime.getYear(), zonedDateTime.getMonth(), zonedDateTime.getDay(),
                zonedDateTime.getHour(), zonedDateTime.getMinute(), zonedDateTime.getSec(),
                zonedDateTime.getMicrosec(), ZoneOffsetUtil.buildOffset(zonedDateTime.getOffset()));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        NZonedDateTime that = (NZonedDateTime) o;
        return zonedDateTime.getYear() == that.getYear()
                && zonedDateTime.getMonth() == that.getMonth()
                && zonedDateTime.getDay() == that.getDay()
                && zonedDateTime.getHour() == that.getHour()
                && zonedDateTime.getMinute() == that.getMinute()
                && zonedDateTime.getSec() == that.getSecond()
                && zonedDateTime.getMicrosec() == that.getMicrosec();
    }

    @Override
    public int hashCode() {
        return Objects.hash(zonedDateTime, getOffset());
    }
}

package com.vesoft.nebula.driver.graph.data;

import com.vesoft.nebula.driver.graph.utils.ZoneOffsetUtil;
import com.vesoft.nebula.proto.common.ZonedTime;
import java.util.Objects;

public class NZonedTime {

    private final ZonedTime zonedTime;

    public NZonedTime(ZonedTime zonedTime) {
        this.zonedTime = zonedTime;
    }

    /**
     * @return utc Time hour
     */
    public int getHour() {
        return zonedTime.getHour();
    }

    /**
     * @return utc Time minute
     */
    public int getMinute() {
        return zonedTime.getMinute();
    }

    /**
     * @return utc Time second
     */
    public int getSecond() {
        return zonedTime.getSec();
    }

    /**
     * @return utc Time microsec
     */
    public int getMicrosec() {
        return zonedTime.getMicrosec();
    }

    /**
     * get zone offset in seconds
     *
     * @return offset
     */
    public int getOffset() {
        return zonedTime.getOffset();
    }


    @Override
    public String toString() {
        return String.format("%02d:%02d:%02d.%06d%s",
                zonedTime.getHour(), zonedTime.getMinute(),
                zonedTime.getSec(), zonedTime.getMicrosec(),
                ZoneOffsetUtil.buildOffset(zonedTime.getOffset()));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        NZonedTime that = (NZonedTime) o;
        return zonedTime.getHour() == that.getHour()
                && zonedTime.getHour() == that.getMinute()
                && zonedTime.getSec() == that.getSecond()
                && zonedTime.getMicrosec() == that.getMicrosec();
    }

    @Override
    public int hashCode() {
        return Objects.hash(zonedTime, getOffset());
    }
}

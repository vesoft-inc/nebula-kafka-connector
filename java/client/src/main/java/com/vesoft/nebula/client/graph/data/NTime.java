/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.LocalTime;
import java.util.Objects;

public class NTime {
    private final LocalTime localTime;

    public NTime(LocalTime localTime) {
        this.localTime = localTime;
    }

    /**
     * @return utc Time hour
     */
    public byte getHour() {
        return localTime.getHour();
    }

    /**
     * @return utc Time minute
     */
    public byte getMinute() {
        return localTime.getMinute();
    }

    /**
     * @return utc Time second
     */
    public byte getSecond() {
        return localTime.getSec();
    }

    /**
     * @return utc Time microsec
     */
    public int getMicrosec() {
        return localTime.getMicrosec();
    }


    @Override
    public String toString() {
        return String.format("%02d:%02d:%02d.%06d",
                localTime.hour, localTime.minute, localTime.sec, localTime.microsec);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        NTime that = (NTime) o;
        return localTime.hour == that.getHour()
                && localTime.minute == that.getMinute()
                && localTime.sec == that.getSecond()
                && localTime.microsec == that.getMicrosec();
    }

    @Override
    public int hashCode() {
        return Objects.hash(localTime);
    }
}

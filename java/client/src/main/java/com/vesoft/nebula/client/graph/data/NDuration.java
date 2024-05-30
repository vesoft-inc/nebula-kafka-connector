/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.common.Duration;
import java.util.Objects;

public class NDuration {
    private final Duration duration;

    public NDuration(Duration duration) {
        this.duration = duration;
    }


    /**
     * @return whether the duration is month-based duration type
     */
    public boolean isMonthBased() {
        return duration.getIsMonthBased();
    }

    /**
     * @return duration year
     */
    public int getYear() {
        return duration.getYear();
    }

    /**
     * @return duration month
     */
    public int getMonth() {
        return duration.getMonth();
    }

    /**
     * @return duration year
     */
    public int getDay() {
        return duration.getDay();
    }

    /**
     * @return duration hour
     */
    public int getHour() {
        return duration.getHour();
    }

    /**
     * @return duration minute
     */
    public int getMinute() {
        return duration.getMinute();
    }

    /**
     * @return duration seconds
     */
    public int getSecond() {
        return duration.getSec();
    }

    /**
     * @return duration microseconds
     */
    public int getMicrosecond() {
        return duration.getMicrosec();
    }

    @Override
    public String toString() {
        StringBuilder durationStr = new StringBuilder();
        durationStr.append("P");
        if (isMonthBased()) {
            if (getYear() != 0) {
                durationStr.append(getYear()).append("Y");
            }
            if (getMonth() != 0) {
                durationStr.append(getMonth()).append("M");
            }
        } else {
            if (getDay() != 0) {
                durationStr.append(getDay()).append("D");
            }
            if (getHour() != 0 || getMinute() != 0 || getSecond() != 0 || getMicrosecond() != 0) {
                durationStr.append("T");
            }
            if (getHour() != 0) {
                durationStr.append(getHour()).append("H");
            }
            if (getMinute() != 0) {
                durationStr.append(getMinute()).append("M");
            }
            if (getSecond() != 0 || getMicrosecond() != 0) {
                if (getMicrosecond() == 0) {
                    durationStr.append(getSecond()).append("S");
                } else {
                    int totalMicroseconds = getSecond() * 1000000 + getMicrosecond();
                    boolean isMinus = totalMicroseconds < 0;
                    if (isMinus) {
                        totalMicroseconds = -totalMicroseconds;
                    }
                    int seconds = totalMicroseconds / 1000000;
                    int microSec = totalMicroseconds % 1000000;
                    if (isMinus) {
                        durationStr
                                .append(String.format("-%d.%06d", seconds, microSec));
                    } else {
                        durationStr
                                .append(String.format("%d.%06d", seconds, microSec));
                    }
                    while (durationStr.charAt(durationStr.length() - 1) == '0') {
                        durationStr.setLength(durationStr.length() - 1);
                    }
                    durationStr.append("S");
                }
            }
        }
        return durationStr.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        NDuration that = (NDuration) o;
        return duration.getIsMonthBased() == that.isMonthBased()
                && duration.getYear() == that.getYear()
                && duration.getMonth() == that.getMonth()
                && duration.getDay() == that.getDay()
                && duration.getHour() == that.getHour()
                && duration.getMinute() == that.getMinute()
                && duration.getSec() == that.getSecond()
                && duration.getMicrosec() == that.getMicrosecond();
    }

    @Override
    public int hashCode() {
        return Objects.hash(duration);
    }
}

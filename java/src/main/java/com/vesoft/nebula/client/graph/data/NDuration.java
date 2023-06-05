/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.Duration;
import java.util.Objects;

public class NDuration {
    private final Duration duration;

    public NDuration(Duration duration) {
        this.duration = duration;
    }

    /**
     * @return duration months
     */
    public int getMonths() {
        return duration.getMonths();
    }

    /**
     * @return duration seconds
     */
    public long getSeconds() {
        return duration.getSeconds();
    }

    /**
     * @return duration microseconds
     */
    public int getMicroseconds() {
        return duration.getMicroseconds();
    }

    @Override
    public String toString() {
        long totalSeconds = duration.seconds + duration.microseconds / 1000000;
        int remainMicroSeconds = duration.microseconds % 1000000;
        String microSends = String.format("%06d", remainMicroSeconds) + "000";
        return String.format("P%dMT%d.%sS", duration.months, totalSeconds, microSends);
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
        return duration.months == that.getMonths()
                && duration.seconds == that.getSeconds()
                && duration.microseconds == that.getMicroseconds();
    }

    @Override
    public int hashCode() {
        return Objects.hash(duration);
    }
}

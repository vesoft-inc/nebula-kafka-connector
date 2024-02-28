/* Copyright (c) 2023 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.data;

import com.vesoft.nebula.proto.graph.Duration;
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
        long totalSeconds = duration.getSeconds() + duration.getMicroseconds() / 1000000;
        int remainMicroSeconds = duration.getMicroseconds() % 1000000;
        String microSends = String.format("%06d", remainMicroSeconds) + "000";
        return String.format("P%dMT%d.%sS", duration.getMonths(), totalSeconds, microSends);
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
        return duration.getMonths() == that.getMonths()
                && duration.getSeconds() == that.getSeconds()
                && duration.getMicroseconds() == that.getMicroseconds();
    }

    @Override
    public int hashCode() {
        return Objects.hash(duration);
    }
}

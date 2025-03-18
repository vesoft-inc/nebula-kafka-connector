/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.util;

import java.util.ArrayList;
import java.util.List;

public class GroupUtils {

    public static <T> List<List<T>> getGroups(List<T> list, int size) {
        List<List<T>> result = new ArrayList<List<T>>();
        if (list == null || list.isEmpty()) {
            return result;
        }
        int remainder = list.size() % size;
        int amount = list.size() / size;
        if (remainder > 0) {
            amount++;
        }
        for (int i = 0; i < amount; i++) {
            List<T> value;
            if (i == (amount - 1)) {
                if (remainder == 0) {
                    value = list.subList(i * size, (i + 1) * size);
                } else {
                    value = list.subList(i * size, i * size + remainder);
                }
            } else {
                value = list.subList(i * size, i * size + size);
            }
            result.add(value);
        }
        return result;
    }
}

/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode;

import static com.vesoft.nebula.driver.graph.decode.DecodeUtils.bytesToInt16;
import static com.vesoft.nebula.driver.graph.decode.struct.SizeConstant.ELEMENT_NUMBER_SIZE_FOR_ANY_VALUE;

import com.google.protobuf.ByteString;
import java.nio.ByteOrder;

public class BytesReader {
    private ByteString data;
    private int        index;

    public BytesReader(ByteString data) {
        this.data = data;
        this.index = 0;
    }

    public ByteString read(int len) {
        ByteString byteString = data.substring(index, index + len);
        index += len;
        return byteString;
    }

    public String readString() {
        StringBuilder sb = new StringBuilder();
        for (int propCharIndex = index; propCharIndex < data.size(); propCharIndex++) {
            if (data.byteAt(propCharIndex) == '\0') {
                index++;
                break;
            }
            sb.append((char) data.byteAt(propCharIndex));
        }
        index += sb.length();
        return sb.toString();
    }

    public String readSizedString(ByteOrder byteOrder) {
        int           length = bytesToInt16(read(ELEMENT_NUMBER_SIZE_FOR_ANY_VALUE), byteOrder);
        StringBuilder sb     = new StringBuilder();
        for (int propCharIndex = index; propCharIndex < index + length; propCharIndex++) {
            sb.append((char) data.byteAt(propCharIndex));
        }
        index += length;
        return sb.toString();
    }
}

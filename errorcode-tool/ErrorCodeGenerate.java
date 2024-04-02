/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Generator to generate the error code for client from Server ErrorCode
 * definition
 */
public class ErrorCodeGenerate {

    public static void main(String[] args) throws IOException {

        String codeFileName = args[0];
        String client = args[1];
        if (client == null) {
            client = "java";
        }
        Pattern pattern = Pattern.compile("DEFINE_ERRORCODE\\((.*), \"(.*)\", (.*?)\\),");

        File file = new File(codeFileName);
        FileReader fr = new FileReader(file);
        BufferedReader br = new BufferedReader(fr);
        String line;

        while ((line = br.readLine()) != null) {
            String error = line.trim();
            if (error.isEmpty() || error.startsWith("//")) {
                // System.out.println();
                continue;
            }

            Matcher matcher = pattern.matcher(error);
            if (matcher.find()) {
                String errorPrefix = getCodePrefix(matcher.group(1));
                if (errorPrefix == null) {
                    continue;
                }
                String errorCode = matcher.group(2);
                String errorName = matcher.group(3);
                String newCode = "";
                switch (client) {
                    case "java":
                        newCode = String.format("%s(\"%s\"),",
                                errorName, errorPrefix + errorCode);
                        break;
                    case "golang":
                        newCode = String.format("ERROR_%s ErrorCode = \"%s\"",
                                errorName, errorPrefix + errorCode);
                        break;
                }
                System.out.println(newCode);
            } else {
                System.out.println("===== cannot parse the error code, " + error);
            }
        }
    }

    public static String getCodePrefix(String code) {
        switch (code.trim()) {
            case "SUCCESSFUL_COMPLETION":
                return "00";
            case "WARNING":
                return "01";
            case "NO_DATA":
                return "02";
            case "INFORMATIONAL":
                return "03";
            case "CONNECTION_EXCEPTION":
                return "08";
            case "DATA_EXCEPTION":
                return "22";
            case "INVALID_TRANSACTION_STATE":
                return "25";
            case "INVALID_TRANSACTION_TERMINATION":
                return "2D";
            case "TRANSACTION_ROLLBACK":
                return "40";
            case "SYNTAX_ERROR_OR_ACCESS_RULE_VIOLATION":
                return "42";
            case "DEPENDENT_OBJECT_ERROR":
                return "G1";
            case "GRAPH_TYPE_VIOLATION":
                return "G2";

            /* Extended */
            /* exceptions */
            case "RUNTIME_ERROR":
                return "NR";
            case "SEMANTIC_ERROR":
                return "NS";
            case "STORAGE_ERROR":
                return "NO";
            case "CATALOG_ERROR":
                return "NC";
            case "DATA_READ_WRITE_ERROR":
                return "ND";
            case "INVALID_PARAMETER":
                return "NI";
            case "GRAPH_COMPUTE_ERROR":
                return "NG";
            case "PLUGIN_ERROR":
                return "NP";
            case "SESSION_ERROR":
                return "NE";
            case "UNSUPPORTED":
                return "NT";
            case "AUTHENTICATE_ERROR":
                return "NH";
            case "JOB_ERROR":
                return "NJ";
            case "LICENSE_ERROR":
                return "NL";
            /* internal errors */
            /*
             * case "METADATA_ERROR":
             * return "NM";
             * case "RPC_ERROR":
             * return "NN";
             * case "KVSTORE_ERROR":
             * return "NK";
             * case "RAFT_ERROR":
             * return "NA";
             * case "SYSTEM_ERROR":
             * return "NY";
             * case "JOB_ERROR":
             * return "NJ";
             * case "LICENSE_ERROR":
             * return "NL";
             * case "UNKNOWN":
             * return "NU";
             */
            default:
                return null;
        }
    }

    public static boolean isInternal(String code) {
        if (code.startsWith("NM")
                || code.startsWith("NN")
                || code.startsWith("NK")
                || code.startsWith("NA")
                || code.startsWith("NY")
                || code.startsWith("NJ")
                || code.startsWith("NL")
                || code.startsWith("NU")) {
            return true;
        }
        return false;
    }
}

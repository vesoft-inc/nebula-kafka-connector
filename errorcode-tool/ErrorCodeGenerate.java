/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Generator to generate the error code for client from Server ErrorCode
 * definition
 */
public class ErrorCodeGenerate {
    public static void main(String[] args) throws IOException {
        String codeFileName = args[0];
        String codeclassFileName = args[1];
        String messageFileName = args[2];
        String docDescFileName = args[3];
        List<ErrorCode> codes = constructErrorCode(codeFileName, codeclassFileName, messageFileName);
        Map<String, CodeDesc> codeDesc = getCodeDesc(docDescFileName);
        List<ErrorCode> codesWithInternalError = update(codes);
        writeCodeForJava(codesWithInternalError);
        writeCodeForGo(codesWithInternalError);
        writeCodeForPython(codesWithInternalError);
        writeCodeForDoc(codesWithInternalError, codeDesc);
        writeCodeForYaml(codesWithInternalError, codeDesc);
        System.out.println("Finished.");
    }


    public static void writeCodeForJava(List<ErrorCode> codes) {
        String javaCodeFile = "ErrorCode.java";
        FileWriter writer = null;
        try {
            writer = new FileWriter(javaCodeFile);
            writer.write(javaPrefix);
            for (ErrorCode code : codes) {
                String codeString = String.format("    %s(\"%s\"),\n", code.getName(), code.getCode());
                writer.write(codeString);
            }
            writer.write(javaSuffix);
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (writer != null) {
                try {
                    writer.flush();
                    writer.close();
                } catch (Exception e) {
                    // ignore
                }
            }
        }
    }


    public static void writeCodeForGo(List<ErrorCode> codes) {
        String javaCodeFile = "error.go";
        FileWriter writer = null;
        try {
            writer = new FileWriter(javaCodeFile);
            writer.write(goPrefix);
            for (ErrorCode code : codes) {
                String codeString = String.format("ERROR_%s ErrorCode = \"%s\"\n", code.getName(), code.getCode());
                writer.write(codeString);
            }
            writer.write(goSuffix);
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (writer != null) {
                try {
                    writer.flush();
                    writer.close();
                } catch (Exception e) {
                    // ignore
                }
            }
        }
    }


    public static void writeCodeForPython(List<ErrorCode> codes) {
        String pythonCodeFile = "_error_code.py";
        FileWriter writer = null;
        try {
            writer = new FileWriter(pythonCodeFile);
            writer.write(pythonPrefix);

            for (ErrorCode code : codes) {
                String codeString = String.format("    %s = '%s'\n", code.getName(), code.getCode());
                writer.write(codeString);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (writer != null) {
                try {
                    writer.flush();
                    writer.close();
                } catch (Exception e) {
                    // ignore
                }
            }
        }
    }


    public static void writeCodeForDoc(List<ErrorCode> codes, Map<String, CodeDesc> codeDescMap) {
        String enDescYaml = "errorcode_doc_en.md";
        String chDescYaml = "errorcode_doc_ch.md";

        StringBuilder englishDesc = new StringBuilder();
        StringBuilder chineseDesc = new StringBuilder();

        englishDesc.append("| Code | Message | Description|\n");
        englishDesc.append("|-----|-----|-----|\n");

        chineseDesc.append("| 错误码 | 错误信息 | 描述 |\n");
        chineseDesc.append("|-----|-----|-----|\n");
        for (ErrorCode code : codes) {
            String desc = null;
            if (codeDescMap.get(code.getCode()) != null) {
                desc = codeDescMap.get(code.getCode()).getEnglishDesc();
            }
            englishDesc.append("|").append(code.getCode())
                    .append("|").append(code.getMessage())
                    .append("|").append(desc).append("|\n");

            desc = null;
            if (codeDescMap.get(code.getCode()) != null) {
                desc = codeDescMap.get(code.getCode()).getChineseDes();
            }
            chineseDesc.append("|").append(code.getCode())
                    .append("|").append(code.getMessage())
                    .append("|").append(desc).append("|\n");
        }
        FileWriter englishDescWriter = null;
        FileWriter chineseDescWriter = null;
        try {
            englishDescWriter = new FileWriter(enDescYaml);
            englishDescWriter.write(englishDesc.toString());

            chineseDescWriter = new FileWriter(chDescYaml);
            chineseDescWriter.write(chineseDesc.toString());
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (englishDescWriter != null) {
                try {
                    englishDescWriter.flush();
                    englishDescWriter.close();
                } catch (Exception e) {
                    // ignore
                }
            }
            if (chineseDescWriter != null) {
                try {
                    chineseDescWriter.flush();
                    chineseDescWriter.close();
                } catch (Exception e) {
                    // ignore
                }
            }
        }
    }


    public static void writeCodeForYaml(List<ErrorCode> codes, Map<String, CodeDesc> codeDescMap) {
        String codeYaml = "errorcode.yaml";

        StringBuilder sb = new StringBuilder();
        for (ErrorCode code : codes) {
            sb.append("- Name: \"").append(code.getName()).append("\"\n");
            sb.append("  Code: \"").append(code.getCode()).append("\"\n");
            sb.append("  Message: \"").append(code.getMessage()).append("\"\n\n");
        }
        FileWriter codeWriter = null;
        try {
            codeWriter = new FileWriter(codeYaml);
            codeWriter.write(sb.toString());

        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (codeWriter != null) {
                try {
                    codeWriter.flush();
                    codeWriter.close();
                } catch (Exception e) {
                    // ignore
                }
            }

        }
    }


    public static List<ErrorCode> constructErrorCode(String codePath, String codeclassPath, String messageFilePath) {

        Map<String, String> codeClassPrefix = getCodePrefix(codeclassPath);
        Map<String, String> codeNameMessage = getCodeMessage(messageFilePath);
        return getErrorCodes(codePath, codeClassPrefix, codeNameMessage);
    }


    public static List<ErrorCode> getErrorCodes(String filePath,
                                                Map<String, String> codeClassPrefix,
                                                Map<String, String> codeNameMessage) {
        List<ErrorCode> errorCodes = new ArrayList<>();
        try {
            // read errorcode file and construct ErrorCode struct
            Pattern pattern = Pattern.compile("DEFINE_ERRORCODE\\((.*), \"(.*)\", (.*?)\\),");

            File file = new File(filePath);
            FileReader fr = new FileReader(file);
            BufferedReader br = new BufferedReader(fr);
            String line;
            while ((line = br.readLine()) != null) {
                String error = line.trim();
                if (error.isEmpty() || error.startsWith("//")) {
                    continue;
                }
                Matcher matcher = pattern.matcher(error);
                if (matcher.find()) {
                    String codeClass = matcher.group(1);
                    String errorSubCode = matcher.group(2);
                    String errorName = matcher.group(3);
                    ErrorCode code = new ErrorCode(errorName, codeClassPrefix.get(codeClass), errorSubCode, codeNameMessage.get(errorName));
                    errorCodes.add(code);
                } else {
                    System.out.println("===== cannot parse the error code, " + error);
                }
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        return errorCodes;
    }


    public static Map<String, String> getCodePrefix(String filePath) {
        Map<String, String> codeClassPrefix = new HashMap<>();
        Pattern pattern = Pattern.compile("(.*)= CLASS\\(\"(.*)\"\\),");
        try {
            BufferedReader br = new BufferedReader(new FileReader(filePath));
            String line;
            while ((line = br.readLine()) != null) {
                if (!line.contains("= CLASS")) {
                    continue;
                }
                String codeClass = line.trim();
                Matcher matcher = pattern.matcher(codeClass);

                if (matcher.matches()) {
                    String key = matcher.group(1).trim();
                    String value = matcher.group(2).trim();
                    codeClassPrefix.put(key, value);
                }
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return codeClassPrefix;
    }


    public static Map<String, String> getCodeMessage(String filePath) {
        Map<String, String> errorCodeMsg = new HashMap<>();
        Pattern pattern = Pattern.compile("EMSG\\((.*), \"(.*)\"\\);");
        try {
            BufferedReader br = new BufferedReader(new FileReader(filePath));
            String line;
            while ((line = br.readLine()) != null) {
                if (!line.trim().startsWith("EMSG")) {
                    continue;
                }
                String codeMsg = line.trim();
                Matcher matcher = pattern.matcher(codeMsg);
                if (matcher.find()) {
                    String name = matcher.group(1);
                    String msg = matcher.group(2);
                    errorCodeMsg.put(name, msg);
                }
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return errorCodeMsg;
    }

    public static Map<String, CodeDesc> getCodeDesc(String filePath) {
        Map<String, CodeDesc> codeDescMap = new HashMap<>();
        try {
            BufferedReader br = new BufferedReader(new FileReader(filePath));
            String line;
            while ((line = br.readLine()) != null) {
                if (!line.trim().startsWith("|") || line.trim().startsWith("| Code") || line.trim().startsWith("|----")) {
                    continue;
                }

                String[] codeDesc = line.trim().split("\\|");
                codeDescMap.put(codeDesc[1].trim(), new CodeDesc(codeDesc[1].trim(), codeDesc[2].trim(), codeDesc[3].trim()));
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return codeDescMap;
    }

    public static List<ErrorCode> update(List<ErrorCode> codes) {
        List<ErrorCode> filterCodes = new ArrayList<>();
        for (ErrorCode code : codes) {
            if (isInternal(code)) {
                code.setMessage("Internal Server Error");
            }
            filterCodes.add(code);
        }
        return filterCodes;

    }

    public static boolean isInternal(ErrorCode code) {
        String prefixCode = code.getPrefixCode();
        if (prefixCode.startsWith("NM")
                || prefixCode.startsWith("NN")
                || prefixCode.startsWith("NK")
                || prefixCode.startsWith("NA")
                || prefixCode.startsWith("NY")
                || prefixCode.startsWith("NW")
                || prefixCode.startsWith("NU")) {
            return true;
        }
        return false;
    }


    static class ErrorCode {
        private String name;
        private String prefixCode;
        private String subCode;
        private String code; // prefixCode + subCode

        private String message;

        ErrorCode(String name, String prefixCode, String subCode, String msg) {
            this.name = name;
            this.prefixCode = prefixCode;
            this.subCode = subCode;
            this.code = prefixCode + subCode;
            this.message = msg;
        }

        public void setName(String name) {
            this.name = name;
        }

        public void setCode(String code) {
            this.code = code;
        }

        public void setMessage(String message) {
            this.message = message;
        }

        public String getName() {
            return name;
        }

        public String getCode() {
            return code;
        }

        public String getMessage() {
            return message;
        }

        public String getPrefixCode() {
            return prefixCode;
        }
    }

    static class CodeDesc {
        private String code;
        private String englishDesc;
        private String chineseDes;

        CodeDesc(String code, String englishDesc, String chineseDes) {
            this.code = code;
            this.englishDesc = englishDesc;
            this.chineseDes = chineseDes;
        }

        public String getCode() {
            return code;
        }

        public String getEnglishDesc() {
            return englishDesc;
        }

        public String getChineseDes() {
            return chineseDes;
        }
    }


    private static String javaPrefix = "package com.vesoft.nebula.driver.graph;\n" +
            "\n" +
            "import java.util.HashMap;\n" +
            "import java.util.Map;\n" +
            "\n" +
            "public enum ErrorCode {\n";
    private static String javaSuffix = "\n\n" +
            "    UNKNOWN_FOR_CLIENT(null);\n" +
            "\n" +
            "\n" +
            "    public String code;\n" +
            "\n" +
            "    ErrorCode(String c) {\n" +
            "        code = c;\n" +
            "    }\n" +
            "\n" +
            "    private static Map<String, ErrorCode> errorCodeMap = new HashMap<>();\n" +
            "\n" +
            "    static {\n" +
            "        for (ErrorCode value : ErrorCode.values()) {\n" +
            "            errorCodeMap.put(value.code, value);\n" +
            "        }\n" +
            "    }\n" +
            "\n" +
            "    public static ErrorCode find(String code) {\n" +
            "\n" +
            "        if (errorCodeMap.containsKey(code)) {\n" +
            "            return errorCodeMap.get(code);\n" +
            "        } else {\n" +
            "            UNKNOWN_FOR_CLIENT.code = code;\n" +
            "            return UNKNOWN_FOR_CLIENT;\n" +
            "        }\n" +
            "    }\n" +
            "\n" +
            "    public boolean isRetryable() {\n" +
            "        return isSessionError() || isRpcError() || isRaftError();\n" +
            "    }\n" +
            "\n" +
            "    public boolean isSemanticError() {\n" +
            "        return code.startsWith(\"NS\");\n" +
            "    }\n" +
            "\n" +
            "    public boolean isSyntaxError() {\n" +
            "        return code.startsWith(\"42\");\n" +
            "    }\n" +
            "\n" +
            "    public boolean isSessionError() {\n" +
            "        return code.startsWith(\"NE\");\n" +
            "    }\n" +
            "\n" +
            "    public boolean isRpcError() {\n" +
            "        return code.startsWith(\"NN\");\n" +
            "    }\n" +
            "\n" +
            "\n" +
            "    public boolean isRaftError() {\n" +
            "        return code.startsWith(\"NA\");\n" +
            "    }\n" +
            "\n" +
            "    public boolean isNoDataError() {\n" +
            "        return code.startsWith(\"02\");\n" +
            "    }\n" +
            "\n" +
            "}\n";


    private static String goPrefix = "package errors\n" +
            "\n" +
            "import (\n" +
            "\t\"fmt\"\n" +
            "\n" +
            "\tgoerr \"github.com/pkg/errors\"\n" +
            ")\n" +
            "\n" +
            "type ErrorCode string\n" +
            "\n" +
            "// TODO add error code in future\n" +
            "type NebulaError struct {\n" +
            "\terr       error\n" +
            "\terrorCode ErrorCode\n" +
            "\terrorMsg  string\n" +
            "}\n" +
            "\n" +
            "type formater interface {\n" +
            "\tFormat(s fmt.State, verb rune)\n" +
            "}\n" +
            "\n" +
            "var (\n" +
            "\t// Error in client side\n" +
            "\t// Tmp error code, would be redefined in future\n" +
            "\tERROR_ADDRESS_NOT_VALID    ErrorCode = \"99000\"\n" +
            "\tERROR_CANNOT_OPEN          ErrorCode = \"99001\"\n" +
            "\tERROR_CONN_IS_BROKEN       ErrorCode = \"99002\"\n" +
            "\tERROR_CONN_CONNECT_TIMEOUT ErrorCode = \"99003\"\n" +
            "\tERROR_CONN_REQUEST_TIMEOUT ErrorCode = \"99004\"\n" +
            "\tERROR_CONN_IS_CLOSED       ErrorCode = \"99005\"\n" +
            "\tERROR_WAIT_POOL_TIMEOUT    ErrorCode = \"99006\"\n" +
            "\tERROR_ILLEGAL              ErrorCode = \"99007\"\n" +
            "\tERROR_TYPE                 ErrorCode = \"99008\"\n" +
            "\tERROR_CLIENT_INTERNEL      ErrorCode = \"99009\" // TODO should be removed\n" +
            "\tERROR_CLIENT_INTERNAL      ErrorCode = \"99009\"\n" +
            "\tERROR_TLS_ERROR            ErrorCode = \"99010\"\n" +
            "\tERROR_UNKNOWN_COLUMN_TYPE  ErrorCode = \"99011\"\n" +
            "\n" +
            "\t// Error in server side\n" +
            "\n";
    private static String goSuffix = ")\n" +
            "\n" +
            "func NewNebulaError(code ErrorCode, format string, args ...interface{}) error {\n" +
            "\tr := &NebulaError{\n" +
            "\t\terr:       goerr.New(fmt.Sprintf(format, args...)),\n" +
            "\t\terrorCode: code,\n" +
            "\t\terrorMsg:  fmt.Sprintf(format, args...),\n" +
            "\t}\n" +
            "\treturn r\n" +
            "}\n" +
            "\n" +
            "func (e *NebulaError) Error() string {\n" +
            "\treturn fmt.Sprintf(\"[%s]: %s\", e.errorCode, e.errorMsg)\n" +
            "}\n" +
            "\n" +
            "func (e *NebulaError) Code() ErrorCode {\n" +
            "\treturn e.errorCode\n" +
            "}\n" +
            "\n" +
            "func (e *NebulaError) ErrorClass() string {\n" +
            "\treturn string(e.errorCode[:2])\n" +
            "}\n" +
            "\n" +
            "func (e *NebulaError) Format(s fmt.State, verb rune) {\n" +
            "\tf := e.err.(formater)\n" +
            "\tf.Format(s, verb)\n" +
            "}\n" +
            "\n" +
            "func Wrap(err error, msg string) error {\n" +
            "\tnbErr, ok := err.(*NebulaError)\n" +
            "\tif !ok {\n" +
            "\t\treturn goerr.Wrap(err, msg)\n" +
            "\t}\n" +
            "\tnbErr.err = goerr.Wrap(nbErr.err, msg)\n" +
            "\tnbErr.errorMsg = fmt.Sprintf(\"%s: %s\", nbErr.errorMsg, msg)\n" +
            "\treturn nbErr\n" +
            "}\n";

    private static String pythonPrefix = "# Generated by ErrorCodeGenerate. DO NOT MANUALLY EDIT.\n" +
            "\n" +
            "from enum import Enum\n" +
            "\n" +
            "class ErrorCode(Enum):\n" +
            "    @classmethod\n" +
            "    def _missing_(cls, value):\n" +
            "        return cls.UNKNOWN\n" +
            "\n";
}

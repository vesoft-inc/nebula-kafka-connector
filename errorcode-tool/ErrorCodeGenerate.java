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
        String                codeFileName           = args[0];
        String                codeclassFileName      = args[1];
        String                messageFileName        = args[2];
        String                docDescFileName        = args[3];
        List<ErrorCode>       codes                  = constructErrorCode(codeFileName, codeclassFileName, messageFileName);
        Map<String, CodeDesc> codeDesc               = getCodeDesc(docDescFileName);
        List<ErrorCode>       codesWithInternalError = update(codes);
        writeCodeForJava(codesWithInternalError);
        writeCodeForGo(codesWithInternalError);
        writeCodeForYaml(codesWithInternalError, codeDesc);
        System.out.println("Finished.");
    }


    public static void writeCodeForJava(List<ErrorCode> codes) {
        String     javaCodeFile = "errorcode_java.txt";
        FileWriter writer       = null;
        try {
            writer = new FileWriter(javaCodeFile);
            for (ErrorCode code : codes) {
                String codeString = String.format("%s(\"%s\"),\n", code.getName(), code.getCode());
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


    public static void writeCodeForGo(List<ErrorCode> codes) {
        String     javaCodeFile = "errorcode_golang.txt";
        FileWriter writer       = null;
        try {
            writer = new FileWriter(javaCodeFile);
            for (ErrorCode code : codes) {
                String codeString = String.format("ERROR_%s ErrorCode = \"%s\"\n", code.getName(), code.getCode());
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

    public static void writeCodeForYaml(List<ErrorCode> codes, Map<String, CodeDesc> codeDescMap) {
        String enDescYaml = "errorcode_doc_en.txt";
        String chDescYaml = "errorcode_doc_ch.txt";

        StringBuilder englishDesc = new StringBuilder();
        StringBuilder chineseDesc = new StringBuilder();
        for (ErrorCode code : codes) {
            englishDesc.append("- Name: ").append(code.getName()).append("\n");
            englishDesc.append("  Code: ").append(code.getCode()).append("\n");
            englishDesc.append("  Message: ").append(code.getMessage()).append("\n");
            String desc = null;
            if (codeDescMap.get(code.getCode()) != null) {
                desc = codeDescMap.get(code.getCode()).getEnglishDesc();
            }
            englishDesc.append("  Description: ").append(desc).append("\n");

            chineseDesc.append("- 错误名: ").append(code.getName()).append("\n");
            chineseDesc.append("  错误码: ").append(code.getCode()).append("\n");
            chineseDesc.append("  错误消息: ").append(code.getMessage()).append("\n");
            desc = null;
            if (codeDescMap.get(code.getCode()) != null) {
                desc = codeDescMap.get(code.getCode()).getChineseDes();
            }
            chineseDesc.append("  描述: ").append(desc).append("\n");
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

            File           file = new File(filePath);
            FileReader     fr   = new FileReader(file);
            BufferedReader br   = new BufferedReader(fr);
            String         line;
            while ((line = br.readLine()) != null) {
                String error = line.trim();
                if (error.isEmpty() || error.startsWith("//")) {
                    continue;
                }
                Matcher matcher = pattern.matcher(error);
                if (matcher.find()) {
                    String    codeClass    = matcher.group(1);
                    String    errorSubCode = matcher.group(2);
                    String    errorName    = matcher.group(3);
                    ErrorCode code         = new ErrorCode(errorName, codeClassPrefix.get(codeClass), errorSubCode, codeNameMessage.get(errorName));
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
        Pattern             pattern         = Pattern.compile("(.*)= CLASS\\(\"(.*)\"\\),");
        try {
            BufferedReader br = new BufferedReader(new FileReader(filePath));
            String         line;
            while ((line = br.readLine()) != null) {
                if (!line.contains("= CLASS")) {
                    continue;
                }
                String  codeClass = line.trim();
                Matcher matcher   = pattern.matcher(codeClass);

                if (matcher.matches()) {
                    String key   = matcher.group(1).trim();
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
        Pattern             pattern      = Pattern.compile("EMSG\\((.*), \"(.*)\"\\);");
        try {
            BufferedReader br = new BufferedReader(new FileReader(filePath));
            String         line;
            while ((line = br.readLine()) != null) {
                if (!line.trim().startsWith("EMSG")) {
                    continue;
                }
                String  codeMsg = line.trim();
                Matcher matcher = pattern.matcher(codeMsg);
                if (matcher.find()) {
                    String name = matcher.group(1);
                    String msg  = matcher.group(2);
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
            String         line;
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
}

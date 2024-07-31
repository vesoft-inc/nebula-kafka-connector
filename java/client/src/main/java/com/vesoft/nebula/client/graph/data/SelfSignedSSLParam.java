package com.vesoft.nebula.client.graph.data;

public class SelfSignedSSLParam extends SSLParam {
    private String crtFilePath;
    private String keyFilePath;
    private String password;

    public SelfSignedSSLParam(String crtFilePath, String keyFilePath, String password) {
        super(SignMode.SELF_SIGNED);
        this.crtFilePath = crtFilePath;
        this.keyFilePath = keyFilePath;
        this.password = password;
    }

    public String getCrtFilePath() {
        return crtFilePath;
    }

    public String getKeyFilePath() {
        return keyFilePath;
    }

    public String getPassword() {
        return password;
    }
}

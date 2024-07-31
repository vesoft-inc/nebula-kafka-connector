package com.vesoft.nebula.client.graph.data;

public class CASignedSSLParam extends SSLParam {
    private String caCrtFilePath;
    private String crtFilePath;
    private String keyFilePath;

    public CASignedSSLParam(String caCrtFilePath, String crtFilePath, String keyFilePath) {
        super(SignMode.CA_SIGNED);
        this.caCrtFilePath = caCrtFilePath;
        this.crtFilePath = crtFilePath;
        this.keyFilePath = keyFilePath;
    }

    public String getCaCrtFilePath() {
        return caCrtFilePath;
    }

    public String getCrtFilePath() {
        return crtFilePath;
    }

    public String getKeyFilePath() {
        return keyFilePath;
    }
}

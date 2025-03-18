/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.connector.mock;

import java.io.Serializable;

public class Follow implements Serializable {
    private String id;
    private String id1;
    private String id2;
    private Long followness;
    private Double likeness;

    public Follow(String id, String id1, String id2, Long followness, Double likeness) {
        this.id = id;
        this.id1 = id1;
        this.id2 = id2;
        this.followness = followness;
        this.likeness = likeness;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getId1() {
        return id1;
    }

    public void setId1(String id1) {
        this.id1 = id1;
    }

    public String getId2() {
        return id2;
    }

    public void setId2(String id2) {
        this.id2 = id2;
    }

    public Long getFollowness() {
        return followness;
    }

    public void setFollowness(Long followness) {
        this.followness = followness;
    }

    public Double getLikeness() {
        return likeness;
    }

    public void setLikeness(Double likeness) {
        this.likeness = likeness;
    }
}

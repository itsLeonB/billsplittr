package com.itsleonb.billsplittr.impl.mapper;

import com.itsleonb.billsplittr.api.entity.user.UserProfile;
import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;

public final class ProfileMapper {
  private ProfileMapper() {
    throw new UnsupportedOperationException("Utility class");
  }

  public static ProfileResponse toResponse(UserProfile profile) {
    if (profile == null) {
      throw new IllegalArgumentException("Profile input is null");
    }

    return ProfileResponse.builder()
      .name(profile.getName())
      .nickname(profile.getNickname())
      .avatar(profile.getAvatar())
      .bio(profile.getBio())
      .build();
  }
}

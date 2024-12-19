package com.itsleonb.billsplittr.impl.mapper;

import com.itsleonb.billsplittr.api.entity.user.UserProfile;
import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;

public class ProfileMapper {
  public static ProfileResponse toResponse(UserProfile profile) {
    return ProfileResponse.builder()
      .name(profile.getName())
      .nickname(profile.getNickname())
      .avatar(profile.getAvatar())
      .bio(profile.getBio())
      .build();
  }
}

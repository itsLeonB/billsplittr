package com.itsleonb.billsplittr.impl.service;

import com.itsleonb.billsplittr.api.entity.user.User;
import com.itsleonb.billsplittr.api.entity.user.UserProfile;
import com.itsleonb.billsplittr.api.exception.NotFoundException;
import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;
import com.itsleonb.billsplittr.api.model.profile.UpdateProfileRequest;
import com.itsleonb.billsplittr.api.repository.user.UserProfileRepository;
import com.itsleonb.billsplittr.api.repository.user.UserRepository;
import com.itsleonb.billsplittr.api.service.ProfileService;
import com.itsleonb.billsplittr.impl.mapper.ProfileMapper;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@AllArgsConstructor(onConstructor_ = @__(@Autowired))
public class ProfileServiceImpl implements ProfileService {
  private final UserRepository userRepository;
  private final UserProfileRepository profileRepository;

  @Override
  public ProfileResponse get(String email) {
    User user = userRepository.findByEmail(email).orElseThrow(() -> new NotFoundException("user not found"));

    UserProfile profile = profileRepository
      .findByUserId(user.getId())
      .orElseThrow(() -> new NotFoundException("profile not found"));

    return ProfileMapper.toResponse(profile);
  }

  @Override
  @Transactional
  public ProfileResponse update(String email, UpdateProfileRequest request) {
    User user = userRepository.findByEmail(email).orElseThrow(() -> new NotFoundException("user not found"));

    UserProfile profile = profileRepository
      .findByUserId(user.getId())
      .orElseThrow(() -> new NotFoundException("profile not found"));

    UserProfile profileToUpdate = profile.toBuilder()
      .name(request.getName())
      .nickname(request.getNickname())
      .avatar(request.getAvatar())
      .bio(request.getBio())
      .build();

    UserProfile updatedProfile = profileRepository.save(profileToUpdate);

    return ProfileMapper.toResponse(updatedProfile);
  }
}

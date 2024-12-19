package com.itsleonb.billsplittr.api.service;

import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;
import com.itsleonb.billsplittr.api.model.profile.UpdateProfileRequest;

public interface ProfileService {
  ProfileResponse get(String email);

  ProfileResponse update(String email, UpdateProfileRequest request);
}

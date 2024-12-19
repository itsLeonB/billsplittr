package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.controller.ProfileController;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;
import com.itsleonb.billsplittr.api.model.profile.UpdateProfileRequest;
import com.itsleonb.billsplittr.api.service.ProfileService;
import com.itsleonb.billsplittr.impl.util.SecurityUtil;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

@RestController
@AllArgsConstructor(onConstructor = @__(@Autowired))
public class ProfileControllerImpl implements ProfileController {
  private final SecurityUtil securityUtil;
  private final ProfileService profileService;

  @Override
  public JsonResponse<ProfileResponse> handleGet() {
    String email = securityUtil.getCurrentUserEmail();
    ProfileResponse response = profileService.get(email);

    return JsonResponse.<ProfileResponse>builder()
      .success(true)
      .data(response)
      .build();
  }

  @Override
  public JsonResponse<ProfileResponse> handleUpdate(UpdateProfileRequest request) {
    String email = securityUtil.getCurrentUserEmail();
    ProfileResponse response = profileService.update(email, request);

    return JsonResponse.<ProfileResponse>builder()
      .success(true)
      .data(response)
      .build();
  }
}

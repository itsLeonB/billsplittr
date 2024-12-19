package com.itsleonb.billsplittr.api.controller;

import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.profile.ProfileResponse;
import com.itsleonb.billsplittr.api.model.profile.UpdateProfileRequest;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

@RequestMapping("/profile")
public interface ProfileController {
  @GetMapping(
    path = "",
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<ProfileResponse> handleGet();

  @PutMapping(
    path = "",
    consumes = MediaType.APPLICATION_JSON_VALUE,
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<ProfileResponse> handleUpdate(@RequestBody UpdateProfileRequest request);
}

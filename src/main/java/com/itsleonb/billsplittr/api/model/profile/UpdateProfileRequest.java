package com.itsleonb.billsplittr.api.model.profile;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UpdateProfileRequest {
  private String name;
  private String nickname;
  private String avatar;
  private String bio;
}

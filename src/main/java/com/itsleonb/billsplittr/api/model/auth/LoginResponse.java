package com.itsleonb.billsplittr.api.model.auth;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class LoginResponse {
  @Builder.Default
  private String type = "Bearer";

  private String token;
}

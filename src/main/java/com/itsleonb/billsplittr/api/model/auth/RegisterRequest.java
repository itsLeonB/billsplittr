package com.itsleonb.billsplittr.api.model.auth;

import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class RegisterRequest {
  @NotBlank
  @Email(message = "Please provide a valid email address")
  private String email;

  @NotBlank
  @Min(value = 8)
  private String password;

  @NotBlank
  private String passwordConfirmed;

  @AssertTrue(message = "Passwords do not match")
  private boolean isPasswordConfirmed() {
    return password.equals(passwordConfirmed);
  }
}

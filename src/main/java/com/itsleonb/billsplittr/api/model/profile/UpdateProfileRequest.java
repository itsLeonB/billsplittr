package com.itsleonb.billsplittr.api.model.profile;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.hibernate.validator.constraints.URL;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UpdateProfileRequest {
  @NotBlank(message = "Name is required")
  @Size(max = 100, message = "Name must be less than 100 characters")
  @Pattern(regexp = "^[\\p{L}\\s.'\\-]+$", message = "Name contains invalid characters")
  private String name;

  @Size(max = 50, message = "Nickname must be less than 50 characters")
  @Pattern(regexp = "^[\\w\\s.\\-]*$", message = "Nickname contains invalid characters")
  private String nickname;

  @URL(message = "Avatar must be a valid URL")
  @Size(max = 255, message = "Avatar URL is too long")
  private String avatar;

  @Size(max = 500, message = "Bio must be less than 500 characters")
  private String bio;
}

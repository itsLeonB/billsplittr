package com.itsleonb.billsplittr.api.model.merchant;

import com.itsleonb.billsplittr.api.constant.MerchantType;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NewMerchantRequest {
  @NotBlank
  private String name;

  @NotNull
  private MerchantType type;

  @NotBlank
  private String address;
}

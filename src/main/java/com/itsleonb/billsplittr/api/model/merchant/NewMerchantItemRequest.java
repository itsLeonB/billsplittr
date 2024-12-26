package com.itsleonb.billsplittr.api.model.merchant;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NewMerchantItemRequest {
  @NotBlank
  private String name;

  @NotNull
  private BigDecimal price;
}

package com.itsleonb.billsplittr.api.model.merchant;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.itsleonb.billsplittr.api.constant.MerchantType;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class MerchantResponse {
  private UUID id;
  private String name;
  private MerchantType type;
  private String address;
}

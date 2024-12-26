package com.itsleonb.billsplittr.impl.mapper;

import com.itsleonb.billsplittr.api.entity.merchant.Merchant;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;

import java.util.List;
import java.util.stream.Collectors;

public class MerchantMapper {
  public static Merchant fromNewRequest(NewMerchantRequest request) {
    return Merchant.builder()
      .name(request.getName())
      .type(request.getType())
      .address(request.getAddress())
      .build();
  }

  public static MerchantResponse toResponse(Merchant merchant) {
    return MerchantResponse.builder()
      .id(merchant.getId())
      .name(merchant.getName())
      .type(merchant.getType())
      .address(merchant.getAddress())
      .build();
  }

  public static List<MerchantResponse> toResponses(List<Merchant> merchants) {
    return merchants.stream()
      .map(MerchantMapper::toResponse)
      .collect(Collectors.toList());
  }
}

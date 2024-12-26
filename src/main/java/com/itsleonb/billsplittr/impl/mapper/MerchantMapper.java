package com.itsleonb.billsplittr.impl.mapper;

import com.itsleonb.billsplittr.api.entity.merchant.Merchant;
import com.itsleonb.billsplittr.api.entity.merchant.MerchantItem;
import com.itsleonb.billsplittr.api.model.merchant.MerchantItemResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantItemRequest;
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
      .items(toItemResponses(merchant.getItems()))
      .build();
  }

  public static List<MerchantResponse> toResponses(List<Merchant> merchants) {
    return merchants.stream()
      .map(MerchantMapper::toResponse)
      .collect(Collectors.toList());
  }

  public static MerchantItem fromNewItemRequest(Merchant merchant, NewMerchantItemRequest request) {
    return MerchantItem.builder()
      .merchant(merchant)
      .name(request.getName())
      .price(request.getPrice())
      .build();
  }

  public static MerchantItemResponse toItemResponse(MerchantItem item) {
    return MerchantItemResponse.builder()
      .id(item.getId())
      .name(item.getName())
      .price(item.getPrice())
      .build();
  }

  public static List<MerchantItemResponse> toItemResponses(List<MerchantItem> items) {
    if (items == null) {
      return null;
    }

    return items.stream()
      .map(MerchantMapper::toItemResponse)
      .collect(Collectors.toList());
  }
}

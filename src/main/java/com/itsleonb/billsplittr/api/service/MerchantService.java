package com.itsleonb.billsplittr.api.service;

import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;

import java.util.List;
import java.util.UUID;

public interface MerchantService {
  MerchantResponse create(NewMerchantRequest request);

  List<MerchantResponse> find(String name);

  MerchantResponse getById(UUID id);
}

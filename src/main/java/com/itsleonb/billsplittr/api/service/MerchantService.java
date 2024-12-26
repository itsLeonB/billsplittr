package com.itsleonb.billsplittr.api.service;

import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;

public interface MerchantService {
  MerchantResponse create(NewMerchantRequest request);
}

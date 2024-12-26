package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.controller.MerchantController;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;
import com.itsleonb.billsplittr.api.service.MerchantService;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

@RestController
@AllArgsConstructor(onConstructor = @__(@Autowired))
public class MerchantControllerImpl implements MerchantController {
  private MerchantService merchantService;

  @Override
  public JsonResponse<MerchantResponse> handleCreate(NewMerchantRequest request) {
    MerchantResponse response = merchantService.create(request);

    return JsonResponse.<MerchantResponse>builder()
      .success(true)
      .data(response)
      .build();
  }
}

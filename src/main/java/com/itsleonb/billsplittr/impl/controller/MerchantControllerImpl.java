package com.itsleonb.billsplittr.impl.controller;

import com.itsleonb.billsplittr.api.controller.MerchantController;
import com.itsleonb.billsplittr.api.exception.BadRequestException;
import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;
import com.itsleonb.billsplittr.api.service.MerchantService;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.UUID;

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

  @Override
  public JsonResponse<List<MerchantResponse>> handleFind(String name) {
    if (name.length() < 3) {
      throw new BadRequestException("Please input at least three characters to search");
    }

    List<MerchantResponse> responses = merchantService.find(name);

    return JsonResponse.<List<MerchantResponse>>builder()
      .success(true)
      .data(responses)
      .build();
  }

  @Override
  public JsonResponse<MerchantResponse> handleGetById(UUID id) {
    MerchantResponse response = merchantService.getById(id);

    return JsonResponse.<MerchantResponse>builder()
      .success(true)
      .data(response)
      .build();
  }
}

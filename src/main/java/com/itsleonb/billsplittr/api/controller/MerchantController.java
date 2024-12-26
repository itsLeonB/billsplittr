package com.itsleonb.billsplittr.api.controller;

import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantItemResponse;
import com.itsleonb.billsplittr.api.model.merchant.MerchantResponse;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantItemRequest;
import com.itsleonb.billsplittr.api.model.merchant.NewMerchantRequest;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;
import java.util.UUID;

@RequestMapping("/merchants")
public interface MerchantController {
  @PostMapping(
    path = "",
    consumes = MediaType.APPLICATION_JSON_VALUE,
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<MerchantResponse> handleCreate(@RequestBody NewMerchantRequest request);

  @GetMapping(
    path = "",
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<List<MerchantResponse>> handleFind(@RequestParam String name);

  @GetMapping(
    path = "/{id}",
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<MerchantResponse> handleGetById(@PathVariable UUID id);

  @PostMapping(
    path = "/{id}/items",
    consumes = MediaType.APPLICATION_JSON_VALUE,
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<MerchantItemResponse> handleCreateItem(
    @PathVariable UUID id,
    @RequestBody NewMerchantItemRequest request
  );
}

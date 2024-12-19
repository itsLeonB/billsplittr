package com.itsleonb.billsplittr.api.controller;

import com.itsleonb.billsplittr.api.model.JsonResponse;
import com.itsleonb.billsplittr.api.model.auth.LoginRequest;
import com.itsleonb.billsplittr.api.model.auth.LoginResponse;
import com.itsleonb.billsplittr.api.model.auth.RegisterRequest;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

@RequestMapping("/auth")
public interface AuthController {
  @PostMapping(
    path = "/register",
    consumes = MediaType.APPLICATION_JSON_VALUE,
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<String> handleRegister(@RequestBody RegisterRequest request);

  @PostMapping(
    path = "/login",
    consumes = MediaType.APPLICATION_JSON_VALUE,
    produces = MediaType.APPLICATION_JSON_VALUE
  )
  JsonResponse<LoginResponse> handleLogin(@RequestBody LoginRequest request);
}
